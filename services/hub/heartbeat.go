package hub

import (
	"fmt"
	"time"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"go.uber.org/zap"
)

const (
	// MissedBeatsAllowed represents the number of missed beats allowed before the
	// client will be considered unresponsive and disconnected.
	MissedBeatsAllowed = 4
)

// A Heartbeater maintains a connection to a server and sends heartbeats,
// reporting back on server failure.
type Heartbeater struct {
	sender      *identity.Server // the server sending heartbeats
	receiver    *identity.Server // the server receiving heartbeats
	interval    time.Duration    // the interval between heartbeats
	conn        *communication.QueueConnection
	shutdown    chan int
	timer       *time.Timer
	log         *zap.Logger
	beatsMissed uint8 // the number of heartbeats missed
}

// NewHeartbeater creates and returns a new *Heartbeater.
func NewHeartbeater(sender *identity.Server, receiver *identity.Server, interval time.Duration, conn *communication.QueueConnection, logger *zap.Logger) *Heartbeater {
	return &Heartbeater{
		sender:      sender,
		receiver:    receiver,
		interval:    interval,
		conn:        conn,
		shutdown:    make(chan int),
		timer:       nil,
		log:         logger,
		beatsMissed: 0,
	}
}

// Start starts the Heartbeater.
func (h *Heartbeater) Start() error {
	h.timer = time.NewTimer(h.interval)
	go h.loop()

	return nil
}

// Shutdown stops the Heartbeater and prevents any further heartbeats
// from being sent.
func (h *Heartbeater) Shutdown() error {
	h.shutdown <- 1
	return nil
}

func (h *Heartbeater) loop() {
	for {
		select {
		case <-h.timer.C:
			h.sendHeartbeat(true)
			// Restart the timer.
			h.timer.Reset(h.interval)
		case <-h.shutdown:
			h.timer.Stop()

			return
		}
	}
}

// sendHeartbeat sends the heartbeat message to the server.
func (h *Heartbeater) sendHeartbeat(responseExpected bool) {
	message := communication.Heartbeat{
		SingleReceiverPacket: communication.SingleReceiverPacket{
			SenderID:   h.sender.ID,
			ReceiverID: h.receiver.ID,
		},
		ResponseExpected: responseExpected,
	}

	if err := h.conn.SendMessage(message); err != nil {
		h.log.Error(fmt.Sprintf("error sending heartbeat to server %s", h.receiver.ID), zap.Error(err))
		return
	}

	if responseExpected {
		h.beatsMissed++
	}

	h.log.Debug(fmt.Sprintf("sending heartbeat to server %s", h.receiver.ID))
}

// HandleHeartbeat handles the receipt of a Heartbeat from the client server.
func (h *Heartbeater) HandleHeartbeat(hb *communication.Heartbeat) {
	h.beatsMissed = 0
}
