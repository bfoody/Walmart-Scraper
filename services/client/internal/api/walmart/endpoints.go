package walmart

import "fmt"

// WalmartAPIBase is the base URL of the Walmart site
// to which endpoints are appended.
const WalmartAPIBase = "https://walmart.com"

var (
	// ItemDetailsPage returns the endpoint for viewing a single item's
	// details.
	ItemDetailsPage = func(itemSlug string, itemID string) string {
		return fmt.Sprintf("%s/ip/%s/%s", WalmartAPIBase, itemSlug, itemID)
	}
)
