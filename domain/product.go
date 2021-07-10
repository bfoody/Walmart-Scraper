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
	ID         string // the entity's unique ID
	LocationID string // the ID of the product's location
	URL        string // the URL of the product on the seller's website
	Slug       string // the slug being use on the seller's website
	Category   string // the product's category on the seller's website
}

// A ProductInfo represents a single crawl of a product and the details scraped
// from the crawl.
type ProductInfo struct {
	ID                 string    // the entity's unique ID
	Timestamp          time.Time // the time at which the info was crawled/logged
	ProductID          string    // the product's ID
	ProductLocationID  string    // the product-location ID
	Price              float32   // the current price of the item in USD
	AvailabilityStatus string    // the availability of the item, eg. "IN_STOCK"
	InStock            bool      // whether or not the item is in stock (AvailabilityStatus but as a boolean)
}
