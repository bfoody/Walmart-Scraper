package receiver

import (
	"time"

	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api/walmart"
)

// itemDetailsToProductInfo converts an ItemDetails to a ProductInfo.
func itemDetailsToProductInfo(productLocationID string, productID string, id walmart.ItemDetails) domain.ProductInfo {
	return domain.ProductInfo{
		ID:                 "",         // will be filled in by database service
		CreatedAt:          time.Now(), // will be filled in by database service
		ProductID:          productID,
		ProductLocationID:  productLocationID,
		Price:              id.Price,
		AvailabilityStatus: id.AvailabilityStatus,
		InStock:            id.InStock,
	}
}

// itemDetailsToProductLocation converts an ItemDetails to a ProductLocation..
func itemDetailsToProductLocation(id walmart.ItemDetails) domain.ProductLocation {
	return domain.ProductLocation{
		ID:         "",
		Name:       id.Name,
		ProductID:  "",
		LocationID: "",
		URL:        walmart.ItemDetailsPage(id.Slug, id.ID),
		LocalID:    id.ID,
		Slug:       id.Slug,
		CategoryID: id.CategoryID,
		Category:   id.Category,
	}
}
