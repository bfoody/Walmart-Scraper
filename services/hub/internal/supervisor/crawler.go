package supervisor

import (
	"github.com/bfoody/Walmart-Scraper/services/hub/internal/service"
	"go.uber.org/zap"
)

// A Crawler distributes and collects requests to crawl related products/add new products
// to the scrape list.
type Crawler struct {
	service service.Service
	log     *zap.Logger
}

// NewCrawler creates and returns a *Crawler.
func NewCrawler(service service.Service, logger *zap.Logger) *Crawler {
	return &Crawler{
		service: service,
		log:     logger,
	}
}
