package database

import (
	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// A ProductInfoRepository provides methods for interfacing with ProductInfos
// stored in the database.
type ProductInfoRepository struct {
	db *sqlx.DB
}

// NewProductInfoRepository creates and returns a *ProductInfoRepository with the supplied database connection.
func NewProductInfoRepository(db *sqlx.DB) *ProductInfoRepository {
	return &ProductInfoRepository{
		db,
	}
}

// FindProductInfoByID finds a single product info by ID, returning an error if nothing is found.
func (r *ProductInfoRepository) FindProductInfoByID(id string) (*domain.ProductInfo, error) {
	var productInfo *domain.ProductInfo
	err := r.db.Get(productInfo, "SELECT * FROM product_infos WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return productInfo, nil
}

// InsertProductInfo inserts a single product into the database, returning the ID on success.
func (r *ProductInfoRepository) InsertProductInfo(productInfo domain.ProductInfo) (string, error) {
	id := uuid.Generate()
	_, err := r.db.Exec("INSERT INTO product_infos (id, timestamp, product_id, product_location_id, price, availability_status, in_stock) VALUES ($1, $2, $3, $4, $5, $6, $7)", id, productInfo.Timestamp, productInfo.ProductID, productInfo.ProductLocationID, productInfo.Price, productInfo.AvailabilityStatus, productInfo.InStock)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateProductInfo updates a single product info in the database by ID.
func (r *ProductInfoRepository) UpdateProductInfo(productInfo domain.ProductInfo) error {
	_, err := r.db.Exec("UPDATE product_infos SET timestamp=$1, product_id=$2, product_location_id=$3, price=$4, availability_status=$5, in_stock=$6 WHERE id=$7", productInfo.Timestamp, productInfo.ProductID, productInfo.ProductLocationID, productInfo.Price, productInfo.AvailabilityStatus, productInfo.InStock, productInfo.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProductInfo deletes a single product info by ID.
func (r *ProductInfoRepository) DeleteProductInfo(id string) error {
	_, err := r.db.Exec("DELETE FROM product_infos WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
