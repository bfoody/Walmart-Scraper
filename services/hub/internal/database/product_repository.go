package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
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
	err := r.db.Get(product, "SELECT * FROM products WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// InsertProduct inserts a single product into the database, returning the ID on success.
func (r *ProductRepository) InsertProduct(product domain.Product) (string, error) {
	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO products (id, common_name) VALUES ($1, $2)", id, product.CommonName)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateProduct updates a single product in the database by ID.
func (r *ProductRepository) UpdateProduct(product domain.Product) error {
	_, err := r.db.Exec("UPDATE products SET common_name = $1 WHERE id = $2", product.CommonName, product.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct deletes a single product by ID.
func (r *ProductRepository) DeleteProduct(id string) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
