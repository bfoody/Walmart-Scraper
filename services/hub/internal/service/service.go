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
	crawlTaskRepository       hub.CrawlTaskRepository
}

// NewService creates and returns a *Service with the provided dependencies.
func NewService(
	productRepository hub.ProductRepository,
	productInfoRepository hub.ProductInfoRepository,
	productLocationRepository hub.ProductLocationRepository,
	scrapeTaskRepository hub.ScrapeTaskRepository,
	crawlTaskRepository hub.CrawlTaskRepository,
) *Service {
	return &Service{
		productRepository,
		productInfoRepository,
		productLocationRepository,
		scrapeTaskRepository,
		crawlTaskRepository,
	}
}

// CreateTask creates a new task using the provided object.
func (s *Service) CreateTask(scrapeTask domain.ScrapeTask) (string, error) {
	return s.scrapeTaskRepository.InsertScrapeTask(scrapeTask)
}

// FetchUpcomingTasks fetches newest tasks with a limit.
func (s *Service) FetchUpcomingTasks(limit uint16) ([]domain.ScrapeTask, error) {
	return s.scrapeTaskRepository.FindUpcomingScrapeTasks(limit)
}

// GetProductLocationByID gets a single ProductLocation using the ID.
func (s *Service) GetProductLocationByID(id string) (*domain.ProductLocation, error) {
	return s.productLocationRepository.FindProductLocationByID(id)
}

// ResolveTask marks the task with the provided ID as completed.
func (s *Service) ResolveTask(id string, newCallback func(st domain.ScrapeTask)) error {
	st, err := s.scrapeTaskRepository.FindScrapeTaskByID(id)
	if err != nil {
		return err
	}

	st.Completed = true

	err = s.scrapeTaskRepository.UpdateScrapeTask(*st)
	if err != nil {
		return err
	}

	// If repeat is disabled, exit without rescheduling.
	if !st.Repeat {
		return nil
	}

	newSt := domain.ScrapeTask{
		ID:                "",
		Completed:         false,
		CreatedAt:         time.Now(),
		ScheduledFor:      time.Now().Add(st.Interval),
		ProductLocationID: st.ProductLocationID,
		Repeat:            st.Repeat,
		Interval:          st.Interval,
	}

	_, err = s.scrapeTaskRepository.InsertScrapeTask(newSt)
	if err != nil {
		return err
	}

	newCallback(newSt)

	return nil
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

// SaveProduct saves a new Product to the database, returning the ID on success.
func (s *Service) SaveProduct(product domain.Product) (string, error) {
	if product.CommonName == "" {
		return "", errors.New("CommonName must not be null")
	}

	return s.productRepository.InsertProduct(product)
}

// SaveProductLocation saves a ProductLocation to the database.
func (s *Service) SaveProductLocation(productLocation domain.ProductLocation) (string, error) {
	if productLocation.LocalID == "" {
		return "", errors.New("LocalID must not be null")
	}

	if productLocation.LocationID == "" {
		return "", errors.New("LocationID must not be null")
	}

	if productLocation.ProductID == "" {
		return "", errors.New("ProductID must not be null")
	}

	return s.productLocationRepository.InsertProductLocation(productLocation)
}

// IsCrawled returns true if an item was already crawled.
func (s *Service) IsCrawled(productLocationId string) (bool, error) {
	_, err := s.crawlTaskRepository.FindCrawlTaskByProductLocationID(productLocationId)
	return err == nil, nil
}

// SaveCrawlTask saves a crawl task with the provided ID.
func (s *Service) SaveCrawlTask(productLocationId string) (string, error) {
	return s.crawlTaskRepository.InsertCrawlTask(domain.CrawlTask{
		ID:                      "",
		CreatedAt:               time.Now(),
		Completed:               true,
		OriginProductLocationID: productLocationId,
	})
}
