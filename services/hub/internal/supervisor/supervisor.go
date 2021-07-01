package supervisor

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/bfoody/Walmart-Scraper/communication"
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
	serverMapMutex *sync.RWMutex
	serverMap      map[string]ServerStatus
	heartbeaters   map[string]*hub.Heartbeater
	statusUpdates  chan communication.StatusUpdate
	heartbeats     chan communication.Heartbeat
	serverDown     chan identity.Server // any servers sent through this channel will be considered offline
	shutdown       chan int
	log            *zap.Logger
}

// New creates and returns a new *Supervisor.
func New(_identity *identity.Server, logger *zap.Logger, conn *communication.QueueConnection) *Supervisor {
	return &Supervisor{
		identity:       _identity,
		conn:           conn,
		serverMapMutex: &sync.RWMutex{},
		serverMap:      map[string]ServerStatus{},
		heartbeaters:   map[string]*hub.Heartbeater{},
		statusUpdates:  make(chan communication.StatusUpdate, 4),
		heartbeats:     make(chan communication.Heartbeat, 4),
		serverDown:     make(chan identity.Server, 4),
		shutdown:       make(chan int),
		log:            logger,
	}
}

// Start starts the Supervisor.
func (s *Supervisor) Start() error {
	s.conn.RegisterStatusUpdateHandler(s.pipeStatusUpdate)
	s.conn.RegisterHeartbeatHandler(s.pipeHeartbeat)

	go s.loop()
	return nil
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

func (s *Supervisor) loop() {
	for {
		select {
		case su := <-s.statusUpdates:
			s.handleStatusUpdate(&su)
		case hb := <-s.heartbeats:
			s.handleHeartbeat(&hb)
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
}

func (s *Supervisor) handleHeartbeat(hb *communication.Heartbeat) {
	if hb.SenderID == s.identity.ID {
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
