package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	httpConfig "github.com/bfoody/Walmart-Scraper/services/client/internal/http"
)

// HTTPClient wraps an http client and provides convenience methods for sending
// GET and POST requests with simulated browser user agents.
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates and returns a new HTTPClient with a pre-configured
// http.Client.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: httpConfig.NewClient(),
	}
}

// User agent of Chrome 89 running on macOS 11.3.
var simulatedHeaders = map[string]string{
	"authority":                         "www.walmart.com",
	"cache-control":                     "max-age=0",
	"accept":                            "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"service-worker-navigation-preload": "true",
	"User-Agent":                        "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_3_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
	"sec-ch-ua":                         "\"Google Chrome\";v=\"89\", \"Chromium\";v=\"89\", \";Not A Brand\";v=\"99\"",
	"sec-ch-ua-mobile":                  "?0",
	"sec-fetch-dest":                    "document",
	"sec-fetch-mode":                    "navigate",
	"sec-fetch-site":                    "none",
	"sec-fetch-user":                    "?1",
}

// addSimulatedHeaders adds various HTTP headers to a request in order to simulate
// a human visiting the site with a browser. This lowers the chance of CAPTCHA traps
// and IP blocks.
func addSimulatedHeaders(req *http.Request) {
	// Iterate simulated headers map and set each header.
	for key, value := range simulatedHeaders {
		req.Header.Set(key, value)
	}
}

// SetProxy sets the HTTP client's proxy using a URL.
func (c *HTTPClient) SetProxy(urlStr string) error {
	if !strings.Contains(urlStr, "://") {
		urlStr = "http://" + urlStr
	}

	proxy, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}

	return nil
}

// Get sends an HTTP GET request to the specified URL, returning an
// *http.Response and an *HTTPError wrapping the http error on failure.
func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, IntoHTTPError(err)
	}

	// Add simulated headers to simulate human browsing.
	addSimulatedHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return resp, IntoHTTPError(err)
	}

	return resp, nil
}

// Post sends an HTTP POST request to the specified URL with the specified body, returning an
// *http.Response and an *HTTPError wrapping the http error on failure.
func (c *HTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	resp, err := c.client.Post(url, contentType, body)
	if err != nil {
		return resp, IntoHTTPError(err)
	}

	return resp, nil
}

// PostJSON sends an HTTP POST request to the specified URL, marshalling the body into JSON, and returning
// an *http.Response and an *HTTPError wrapping the http error on failure.
func (c *HTTPClient) PostJSON(url string, body interface{}) (*http.Response, error) {
	// TODO: better error handling here
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode(body)

	resp, err := c.client.Post(url, "application/json", b)
	if err != nil {
		return resp, IntoHTTPError(err)
	}

	return resp, nil
}
