package domain

import "time"

// A ProductInfoDiff represents a change in a product's status from one
// ProductInfo to another.
type ProductInfoDiff struct {
	ID                string    `db:"id"`
	CreatedAt         time.Time `db:"created_at"`
	ProductID         string    `db:"product_id"`
	ProductLocationID string    `db:"product_location_id"`
	OldTimestamp      time.Time `db:"old_timestamp"`
	NewTimestamp      time.Time `db:"new_timestamp"`
	OldProductInfoID  string    `db:"old_product_info_id"`
	NewProductInfoID  string    `db:"new_product_info_id"`
	OldInStockValue   bool      `db:"old_in_stock_value"`
	NewInStockValue   bool      `db:"new_in_stock_value"`
}

// An AverageSellout stores the average sellout time for a single
// product.
type AverageSellout struct {
	ID                          string        `db:"id"`
	CreatedAt                   time.Time     `db:"created_at"`
	UpdatedAt                   time.Time     `db:"updated_at"`
	ProductID                   string        `db:"product_id"`
	ProductLocationID           string        `db:"product_location_id"`
	AverageAvailabilityDuration time.Duration `db:"average_availability_duration"`
	AveragedCount               int64         `db:"averaged_count"`
}
