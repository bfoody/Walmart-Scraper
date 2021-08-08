package walmart_test

import (
	"testing"

	"github.com/bfoody/Walmart-Scraper/services/client/internal/api"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api/walmart"
)

const (
	testItemName = "onn. 32\" Class HD (720P) Roku Smart LED TV (100012589)"
)

// TestGetItemDetails tests the GetItemDetails scraping method.
func TestGetItemDetails(t *testing.T) {
	c := walmart.NewClient(api.NewHTTPClient())
	item, err := c.GetItemDetails("onn-32-Class-HD-720P-Roku-Smart-LED-TV-100012589", "314022535")
	if err != nil {
		t.Fatal(err)
	}

	if item.Name != testItemName {
		t.Errorf("expected item.Name == %s, got %s", testItemName, item.Name)
	}
}

// TestGetItemRecommendations tests the GetItemRelatedItems scraping method.
func TestGetItemRecommendations(t *testing.T) {
	c := walmart.NewClient(api.NewHTTPClient())
	item, err := c.GetItemDetails("onn-32-Class-HD-720P-Roku-Smart-LED-TV-100012589", "314022535")
	if err != nil {
		t.Fatal(err)
	}

	if item.Name != testItemName {
		t.Errorf("expected item.Name == %s, got %s", testItemName, item.Name)
	}

	items, err := c.GetItemRelatedItems(item.ID, item.CategoryID, item.Category, item.Name)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) < 1 {
		t.Error("`items` array is empty")
	}
}
