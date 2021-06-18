package receiver

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"go.uber.org/zap"
)

// A Receiver processes and responds to messages from the hub server.
type Receiver struct {
	identity         *identity.Server
	heartbeats       chan communication.Heartbeat
	newHubIdentities chan identity.Server
	hub              *identity.Server // the hub that the client is currently connected to
	conn             *communication.QueueConnection
	hubWelcomes      chan communication.HubWelcome
	shutdown         chan int
	log              *zap.Logger
}

// New creates and returns a new *Receiver.
func New(_identity *identity.Server, logger *zap.Logger, conn *communication.QueueConnection) *Receiver {
	return &Receiver{
		identity:         _identity,
		newHubIdentities: make(chan identity.Server),
		hub:              nil,
		conn:             conn,
		shutdown:         make(chan int),
		log:              logger,
	}
}

// Start starts the Receiver and enters the main loop in a Goroutine.
func (r *Receiver) Start() error {
	r.conn.RegisterHubWelcomeHandler(r.pipeHubWelcome)
	r.conn.RegisterHeartbeatHandler(r.pipeHeartbeat)

	go r.loop()
	return nil
}

func (r *Receiver) Shutdown() error {
	r.shutdown <- 1

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

func (r *Receiver) loop() {
	for {
		select {
		case hw := <-r.hubWelcomes:
			r.handleHubWelcome(&hw)
		case hub := <-r.newHubIdentities:
			r.switchHub(&hub)
		case hb := <-r.heartbeats:
			r.handleHeartbeat(&hb)
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

// switchHub switches the client to communicate with the specified hub identity.
func (r *Receiver) switchHub(hub *identity.Server) {
	r.log.Info(fmt.Sprintf("switching hub to hub %s", hub.ID))
	r.hub = hub
}

// cleanup prepares the Receiver for shutdown and notifies
// the hub that the client is going away.
func (r *Receiver) cleanup() {

}
