package supervisor

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/identity"
	"github.com/bfoody/Walmart-Scraper/services/hub"
	"go.uber.org/zap"
)

// HeartbeatInterval is the amount of time between each heartbeat.
const HeartbeatInterval = 3 * time.Second

// A ServerMap stores a list of servers and their statuses.
type ServerMap map[string]ServerStatus

// A Supervisor maintains a list of currently connected servers and their
// statuses.
type Supervisor struct {
	identity       *identity.Server
	conn           *communication.QueueConnection
	service        hub.Service
	serverMapMutex *sync.RWMutex
	serverMap      map[string]ServerStatus
	heartbeaters   map[string]*hub.Heartbeater
	statusUpdates  chan communication.StatusUpdate
	heartbeats     chan communication.Heartbeat
	goingAways     chan communication.GoingAway
	infoRetrieved  chan communication.InfoRetrieved
	serverDown     chan identity.Server // any servers sent through this channel will be considered offline
	shutdown       chan int
	log            *zap.Logger
	taskManager    *TaskManager
	roundRobin     *RoundRobin
}

// New creates and returns a new *Supervisor.
func New(_identity *identity.Server, logger *zap.Logger, conn *communication.QueueConnection, service hub.Service) *Supervisor {
	tm := NewTaskManager(service)

	return &Supervisor{
		identity:       _identity,
		conn:           conn,
		service:        service,
		serverMapMutex: &sync.RWMutex{},
		serverMap:      map[string]ServerStatus{},
		heartbeaters:   map[string]*hub.Heartbeater{},
		statusUpdates:  make(chan communication.StatusUpdate, 4),
		heartbeats:     make(chan communication.Heartbeat, 4),
		goingAways:     make(chan communication.GoingAway, 4),
		infoRetrieved:  make(chan communication.InfoRetrieved, 4),
		serverDown:     make(chan identity.Server, 4),
		shutdown:       make(chan int),
		log:            logger,
		taskManager:    tm,
		roundRobin:     NewRoundRobin(),
	}
}

// Start starts the Supervisor.
func (s *Supervisor) Start() error {
	s.conn.RegisterStatusUpdateHandler(s.pipeStatusUpdate)
	s.conn.RegisterHeartbeatHandler(s.pipeHeartbeat)
	s.conn.RegisterGoingAwayHandler(s.pipeGoingAway)
	s.conn.RegisterInfoRetrievedHandler(s.pipeInfoRetrieved)

	err := s.taskManager.Initialize()
	if err != nil {
		return err
	}
	s.taskManager.Start(s.taskCallback)

	go s.loop()
	return nil
}

// taskCallback is called by the TaskManager when a task is due to be dispatched.
func (s *Supervisor) taskCallback(task domain.ScrapeTask) {
	s.distributeTask(task)
}

// distributeTask distributes a task to a client server in a round-robin fashion.
func (s *Supervisor) distributeTask(task domain.ScrapeTask) {
	s.serverMapMutex.RLock()
	defer s.serverMapMutex.RUnlock()

	// Create an array of server IDs to choose from.
	serverIDArray := []string{}
	for id := range s.serverMap {
		serverIDArray = append(serverIDArray, id)
	}

	// Get the ID of the server chosen by round-robin.
	idx := s.roundRobin.Next(uint(len(serverIDArray)))
	id := serverIDArray[idx]

	go func() {
		// TODO: handle error
		pl, err := s.service.GetProductLocationByID(task.ProductLocationID)
		if err != nil {
			s.log.Error("Error getting ProductLocation for TaskFulfillmentRequest", zap.String("productLocationID", task.ProductLocationID), zap.Error(err))
			return
		}

		req := communication.TaskFulfillmentRequest{
			SingleReceiverPacket: communication.SingleReceiverPacket{
				SenderID:   s.identity.ID,
				ReceiverID: id,
			},
			TaskID:          task.ID,
			ProductLocation: *pl,
		}

		err = s.conn.SendMessage(req)
		if err != nil {
			s.log.Error("Error sending TaskFulfillmentRequest to server", zap.String("serverID", id), zap.Error(err))
		}
	}()
}

// Shutdown shuts down the Supervisor.
func (s *Supervisor) Shutdown() error {
	s.shutdown <- 1
	return nil
}

// pipeStatusUpdate pipes a StatusUpdate into the supervisor.
func (s *Supervisor) pipeStatusUpdate(su *communication.StatusUpdate) {
	s.statusUpdates <- *su
}

// pipeHeartbeat pipes a Heartbeat into the supervisor.
func (s *Supervisor) pipeHeartbeat(hb *communication.Heartbeat) {
	s.heartbeats <- *hb
}

// pipeGoingAway pipes a GoingAway into the supervisor.
func (s *Supervisor) pipeGoingAway(ga *communication.GoingAway) {
	s.goingAways <- *ga
}

// pipeInfoRetrieved pipes an InfoRetrieved into the supervisor.
func (s *Supervisor) pipeInfoRetrieved(ir *communication.InfoRetrieved) {
	s.infoRetrieved <- *ir
}

func (s *Supervisor) loop() {
	for {
		select {
		case su := <-s.statusUpdates:
			s.handleStatusUpdate(&su)
		case hb := <-s.heartbeats:
			s.handleHeartbeat(&hb)
		case ga := <-s.goingAways:
			s.handleGoingAway(&ga)
		case ir := <-s.infoRetrieved:
			s.handleInfoRetrieved(&ir)
		case server := <-s.serverDown:
			s.terminateServer(&server)
		case <-s.shutdown:
			s.cleanup()
			return
		}
	}
}

// cleanup gracefully shuts down the Supervisor.
func (s *Supervisor) cleanup() {
	for id, hb := range s.heartbeaters {
		if hb != nil {
			err := hb.Shutdown()
			if err != nil {
				s.log.Error(fmt.Sprintf("error occurred while shutting down heartbeater for server %s", id), zap.Error(err))
			}
		}
	}
}

func (s *Supervisor) handleStatusUpdate(su *communication.StatusUpdate) {
	server := identity.NewClient(su.SenderID)
	status := ServerStatus{
		AvailableForWork: su.AvailableForWork,
	}

	// Compare the new status with the old one and log changes.
	if oldStatus, ok := s.serverMap[su.SenderID]; !ok || !reflect.DeepEqual(status, oldStatus) {
		s.log.Info(
			fmt.Sprintf("status changed for server %s", su.SenderID),
			zap.String("status", fmt.Sprintf("%+v", status)),
		)
	}

	if _, ok := s.serverMap[su.SenderID]; !ok {
		s.heartbeaters[server.ID] = hub.NewHeartbeater(s.identity, server, HeartbeatInterval, s.serverDown, s.conn, s.log)
		if err := s.heartbeaters[su.SenderID].Start(); err != nil {
			s.log.Error(
				fmt.Sprintf("error occurred starting heartbeater for server %s", su.SenderID),
				zap.Error(err),
			)
		}
	}

	s.serverMapMutex.Lock()
	defer s.serverMapMutex.Unlock()

	// Replace the status with the new one.
	s.serverMap[su.SenderID] = status

	err := s.conn.SendMessage(communication.HubWelcome{
		SingleReceiverPacket: communication.SingleReceiverPacket{
			SenderID:   s.identity.ID,
			ReceiverID: su.SenderID,
		},
	})
	if err != nil {
		s.log.Error(
			fmt.Sprintf("error occurred sending HubWelcome to server %s", su.SenderID),
			zap.Error(err),
		)
	}
}

func (s *Supervisor) handleHeartbeat(hb *communication.Heartbeat) {
	if hb.SenderID == s.identity.ID || hb.ReceiverID != s.identity.ID {
		return
	}

	server := identity.NewClient(hb.SenderID)
	h, ok := s.heartbeaters[server.ID]
	if !ok {
		s.log.Error(fmt.Sprintf("server %s missing heartbeater, this shouldn't happen", server.ID))
		return
	}

	h.HandleHeartbeat(hb)

	s.log.Debug(
		fmt.Sprintf("received heartbeat from server %s", server.ID),
	)
}

func (s *Supervisor) handleGoingAway(ga *communication.GoingAway) {
	if ga.SenderID == s.identity.ID || ga.ReceiverID != s.identity.ID {
		return
	}

	server := identity.NewClient(ga.SenderID)
	s.serverDown <- *server
}

func (s *Supervisor) handleInfoRetrieved(ir *communication.InfoRetrieved) {
	if ir.SenderID == s.identity.ID || ir.ReceiverID != s.identity.ID {
		return
	}

	pi := ir.ProductInfo

	id, err := s.service.SaveProductInfo(pi)
	if err != nil {
		s.log.Error(fmt.Sprintf("error saving product info for task %s", ir.TaskID), zap.Error(err))
		return
	}

	err = s.service.ResolveTask(ir.TaskID)
	if err != nil {
		s.log.Error(fmt.Sprintf("error resolving task %s", ir.TaskID), zap.Error(err))
	}

	s.log.Debug("product info saved for task", zap.String("taskId", ir.TaskID), zap.String("productInfoId", id))
}

// terminateServer removes a single server from the supervisor and shuts down all listeners
// attached to it.
func (s *Supervisor) terminateServer(server *identity.Server) {
	s.serverMapMutex.Lock()
	defer s.serverMapMutex.Unlock()

	s.log.Debug(fmt.Sprintf("disconnecting from server %s", server.ID))
	err := s.heartbeaters[server.ID].Shutdown()
	if err != nil {
		s.log.Error(fmt.Sprintf("error occurred while shutting down heartbeater for server %s", server.ID), zap.Error(err))
	}

	delete(s.serverMap, server.ID)

	s.log.Info(fmt.Sprintf("disconnected from server %s", server.ID))
}
