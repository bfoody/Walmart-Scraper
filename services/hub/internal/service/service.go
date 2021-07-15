package service

import "github.com/bfoody/Walmart-Scraper/services/hub"

// A Service handles storage and retrieval of Product information, as well as tasks.
type Service struct {
	productRepository         hub.ProductRepository
	productInfoRepository     hub.ProductInfoRepository
	productLocationRepository hub.ProductLocationRepository
	scrapeTaskRepository      hub.ScrapeTaskRepository
}

// NewService creates and returns a *Service with the provided dependencies.
func NewService(
	productRepository hub.ProductRepository,
	productInfoRepository hub.ProductInfoRepository,
	productLocationRepository hub.ProductLocationRepository,
	scrapeTaskRepository hub.ScrapeTaskRepository,
) *Service {
	return &Service{
		productRepository,
		productInfoRepository,
		productLocationRepository,
		scrapeTaskRepository,
	}
}
