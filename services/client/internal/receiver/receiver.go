package receiver

import (
	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
	"go.uber.org/zap"
)

// A Receiver processes and responds to messages from the hub server.
type Receiver struct {
	identity *identity.Server
	conn     *communication.QueueConnection
	shutdown chan int
	log      *zap.Logger
}

// New creates and returns a new *Receiver..
func New(_identity *identity.Server, logger *zap.Logger, conn *communication.QueueConnection) *Receiver {
	return &Receiver{
		identity: _identity,
		conn:     conn,
		shutdown: make(chan int),
		log:      logger,
	}
}

// Start starts the Receiver and enters the main loop in a Goroutine.
func (r *Receiver) Start() error {
	go r.loop()
	return nil
}

func (r *Receiver) loop() {
	for {
		select {
		case <-r.shutdown:
			r.cleanup()
			return
		}
	}
}

// cleanup prepares the Receiver for shutdown and notifies
// the hub that the client is going away.
func (r *Receiver) cleanup() {

}
