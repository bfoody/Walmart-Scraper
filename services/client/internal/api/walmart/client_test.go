package walmart_test

import (
	"testing"

	"github.com/bfoody/Walmart-Scraper/services/client/internal/api"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api/walmart"
)

// TestGetItemDetails tests the GetItemDetails scraping method.
func TestGetItemDetails(t *testing.T) {
	c := walmart.NewClient(api.NewHTTPClient())
	c.GetItemDetails("onn-32-Class-HD-720P-Roku-Smart-LED-TV-100012589", "314022535")
	t.Error()
}
