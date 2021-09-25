package http

import (
	"github.com/bfoody/Walmart-Scraper/services/hub"
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/supervisor"
	"github.com/go-chi/chi/v5"
)

// Handlers provides an HTTP API.
type Handlers struct {
	service    hub.Service
	supervisor *supervisor.Supervisor
}

// NewHandlers creates and returns a *Handlers.
func NewHandlers(service hub.Service, supervisor *supervisor.Supervisor) *Handlers {
	return &Handlers{
		service,
		supervisor,
	}
}

func (h *Handlers) Router() *chi.Mux {
	r := chi.NewRouter()

	return r
}

func (h *Handlers) 
