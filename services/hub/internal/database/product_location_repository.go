package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// A ProductLocationRepository provides methods for interacting with ProductLocations in the
// database.
type ProductLocationRepository struct {
	db *sqlx.DB
}

// NewProductLocationRepository creates and returns a *ProductLocationRepository with the supplied
// database connection.
func NewProductLocationRepository(db *sqlx.DB) *ProductLocationRepository {
	return &ProductLocationRepository{db}
}

// FindProductLocationByID finds a single product location by ID, returning an error
// if nothing is found.
func (r *ProductLocationRepository) FindProductLocationByID(id string) (*domain.ProductLocation, error) {
	productLocation := &domain.ProductLocation{}
	err := r.db.Get(productLocation, "SELECT * FROM product_locations WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return productLocation, nil
}

// FindProductLocationsByProductID finds multiple product locations by product ID,
// returning an empty array if nothing is found.
func (r *ProductLocationRepository) FindProductLocationsByProductID(id string) ([]domain.ProductLocation, error) {
	var productLocations []domain.ProductLocation
	err := r.db.Select(productLocations, "SELECT * FROM product_locations WHERE product_id=$1", id)
	if err != nil {
		return nil, err
	}

	return productLocations, nil
}

// FindProductLocationsByLocationID finds multiple product locations by location ID,
// returning an empty array if nothing is found.
func (r *ProductLocationRepository) FindProductLocationsByLocationID(id string) ([]domain.ProductLocation, error) {
	var productLocations []domain.ProductLocation
	err := r.db.Select(productLocations, "SELECT * FROM product_locations WHERE location_id=$1", id)
	if err != nil {
		return nil, err
	}

	return productLocations, nil
}

// FindProductLocationByProductAndLocationID finds a single product location by both a product
// and location ID, returning an error if nothing is found.
func (r *ProductLocationRepository) FindProductLocationByProductAndLocationID(productID, locationID string) (*domain.ProductLocation, error) {
	var productLocation *domain.ProductLocation
	err := r.db.Get(productLocation, "SELECT * FROM product_locations WHERE product_id=$1 AND location_id=$2", productID, locationID)
	if err != nil {
		return nil, err
	}

	return productLocation, nil
}

// InsertProductLocation inserts a single product location into the database,
// returning the ID on success.
func (r *ProductLocationRepository) InsertProductLocation(productLocation domain.ProductLocation) (string, error) {
	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO product_locations (id, location_id, url, slug, category) VALUES ($1, $2, $3, $4, $5)", id, productLocation.LocationID, productLocation.URL, productLocation.Slug, productLocation.Category)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateProductLocation updates a single product location in the database by ID.
func (r *ProductLocationRepository) UpdateProductLocation(productLocation domain.ProductLocation) error {
	_, err := r.db.Exec("UPDATE product_locations SET location_id=$1, url=$2, slug=$3, category=$4 WHERE id=$5", productLocation.LocationID, productLocation.URL, productLocation.Slug, productLocation.Category, productLocation.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProductLocation deletes a single product location by ID.
func (r *ProductLocationRepository) DeleteProductLocation(id string) error {
	_, err := r.db.Exec("DELETE FROM product_locations WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
