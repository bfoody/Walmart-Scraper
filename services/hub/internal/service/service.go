package service

import (
	"errors"
	"time"

	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/services/hub"
)

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

// ResolveTask marks the task with the provided ID as completed.
func (s *Service) ResolveTask(id string) error {
	st, err := s.scrapeTaskRepository.FindScrapeTaskByID(id)
	if err != nil {
		return err
	}

	st.Completed = true

	return s.scrapeTaskRepository.UpdateScrapeTask(*st)
}

// SaveProductInfo saves a new ProductInfo to the database, returning the ID on success.
func (s *Service) SaveProductInfo(productInfo domain.ProductInfo) (string, error) {
	productInfo.CreatedAt = time.Now()

	if productInfo.ProductID == "" {
		return "", errors.New("ProductID must not be null")
	}

	if productInfo.ProductLocationID == "" {
		return "", errors.New("ProductLocationID must not be null")
	}

	if productInfo.AvailabilityStatus == "" {
		return "", errors.New("AvailabilityStatus must not be null")
	}

	return s.productInfoRepository.InsertProductInfo(productInfo)
}
