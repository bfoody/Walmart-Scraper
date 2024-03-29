package hub

import "github.com/bfoody/Walmart-Scraper/domain"

// A ProductRepository provides methods for interfacing with Products stored
// in the database.
type ProductRepository interface {
	// FindProductByID finds a single product by ID, returning an error if nothing is found.
	FindProductByID(id string) (*domain.Product, error)
	// InsertProduct inserts a single product into the database, returning the ID on success.
	InsertProduct(product domain.Product) (string, error)
	// UpdateProduct updates a single product in the database by ID.
	UpdateProduct(product domain.Product) error
	// DeleteProduct deletes a single product by ID.
	DeleteProduct(id string) error
}

// A ProductLocationRepository provides methods for interfacing with ProductLocations
// stored in the database.
type ProductLocationRepository interface {
	// FindProductLocationByID finds a single product location by ID, returning an error
	// if nothing is found.
	FindProductLocationByID(id string) (*domain.ProductLocation, error)
	// FindProductLocationsByProductID finds multiple product locations by product ID,
	// returning an empty array if nothing is found.
	FindProductLocationsByProductID(id string) ([]domain.ProductLocation, error)
	// FindProductLocationsByLocationID finds multiple product locations by location ID,
	// returning an empty array if nothing is found.
	FindProductLocationsByLocationID(id string) ([]domain.ProductLocation, error)
	// FindProductLocationByProductAndLocationID finds a single product location by both a product
	// and location ID, returning an error if nothing is found.
	FindProductLocationByProductAndLocationID(productID, locationID string) (*domain.ProductLocation, error)
	// InsertProductLocation inserts a single product location into the database,
	// returning the ID on success.
	InsertProductLocation(productLocation domain.ProductLocation) (string, error)
	// UpdateProductLocation updates a single product location in the database by ID.
	UpdateProductLocation(productLocation domain.ProductLocation) error
	// DeleteProductLocation deletes a single product location by ID.
	DeleteProductLocation(id string) error
}

// A ProductInfoRepository provides methods for interfacing with ProductInfos
// stored in the database.
type ProductInfoRepository interface {
	// FindProductInfoByID finds a single product info by ID, returning an error if nothing is found.
	FindProductInfoByID(id string) (*domain.ProductInfo, error)
	// InsertProductInfo inserts a single product into the database, returning the ID on success.
	InsertProductInfo(productInfo domain.ProductInfo) (string, error)
	// UpdateProductInfo updates a single product info in the database by ID.
	UpdateProductInfo(productInfo domain.ProductInfo) error
	// DeleteProductInfo deletes a single product info by ID.
	DeleteProductInfo(id string) error
}

// A ScrapeTaskRepository provides methods for interfacing with ScrapeTasks
// stored in the database.
type ScrapeTaskRepository interface {
	// FindScrapeTaskByID finds a single scrape task by ID, returning an error if nothing is found.
	FindScrapeTaskByID(id string) (*domain.ScrapeTask, error)
	// FindUpcomingScrapeTasks returns due tasks closest to the current time, using the supplied
	// limit.
	FindUpcomingScrapeTasks(limit uint16) ([]domain.ScrapeTask, error)
	// FindScrapeTasksByProductLocationID finds scrape tasks by ProductLocationID, returning a
	// blank array if nothing is found.
	FindScrapeTasksByProductLocationID(id string) ([]domain.ScrapeTask, error)
	// FindScrapeTasksByProductID finds scrape tasks by Product ID, returning a
	// blank array if nothing is found.
	FindScrapeTasksByProductID(id string) ([]domain.ScrapeTask, error)
	// FindScrapeTasksByLocationID finds scrape tasks by Location ID, returning a
	// blank array if nothing is found.
	FindScrapeTasksByLocationID(id string) ([]domain.ScrapeTask, error)
	// InsertScrapeTask inserts a single scrape task into the database, returning the ID on success.
	InsertScrapeTask(scrapeTask domain.ScrapeTask) (string, error)
	// UpdateScrapeTask updates a single scrape task in the database by ID.
	UpdateScrapeTask(scrapeTask domain.ScrapeTask) error
	// DeleteScrapeTask deletes a single scrape task by ID.
	DeleteScrapeTask(id string) error
}

// A CrawlTaskRepository provides methods for interfacing with CrawlTasks
// stored in the database.
type CrawlTaskRepository interface {
	// FindCrawlTaskByID finds a single crawl task by ID, returning an error if nothing is found.
	FindCrawlTaskByID(id string) (*domain.CrawlTask, error)
	// FindCrawlTaskByID finds a single crawl task by ProductLocation ID, returning an error if nothing is found.
	FindCrawlTaskByProductLocationID(id string) (*domain.CrawlTask, error)
	// InsertCrawlTask inserts a single crawl task into the database, returning the ID on success.
	InsertCrawlTask(crawlTask domain.CrawlTask) (string, error)
	// UpdateCrawlTask updates a single crawl task in the database by ID.
	UpdateCrawlTask(crawlTask domain.CrawlTask) error
	// DeleteCrawlTask deletes a single crawl task by ID.
	DeleteCrawlTask(id string) error
}

// A Service provides abstractions for interacting with product and task data in the database.
type Service interface {
	// ResolveTask marks the task with the provided ID as completed.
	ResolveTask(id string, newCallback func(st domain.ScrapeTask)) error
	// SaveProductInfo saves a new ProductInfo to the database, returning the ID on success.
	SaveProductInfo(productInfo domain.ProductInfo) (string, error)
	// SaveProduct saves a new Product to the database, returning the ID on success.
	SaveProduct(product domain.Product) (string, error)
	// CreateTask creates a new task using the provided object.
	CreateTask(scrapeTask domain.ScrapeTask) (string, error)
	// FetchUpcomingTasks fetches newest tasks with a limit.
	FetchUpcomingTasks(limit uint16) ([]domain.ScrapeTask, error)
	// GetProductLocationByID gets a single ProductLocation using the ID.
	GetProductLocationByID(id string) (*domain.ProductLocation, error)
	// SaveProductLocation saves a ProductLocation to the database.
	SaveProductLocation(productLocation domain.ProductLocation) (string, error)
	// IsCrawled returns true if an item was already crawled.
	IsCrawled(productLocationId string) (bool, error)
	// SaveCrawlTask saves a crawl task with the provided ID.
	SaveCrawlTask(productLocationId string) (string, error)
}
