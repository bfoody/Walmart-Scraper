package hub

import (
	"fmt"
	"time"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"go.uber.org/zap"
)

// A Heartbeater maintains a connection to a server and sends heartbeats,
// reporting back on server failure.
type Heartbeater struct {
	sender   *identity.Server // the server sending heartbeats
	receiver *identity.Server // the server receiving heartbeats
	interval time.Duration    // the interval between heartbeats
	conn     *communication.QueueConnection
	shutdown chan int
	timer    *time.Timer
	log      *zap.Logger
}

// NewHeartbeater creates and returns a new *Heartbeater.
func NewHeartbeater(sender *identity.Server, receiver *identity.Server, interval time.Duration, conn *communication.QueueConnection, logger *zap.Logger) *Heartbeater {
	return &Heartbeater{
		sender:   sender,
		receiver: receiver,
		interval: interval,
		conn:     conn,
		shutdown: make(chan int),
		timer:    nil,
		log:      logger,
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
			h.sendHeartbeat()
			// Restart the timer.
			h.timer.Reset(h.interval)
		case <-h.shutdown:
			h.timer.Stop()

			return
		}
	}
}

// sendHeartbeat sends the heartbeat message to the server.
func (h *Heartbeater) sendHeartbeat() {
	message := communication.Heartbeat{
		SingleReceiverPacket: communication.SingleReceiverPacket{
			SenderID:   h.sender.ID,
			ReceiverID: h.receiver.ID,
		},
		ResponseExpected: true,
	}

	if err := h.conn.SendMessage(message); err != nil {
		h.log.Error(fmt.Sprintf("error sending heartbeat to server %s", h.receiver.ID), zap.Error(err))
		return
	}

	h.log.Debug(fmt.Sprintf("sending heartbeat to server %s", h.receiver.ID))
}
