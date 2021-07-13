package database

import (
	"fmt"

	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// A ScrapeTaskRepository provides methods for interacting with ScrapeTasks in
// the database.
type ScrapeTaskRepository struct {
	db *sqlx.DB
}

// NewScrapeTaskRepository creates and returns a *ScrapeTaskRepository from the supplied
// database connection.
func NewScrapeTaskRepository(db *sqlx.DB) *ScrapeTaskRepository {
	return &ScrapeTaskRepository{
		db,
	}
}

// FindScrapeTaskByID finds a single scrape task by ID, returning an error if nothing is found.
func (r *ScrapeTaskRepository) FindScrapeTaskByID(id string) (*domain.ScrapeTask, error) {
	var scrapeTask *domain.ScrapeTask
	err := r.db.Get(scrapeTask, "SELECT * FROM scrape_tasks WHERE id=$1 ORDER BY scheduled_for", id)
	if err != nil {
		return nil, err
	}

	return scrapeTask, nil
}

// FindUpcomingScrapeTasks returns due tasks closest to the current time, using the supplied
// limit.
func (r *ScrapeTaskRepository) FindUpcomingScrapeTasks(limit uint8) ([]domain.ScrapeTask, error) {
	var scrapeTasks []domain.ScrapeTask
	err := r.db.Select(&scrapeTasks, fmt.Sprintf("SELECT * FROM scrape_tasks WHERE completed=FALSE ORDER BY scheduled_for LIMIT %d", limit))
	if err != nil {
		return nil, err
	}

	return scrapeTasks, nil
}

// FindScrapeTasksByProductLocationID finds scrape tasks by ProductLocationID, returning a
// blank array if nothing is found.
func (r *ScrapeTaskRepository) FindScrapeTasksByProductLocationID(id string) ([]domain.ScrapeTask, error) {
	var scrapeTasks []domain.ScrapeTask
	err := r.db.Select(&scrapeTasks, "SELECT * FROM scrape_tasks WHERE product_location_id=$1 ORDER BY scheduled_for", id)
	if err != nil {
		return nil, err
	}

	return scrapeTasks, nil
}

// FindScrapeTasksByProductID finds scrape tasks by Product ID, returning a
// blank array if nothing is found.
func (r *ScrapeTaskRepository) FindScrapeTasksByProductID(id string) ([]domain.ScrapeTask, error) {
	var scrapeTasks []domain.ScrapeTask
	err := r.db.Select(&scrapeTasks, "SELECT * FROM scrape_tasks JOIN product_locations ON scrape_tasks.product_location_id=product_locations.id WHERE product_locations.product_id=$1 ORDER BY scheduled_for", id)
	if err != nil {
		return nil, err
	}

	return scrapeTasks, nil
}

// FindScrapeTasksByLocationID finds scrape tasks by Location ID, returning a
// blank array if nothing is found.
func (r *ScrapeTaskRepository) FindScrapeTasksByLocationID(id string) ([]domain.ScrapeTask, error) {
	var scrapeTasks []domain.ScrapeTask
	err := r.db.Select(&scrapeTasks, "SELECT * FROM scrape_tasks JOIN product_locations ON scrape_tasks.product_location_id=product_locations.id WHERE product_locations.location_id=$1 ORDER BY scheduled_for", id)
	if err != nil {
		return nil, err
	}

	return scrapeTasks, nil
}

// InsertScrapeTask inserts a single scrape task into the database, returning the ID on success.
func (r *ScrapeTaskRepository) InsertScrapeTask(scrapeTask domain.ScrapeTask) (string, error) {
	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO scrape_tasks (id, completed, timestamp, scheduled_for, product_location_id, repeat, interval) VALUES ($1, $2, $3, $4, $5, $6, $7)", id, scrapeTask.Completed, scrapeTask.Timestamp, scrapeTask.ScheduledFor, scrapeTask.ProductLocationID, scrapeTask.Repeat, scrapeTask.Interval)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateScrapeTask updates a single scrape task in the database by ID.
func (r *ScrapeTaskRepository) UpdateScrapeTask(scrapeTask domain.ScrapeTask) error {
	_, err := r.db.Exec("UPDATE scrape_tasks SET completed=$1, timestamp=$2, scheduled_for=$3, product_location_id=$4, repeat=$5, interval=$6 WHERE id=$7", scrapeTask.Completed, scrapeTask.Timestamp, scrapeTask.ScheduledFor, scrapeTask.ProductLocationID, scrapeTask.Repeat, scrapeTask.Interval, scrapeTask.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteScrapeTask deletes a single scrape task by ID.
func (r *ScrapeTaskRepository) DeleteScrapeTask(id string) error {
	_, err := r.db.Exec("DELETE FROM scrape_tasks WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
