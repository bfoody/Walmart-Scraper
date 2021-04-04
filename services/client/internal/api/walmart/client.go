package walmart

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api"
)

type Client struct {
	client *api.HTTPClient
}

// NewClient creates and returns a new Walmart API Client.
func NewClient(client *api.HTTPClient) *Client {
	return &Client{
		client,
	}
}

// GetItemDetails scrapes the details page for a single item.
func (c *Client) GetItemDetails(itemSlug, itemID string) (*ItemDetails, error) {
	// Fetch the item page.
	resp, err := c.client.Get(ItemDetailsPage(itemSlug, itemID))
	if err != nil {
		// Return the HTTPError.
		return nil, err
	}

	// Read the HTML body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, api.NewAPIError(resp, "failed to read response body", "io_error", err)
	}

	// Parse the HTML with htmlquery.
	html := string(body)
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return nil, api.NewAPIError(resp, "failed to parse html", "html_parse_error", err)
	}

	// Match the `item` JSON in the webpage with this type.
	type queryType struct {
		Item struct {
			Product struct {
				MidasContext struct {
					Brand                       string  `json:"brand"`
					CategoryPathID              string  `json:"categoryPathId"`
					CategoryPathName            string  `json:"categoryPathName"`
					FreeShipping                bool    `json:"freeShipping"`
					InStore                     bool    `json:"inStore"`
					IsTwoDayDeliveryTextEnabled bool    `json:"isTwoDayDeliveryTextEnabled"`
					ItemID                      string  `json:"itemId"`
					Manufacturer                string  `json:"manufacturer"`
					Online                      bool    `json:"online"`
					PageType                    string  `json:"item"`
					Preorder                    bool    `json:"preorder"`
					Price                       float32 `json:"price"`
					Query                       string  `json:"query"`
				} `json:"midasContext"`
				BuyBox struct {
					Products []struct {
						AvailabilityStatus string `json:"availabilityStatus"`
						PickUpMethod       string `json:"pickUpMethod"`
					} `json:"products"`
				} `json:"buyBox"`
			} `json:"product"`
			Query string `json:"query"`
		}
	}

	// Find the JSON payload inside the script tag with the ID of item.
	name := htmlquery.FindOne(doc, "//script[@id=\"item\"]")
	jsonStr := htmlquery.InnerText(name)

	// Decode the JSON value into the struct.
	query := queryType{}
	err = json.NewDecoder(strings.NewReader(jsonStr)).Decode(&query)
	if err != nil {
		return nil, api.NewAPIError(resp, "failed to decode item json payload", "deserialization_error", err)
	}

	if len(query.Item.Product.BuyBox.Products) < 1 {
		return nil, api.NewAPIError(resp, "item.product.buyBox.products array empty", "parse_error", nil)
	}

	return &ItemDetails{
		ID:                 itemID,
		Slug:               itemSlug,
		Name:               query.Item.Product.MidasContext.Query,
		Category:           query.Item.Product.MidasContext.CategoryPathName,
		Price:              query.Item.Product.MidasContext.Price,
		AvailabilityStatus: query.Item.Product.BuyBox.Products[0].AvailabilityStatus,
		InStock:            query.Item.Product.BuyBox.Products[0].AvailabilityStatus == "IN_STOCK",
	}, nil
}
