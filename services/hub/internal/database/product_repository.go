package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/jmoiron/sqlx"
)

// A ProductRepository provides methods for interacting with Products in the
// database.
type ProductRepository struct {
	db *sqlx.DB
}

// NewProductRepository creates and returns a *ProductRepository with the supplied
// database connection.
func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

// FindProductByID finds a single product by ID, returning an error if nothing is found.
func (r *ProductRepository) FindProductByID(id string) (*domain.Product, error) {
	var product *domain.Product
	err := r.db.Get(&product, "SELECT * FROM products WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
