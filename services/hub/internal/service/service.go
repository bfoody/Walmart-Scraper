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
	productInfoDiffRepository hub.ProductInfoDiffRepository
	averageSelloutRepository  hub.AverageSelloutRepository
	productLocationRepository hub.ProductLocationRepository
	scrapeTaskRepository      hub.ScrapeTaskRepository
	crawlTaskRepository       hub.CrawlTaskRepository
}

// NewService creates and returns a *Service with the provided dependencies.
func NewService(
	productRepository hub.ProductRepository,
	productInfoRepository hub.ProductInfoRepository,
	productInfoDiffRepository hub.ProductInfoDiffRepository,
	averageSelloutRepository hub.AverageSelloutRepository,
	productLocationRepository hub.ProductLocationRepository,
	scrapeTaskRepository hub.ScrapeTaskRepository,
	crawlTaskRepository hub.CrawlTaskRepository,
) *Service {
	return &Service{
		productRepository,
		productInfoRepository,
		productInfoDiffRepository,
		averageSelloutRepository,
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

	id, err := s.productInfoRepository.InsertProductInfo(productInfo)
	if err != nil {
		return "", err
	}

	productInfo.ID = id
	go s.runStatistics(productInfo)

	return id, nil
}

// runStatistics calculates and saves all statistic information calculated
// after a ProductInfo being saved.
func (s *Service) runStatistics(productInfo domain.ProductInfo) error {
	lastProductInfos, err := s.productInfoRepository.FindProductInfosByProductID(productInfo.ProductID, 2)
	if err != nil {
		return err
	}

	if len(lastProductInfos) < 2 {
		return errors.New("not enough product infos to calculate diff")
	}

	// Get the ProductInfo before the one that was just saved.
	lastProductInfo := lastProductInfos[1]

	if lastProductInfo.InStock == productInfo.InStock {
		return nil
	}

	pid := domain.ProductInfoDiff{
		ID:                "",
		CreatedAt:         time.Now(),
		ProductID:         productInfo.ProductID,
		ProductLocationID: productInfo.ProductLocationID,
		OldTimestamp:      lastProductInfo.CreatedAt,
		NewTimestamp:      productInfo.CreatedAt,
		OldProductInfoID:  lastProductInfo.ID,
		NewProductInfoID:  productInfo.ID,
	}
	_, err = s.productInfoDiffRepository.InsertProductInfoDiff(pid)
	if err != nil {
		return err
	}

	// Return if the product did not change from in stock to out of stock,
	// nothing more to do in that case.
	if !(lastProductInfo.InStock && !productInfo.InStock) {
		return nil
	}

	lastAverageSellout, err := s.averageSelloutRepository.FindAverageSelloutByProductID(productInfo.ProductID)
	if err != nil {
		return err
	}

	delta := productInfo.CreatedAt.Sub(lastProductInfo.CreatedAt)
	sum := int64(lastAverageSellout.AverageAvailabilityDuration + delta)
	n := lastAverageSellout.AveragedCount + 1
	average := time.Duration(sum / n)

	lastAverageSellout.UpdatedAt = time.Now()
	lastAverageSellout.AverageAvailabilityDuration = average
	lastAverageSellout.AveragedCount = n

	return s.averageSelloutRepository.UpdateAverageSellout(*lastAverageSellout)
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
