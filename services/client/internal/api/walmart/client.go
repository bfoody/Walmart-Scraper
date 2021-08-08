package walmart

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
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
		CategoryID:         query.Item.Product.MidasContext.CategoryPathID,
		Price:              query.Item.Product.MidasContext.Price,
		AvailabilityStatus: query.Item.Product.BuyBox.Products[0].AvailabilityStatus,
		InStock:            query.Item.Product.BuyBox.Products[0].AvailabilityStatus == "IN_STOCK",
	}, nil
}

var SlugRegex = regexp.MustCompile("/\\/ip\\/(.{1,})\\//g")

// GetItemRelatedItems scrapes related items for a specific item, returning an array of items.
func (c *Client) GetItemRelatedItems(itemID, categoryID, categoryPath, itemName string) ([]ItemDetails, error) {
	// Fetch the item page.
	resp, err := c.client.Get(ItemRecommendations(itemID, categoryID, categoryPath, itemName))
	if err != nil {
		// Return the HTTPError.
		return nil, err
	}

	// Read the HTML body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, api.NewAPIError(resp, "failed to read response body", "io_error", err)
	}

	jsonStr := string(body)

	type queryType struct {
		StatusCode int // should be "200"
		Modules    []struct {
			Name    string
			Configs map[string]struct {
				Title    string
				Products []struct {
					ID struct {
						ProductID string `json:"productId"`
					}
					Price struct {
						CurrentPrice float32
					}
					ProductName        string
					ProductURL         string `json:"productUrl"`
					Category           string
					AvailabilityStatus string
				}
			}
		}
	}

	// Decode the JSON value into the struct.
	query := queryType{}
	err = json.NewDecoder(strings.NewReader(jsonStr)).Decode(&query)
	if err != nil {
		return nil, api.NewAPIError(resp, "failed to decode item json payload", "deserialization_error", err)
	}

	if query.StatusCode != 200 {
		return nil, api.NewAPIError(resp, "API returned non-200 status code", "server_error", nil)
	}

	items := map[string]ItemDetails{}

	for _, module := range query.Modules {
		for _, config := range module.Configs {
			for _, product := range config.Products {
				if product.ID.ProductID == "" {
					continue
				}

				slug := SlugRegex.FindString(product.ProductURL)

				items[product.ID.ProductID] = ItemDetails{
					ID:                 product.ID.ProductID,
					Slug:               slug,
					Name:               product.ProductName,
					Category:           product.Category,
					CategoryID:         "",
					Price:              product.Price.CurrentPrice,
					AvailabilityStatus: product.AvailabilityStatus,
					InStock:            product.AvailabilityStatus == "IN_STOCK",
				}
			}
		}
	}

	itemsArray := []ItemDetails{}
	for _, item := range items {
		itemsArray = append(itemsArray, item)
	}

	return itemsArray, nil
}
