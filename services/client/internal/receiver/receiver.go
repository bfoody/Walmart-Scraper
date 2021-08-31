package receiver

import (
	"fmt"
	"sync"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"go.uber.org/zap"
)

const (
	// ReasonShuttingDown signifies that the server is gracefully shutting down.
	ReasonShuttingDown = "SHUTTING_DOWN"
)

// A Receiver processes and responds to messages from the hub server.
type Receiver struct {
	identity                *identity.Server
	heartbeats              chan communication.Heartbeat
	newHubIdentities        chan identity.Server
	hub                     *identity.Server // the hub that the client is currently connected to
	conn                    *communication.QueueConnection
	taskService             *TaskService
	hubWelcomes             chan communication.HubWelcome
	taskFulfillmentRequests chan communication.TaskFulfillmentRequest
	shutdown                chan int
	shutdownWg              *sync.WaitGroup
	log                     *zap.Logger
}

// New creates and returns a new *Receiver.
func New(_identity *identity.Server, logger *zap.Logger, conn *communication.QueueConnection) *Receiver {
	return &Receiver{
		identity:                _identity,
		heartbeats:              make(chan communication.Heartbeat),
		newHubIdentities:        make(chan identity.Server, 4),
		hub:                     nil,
		conn:                    conn,
		taskService:             NewTaskService(logger),
		hubWelcomes:             make(chan communication.HubWelcome, 4),
		taskFulfillmentRequests: make(chan communication.TaskFulfillmentRequest, 4),
		shutdown:                make(chan int),
		shutdownWg:              &sync.WaitGroup{},
		log:                     logger,
	}
}

// Start starts the Receiver and enters the main loop in a Goroutine.
func (r *Receiver) Start() error {
	r.conn.RegisterHubWelcomeHandler(r.pipeHubWelcome)
	r.conn.RegisterHeartbeatHandler(r.pipeHeartbeat)
	r.conn.RegisterTaskFulfillmentRequest(r.pipeTaskFulfillmentRequest)

	go r.loop()
	return nil
}

func (r *Receiver) Shutdown() error {
	r.shutdownWg.Add(1)
	r.shutdown <- 1

	r.shutdownWg.Wait()
	return nil
}

// pipeHubWelcome pipes a HubWelcome into the receiver.
func (r *Receiver) pipeHubWelcome(hw *communication.HubWelcome) {
	r.hubWelcomes <- *hw
}

// pipeHeartbeat pipes a Heartbeat into the receiver.
func (r *Receiver) pipeHeartbeat(hb *communication.Heartbeat) {
	r.heartbeats <- *hb
}

// pipeTaskFulfillmentRequest pipes a TaskFulfillmentRequest into the receiver.
func (r *Receiver) pipeTaskFulfillmentRequest(tfr *communication.TaskFulfillmentRequest) {
	r.taskFulfillmentRequests <- *tfr
}

func (r *Receiver) loop() {
	for {
		select {
		case hw := <-r.hubWelcomes:
			r.handleHubWelcome(&hw)
		case hub := <-r.newHubIdentities:
			r.switchHub(&hub)
		case hb := <-r.heartbeats:
			r.handleHeartbeat(&hb)
		case tfr := <-r.taskFulfillmentRequests:
			r.handleTaskFulfillmentRequest(&tfr)
		case <-r.shutdown:
			r.cleanup()
			return
		}
	}
}

func (r *Receiver) handleHubWelcome(hw *communication.HubWelcome) {
	r.newHubIdentities <- *identity.NewHub(hw.SenderID)
}

func (r *Receiver) handleHeartbeat(hb *communication.Heartbeat) {
	// TODO: check receiver ID in a better way
	if hb.ReceiverID != r.identity.ID {
		return
	}

	message := communication.Heartbeat{
		SingleReceiverPacket: communication.SingleReceiverPacket{
			SenderID:   r.identity.ID,
			ReceiverID: hb.SenderID,
		},
		ResponseExpected: false,
	}

	if err := r.conn.SendMessage(message); err != nil {
		r.log.Error(fmt.Sprintf("error sending heartbeat to hub %s", hb.SenderID), zap.Error(err))
	}

	r.log.Debug(fmt.Sprintf("sending heartbeat to hub %s", hb.SenderID))
}

func (r *Receiver) handleTaskFulfillmentRequest(tfr *communication.TaskFulfillmentRequest) {
	// TODO: check receiver ID in a better way
	if tfr.ReceiverID != r.identity.ID {
		return
	}

	// TODO: possibly implement a thread pool for this?
	go r.runTask(tfr)
}

func (r *Receiver) runTask(tfr *communication.TaskFulfillmentRequest) {
	pi, err := r.taskService.FetchProductInfo(&tfr.ProductLocation)
	if err != nil {
		r.log.Error(
			"couldn't fetch product info, rescheduling to next interval",
			zap.String("productLocationId", tfr.ProductLocation.ID),
			zap.Error(err),
		)

		return
	}

	ir := communication.InfoRetrieved{
		SingleReceiverPacket: communication.SingleReceiverPacket{
			SenderID:   r.identity.ID,
			ReceiverID: r.hub.ID,
		},
		TaskID:      tfr.TaskID,
		ProductInfo: *pi,
	}

	err = r.conn.SendMessage(ir)
	if err != nil {
		r.log.Error(
			"couldn't send InfoRetrieved message to hub",
			zap.String("productLocationId", tfr.ProductLocation.ID),
			zap.String("hubId", r.hub.ID),
			zap.Error(err),
		)
	}
}

// switchHub switches the client to communicate with the specified hub identity.
func (r *Receiver) switchHub(hub *identity.Server) {
	r.log.Info(fmt.Sprintf("switching hub to hub %s", hub.ID))
	r.hub = hub
}

// cleanup prepares the Receiver for shutdown and notifies
// the hub that the client is going away.
func (r *Receiver) cleanup() {
	defer r.shutdownWg.Done()

	err := r.conn.SendMessage(communication.GoingAway{
		SingleReceiverPacket: communication.SingleReceiverPacket{
			SenderID:   r.identity.ID,
			ReceiverID: r.hub.ID,
		},
		Reason: ReasonShuttingDown,
	})
	if err != nil {
		r.log.Error(
			fmt.Sprintf("error sending GoingAway to hub %s", r.hub.ID),
			zap.Error(err),
		)
	}
}
