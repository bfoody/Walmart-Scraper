package walmart

import (
	"fmt"
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
func (c *Client) GetItemDetails(itemSlug, itemID string) error {
	resp, err := c.client.Get(ItemDetailsPage(itemSlug, itemID))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	html := string(body)

	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return err
	}

	name := htmlquery.FindOne(doc, "//h1")

	fmt.Println(html)
	fmt.Println(htmlquery.InnerText(name))

	return nil
}
