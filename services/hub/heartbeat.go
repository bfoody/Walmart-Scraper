package hub

import (
	"time"

	"github.com/bfoody/Walmart-Scraper/communication"
)

// A Heartbeater maintains a connection to a server and sends heartbeats,
// reporting back on server failure.
type Heartbeater struct {
	serverID  string        // the server being sent heartbeats
	interval  time.Duration // the interval between heartbeats
	conn      *communication.QueueConnection
	queueName string
	shutdown  chan int
	timer     *time.Timer
}

// NewHeartbeater creates and returns a new *Heartbeater.
func NewHeartbeater(serverID string, interval time.Duration, conn *communication.QueueConnection, queueName string) *Heartbeater {
	return &Heartbeater{
		serverID:  serverID,
		interval:  interval,
		conn:      conn,
		queueName: queueName,
		shutdown:  make(chan int),
		timer:     nil,
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
	h.conn.SendMessage()
}
