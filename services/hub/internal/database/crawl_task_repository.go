package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// A CrawlTaskRepository provides methods for interacting with CrawlTasks in
// the database.
type CrawlTaskRepository struct {
	db *sqlx.DB
}

// NewCrawlTaskRepository creates and returns a *CrawlTaskRepository from the supplied
// database connection.
func NewCrawlTaskRepository(db *sqlx.DB) *CrawlTaskRepository {
	return &CrawlTaskRepository{
		db,
	}
}

// FindCrawlTaskByID finds a single crawl task by ID, returning an error if nothing is found.
func (r *CrawlTaskRepository) FindCrawlTaskByID(id string) (*domain.CrawlTask, error) {
	crawlTask := &domain.CrawlTask{}
	err := r.db.Get(crawlTask, "SELECT * FROM crawl_tasks WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return crawlTask, nil
}

// FindCrawlTaskByProductLocationID finds a single crawl task by ProductLocation ID, returning an error if nothing is found.
func (r *CrawlTaskRepository) FindCrawlTaskByProductLocationID(id string) (*domain.CrawlTask, error) {
	crawlTask := &domain.CrawlTask{}
	err := r.db.Get(crawlTask, "SELECT * FROM crawl_tasks WHERE origin_product_location_id=$1", id)
	if err != nil {
		return nil, err
	}

	return crawlTask, nil
}

// InsertCrawlTask inserts a single crawl task into the database, returning the ID on success.
func (r *CrawlTaskRepository) InsertCrawlTask(crawlTask domain.CrawlTask) (string, error) {

	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO crawl_tasks (id, completed, created_at, origin_product_location_id) VALUES ($1, $2, $3, $4)", id, crawlTask.Completed, crawlTask.CreatedAt, crawlTask.OriginProductLocationID)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateCrawlTask updates a single crawl task in the database by ID.
func (r *CrawlTaskRepository) UpdateCrawlTask(crawlTask domain.CrawlTask) error {
	_, err := r.db.Exec("UPDATE crawl_tasks SET completed=$1, created_at=$2, origin_product_location_id=$3 WHERE ID=$4", crawlTask.Completed, crawlTask.CreatedAt, crawlTask.OriginProductLocationID, crawlTask.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCrawlTask deletes a single crawl task by ID.
func (r *CrawlTaskRepository) DeleteCrawlTask(id string) error {
	_, err := r.db.Exec("DELETE FROM scrape_tasks WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
