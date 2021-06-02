package receiver

import (
	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/identity"
)

// A Receiver processes and responds to messages from the hub server.
type Receiver struct {
	identity *identity.Server
	conn     *communication.QueueConnection
	shutdown chan int
	log      *zap.Logger
}

// New creates and returns a new *Receiver..
func New(_identity *identity.Server, logger *zap.Logger, conn *communication.QueueConnection) *Supervisor {
	return &Supervisor{
		identity: _identity,
		conn:     conn,
		shutdown: make(chan int),
		log:      logger,
	}
}
