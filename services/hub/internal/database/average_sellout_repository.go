package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// A AverageSelloutRepository provides methods for interacting with AverageSellouts in
// the database.
type AverageSelloutRepository struct {
	db *sqlx.DB
}

// NewAverageSelloutRepository creates and returns a *AverageSelloutRepository from the supplied
// database connection.
func NewAverageSelloutRepository(db *sqlx.DB) *AverageSelloutRepository {
	return &AverageSelloutRepository{
		db,
	}
}

// FindAverageSelloutByID finds a single AverageSellout by ID, returning an error if nothing is found.
func (r *AverageSelloutRepository) FindAverageSelloutByID(id string) (*domain.AverageSellout, error) {
	averageSellout := &domain.AverageSellout{}
	err := r.db.Get(averageSellout, "SELECT * FROM average_sellouts WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return averageSellout, nil
}

// FindAverageSelloutByProductID finds a single AverageSellout by Product ID, returning an error if nothing
// is found.
func (r *AverageSelloutRepository) FindAverageSelloutByProductID(id string) (*domain.AverageSellout, error) {
	averageSellout := &domain.AverageSellout{}
	err := r.db.Get(averageSellout, "SELECT * FROM average_sellouts WHERE product_id=$1", id)
	if err != nil {
		return nil, err
	}

	return averageSellout, nil
}

// FindAverageSelloutsByProductLocationID finds AverageSellouts by ProductLocation ID, returning
// a blank array if nothing is found.
func (r *AverageSelloutRepository) FindAverageSelloutsByProductLocationID(id string) ([]domain.AverageSellout, error) {
	var averageSellouts []domain.AverageSellout
	err := r.db.Select(&averageSellouts, "SELECT * FROM average_sellouts WHERE product_location_id=$1", id)
	if err != nil {
		return nil, err
	}

	return averageSellouts, nil
}

// InsertAverageSellout inserts a AverageSellout into the database, returning the ID on success.
func (r *AverageSelloutRepository) InsertAverageSellout(averageSellout domain.AverageSellout) (string, error) {
	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO average_sellouts (id, created_at, updated_at, product_id, product_location_id, average_availability_duration, averaged_count) VALUES ($1, $2, $3, $4, $5, $6, $7)", id, averageSellout.CreatedAt, averageSellout.UpdatedAt, averageSellout.ProductID, averageSellout.ProductLocationID, averageSellout.AverageAvailabilityDuration, averageSellout.AveragedCount)
	if err != nil {
		return "", err
	}
	return id, nil
}

// UpdateAverageSellout updates a AverageSellout in the database by ID.
func (r *AverageSelloutRepository) UpdateAverageSellout(averageSellout domain.AverageSellout) error {
	_, err := r.db.Exec("UPDATE average_sellouts SET created_at=$1, updated_at=$2, product_id=$3, product_location_id=$4, average_availability_duration=$5, averaged_count=$6 WHERE id=$7", averageSellout.CreatedAt, averageSellout.UpdatedAt, averageSellout.ProductID, averageSellout.ProductLocationID, averageSellout.AveragedCount)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAverageSellout deletes a AverageSellout by ID.
func (r *AverageSelloutRepository) DeleteAverageSellout(id string) error {
	_, err := r.db.Exec("DELETE FROM average_sellouts WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
