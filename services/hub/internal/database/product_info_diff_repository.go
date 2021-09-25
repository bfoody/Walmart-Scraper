package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// A ProductInfoDiffRepository provides methods for interacting with ProductInfoDiffs in
// the database.
type ProductInfoDiffRepository struct {
	db *sqlx.DB
}

// NewProductInfoDiffRepository creates and returns a *ProductInfoDiffRepository from the supplied
// database connection.
func NewProductInfoDiffRepository(db *sqlx.DB) *ProductInfoDiffRepository {
	return &ProductInfoDiffRepository{
		db,
	}
}

// FindProductInfoDiffByID finds a single ProductInfoDiff by ID, returning an error if nothing is found.
func (r *ProductInfoDiffRepository) FindProductInfoDiffByID(id string) (*domain.ProductInfoDiff, error) {
	productInfoDiff := &domain.ProductInfoDiff{}
	err := r.db.Get(productInfoDiff, "SELECT * FROM product_info_diffs WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return productInfoDiff, nil
}

// FindProductInfoDiffsByProductID finds ProductInfoDiffs by Product ID, returning
// a blank array if nothing is found.
func (r *ProductInfoDiffRepository) FindProductInfoDiffsByProductID(id string) ([]domain.ProductInfoDiff, error) {
	var productInfoDiffs []domain.ProductInfoDiff
	err := r.db.Select(&productInfoDiffs, "SELECT * FROM product_info_diffs WHERE product_id=$1", id)
	if err != nil {
		return nil, err
	}

	return productInfoDiffs, nil
}

// FindProductInfoDiffsByProductLocationID finds ProductInfoDiffs by ProductLocation ID, returning
// a blank array if nothing is found.
func (r *ProductInfoDiffRepository) FindProductInfoDiffsByProductLocationID(id string) ([]domain.ProductInfoDiff, error) {
	var productInfoDiffs []domain.ProductInfoDiff
	err := r.db.Select(&productInfoDiffs, "SELECT * FROM product_info_diffs WHERE product_location_id=$1", id)
	if err != nil {
		return nil, err
	}

	return productInfoDiffs, nil
}

// InsertProductInfoDiff inserts a ProductInfoDiff into the database, returning the ID on success.
func (r *ProductInfoDiffRepository) InsertProductInfoDiff(productInfoDiff domain.ProductInfoDiff) (string, error) {
	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO product_info_diffs (id, created_at, product_id, product_location_id, old_timestamp, new_timestamp, old_product_info_id, new_product_info_id, old_in_stock_value, new_in_stock_value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", id, productInfoDiff.CreatedAt, productInfoDiff.ProductID, productInfoDiff.ProductLocationID, productInfoDiff.OldTimestamp, productInfoDiff.NewTimestamp, productInfoDiff.OldProductInfoID, productInfoDiff.NewProductInfoID, productInfoDiff.OldInStockValue, productInfoDiff.NewInStockValue)
	if err != nil {
		return "", err
	}
	return id, nil
}

// UpdateProductInfoDiff updates a ProductInfoDiff in the database by ID.
func (r *ProductInfoDiffRepository) UpdateProductInfoDiff(productInfoDiff domain.ProductInfoDiff) error {
	_, err := r.db.Exec("UPDATE product_info_diffs SET created_at=$1, product_id=$2, product_location_id=$3, old_timestamp=$4, new_timestamp=$5, old_product_info_id=$6, new_product_info_id=$7, old_in_stock_value=$8, new_in_stock_value=$9 WHERE id=$10", productInfoDiff.CreatedAt, productInfoDiff.ProductID, productInfoDiff.ProductLocationID, productInfoDiff.OldTimestamp, productInfoDiff.NewTimestamp, productInfoDiff.OldProductInfoID, productInfoDiff.NewProductInfoID, productInfoDiff.OldInStockValue, productInfoDiff.NewInStockValue, productInfoDiff.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProductInfoDiff deletes a ProductInfoDiff by ID.
func (r *ProductInfoDiffRepository) DeleteProductInfoDiff(id string) error {
	_, err := r.db.Exec("DELETE FROM product_info_diffs WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
