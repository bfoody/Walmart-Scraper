package supervisor

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/bfoody/Walmart-Scraper/communication"
	"go.uber.org/zap"
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
	log            *zap.Logger
}

// New creates and returns a new *Supervisor.
func New(logger *zap.Logger) *Supervisor {
	return &Supervisor{
		serverMapMutex: &sync.RWMutex{},
		serverMap:      map[string]ServerStatus{},
		statusUpdates:  make(chan communication.StatusUpdate, 4),
		shutdown:       make(chan int),
		log:            logger,
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

// PipeStatusUpdate pipes a StatusUpdate into the supervisor.
func (s *Supervisor) PipeStatusUpdate(su communication.StatusUpdate) {
	s.statusUpdates <- su
}

func (s *Supervisor) loop() {
	for {
		select {
		case su := <-s.statusUpdates:
			s.handleStatusUpdate(&su)
		case <-s.shutdown:
			return
		}
	}
}

func (s *Supervisor) handleStatusUpdate(su *communication.StatusUpdate) {
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

	// Replace the status with the new one.
	s.serverMap[su.SenderID] = status
}
