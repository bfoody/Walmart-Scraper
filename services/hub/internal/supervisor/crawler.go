package supervisor

import (
	"math/rand"
	"time"

	"github.com/bfoody/Walmart-Scraper/communication"
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/services/hub"
	"go.uber.org/zap"
)

// A Crawler distributes and collects requests to crawl related products/add new products
// to the scrape list.
type Crawler struct {
	service      hub.Service
	log          *zap.Logger
	crawledCache map[string]bool // a cache of already crawled items
	callback     func(productLocationID string)
	taskManager  *TaskManager
}

// NewCrawler creates and returns a *Crawler.
func NewCrawler(service hub.Service, logger *zap.Logger, callback func(productLocationID string), taskManager *TaskManager) *Crawler {
	return &Crawler{
		service:      service,
		log:          logger,
		crawledCache: map[string]bool{},
		callback:     callback,
		taskManager:  taskManager,
	}
}

// AttemptCrawl attempts a crawl from a product ID.
func (c *Crawler) AttemptCrawl(id string) {
	if _, ok := c.crawledCache[id]; ok {
		// Already crawled, exit.
		return
	}

	crawled, err := c.service.IsCrawled(id)
	if err != nil {
		c.log.Error("error crawling", zap.Error(err))
		return
	}

	if crawled {
		return
	}

	c.callback(id)
}

// PipeRetrieval receives a crawl.
func (c *Crawler) PipeRetrieval(cr *communication.CrawlRetrieved) {
	_, err := c.service.SaveCrawlTask(cr.ProductLocationID)
	if err != nil {
		c.log.Error("error saving CrawlTask", zap.Error(err))
	}
	c.crawledCache[cr.ProductLocationID] = true

	for _, item := range cr.Recommendations {
		id, err := c.service.SaveProduct(domain.Product{
			ID:         "",
			CommonName: item.Name,
		})
		if err != nil {
			c.log.Error("error saving Product", zap.Error(err))
			continue
		}

		plID, err := c.service.SaveProductLocation(domain.ProductLocation{
			ID:         "",
			Name:       item.Name,
			ProductID:  id,
			LocationID: "8e1922b0-6c12-4bd6-944e-9f87d0b15359",
			URL:        item.URL,
			LocalID:    item.LocalID,
			Slug:       item.Slug,
			CategoryID: item.CategoryID,
			Category:   item.Category,
		})
		if err != nil {
			c.log.Error("error saving ProductLocation", zap.Error(err))
			continue
		}

		st := domain.ScrapeTask{
			ID:                "",
			Completed:         false,
			CreatedAt:         time.Now(),
			ScheduledFor:      time.Now().Add(time.Duration(120+rand.Intn(120)) * time.Second),
			ProductLocationID: plID,
			Repeat:            true,
			Interval:          120 * time.Second,
		}

		stID, err := c.service.CreateTask(st)
		if err != nil {
			c.log.Error("error saving ScrapeTask", zap.Error(err))
			continue
		}

		st.ID = stID

		c.taskManager.pushTaskToQueue(st)
	}
}
