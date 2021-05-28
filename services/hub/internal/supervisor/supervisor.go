package supervisor

import (
	"fmt"
	"sync"

	"github.com/bfoody/Walmart-Scraper/communication"
)

// A ServerMap stores a list of servers and their statuses.
type ServerMap map[string]ServerStatus

// A Supervisor maintains a list of currently connected servers and their
// statuses.
type Supervisor struct {
	serverMapMutex *sync.RWMutex
	serverMap      map[string]ServerStatus
	statusUpdates  chan communication.StatusUpdate
	shutdown       chan int
}

// New creates and returns a new *Supervisor.
func New() *Supervisor {
	return &Supervisor{
		serverMapMutex: &sync.RWMutex{},
		serverMap:      map[string]ServerStatus{},
		statusUpdates:  make(chan communication.StatusUpdate, 4),
		shutdown:       make(chan int),
	}
}

// Start starts the Supervisor.
func (s *Supervisor) Start() error {
	go s.loop()
	return nil
}

// Shutdown shuts down the Supervisor.
func (s *Supervisor) Shutdown() error {
	s.shutdown <- 1
	return nil
}

func (s *Supervisor) loop() {
	for {
		select {
		case su := <-s.statusUpdates:
			fmt.Println(su)
		case <-s.shutdown:
			return
		}
	}
}
