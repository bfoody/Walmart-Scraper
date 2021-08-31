package domain

import "time"

// A Product represents a single product, detached from a seller or price.
type Product struct {
	ID         string // the entity's unique ID
	CommonName string // a common name shown for the product throughout the UI
}

// A ProductLocation describes a location where a product is being sold, used to
// track different sellers of a product.
type ProductLocation struct {
	ID         string `db:"id"` // the entity's unique ID
	ProductID  string `db:"product_id"`
	LocationID string `db:"location_id"` // the ID of the product's location
	URL        string `db:"url"`         // the URL of the product on the seller's website
	LocalID    string `db:"local_id"`    // the ID used by the seller for the product
	Slug       string `db:"slug"`        // the slug being use on the seller's website
	Category   string `db:"category"`    // the product's category on the seller's website
}

// A ProductInfo represents a single crawl of a product and the details scraped
// from the crawl.
type ProductInfo struct {
	ID                 string    // the entity's unique ID
	CreatedAt          time.Time // the time at which the info was crawled/logged
	ProductID          string    // the product's ID
	ProductLocationID  string    // the product-location ID
	Price              float32   // the current price of the item in USD
	AvailabilityStatus string    // the availability of the item, eg. "IN_STOCK"
	InStock            bool      // whether or not the item is in stock (AvailabilityStatus but as a boolean)
}
