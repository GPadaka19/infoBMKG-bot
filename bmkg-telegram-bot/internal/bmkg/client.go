package bmkg

import (
	"bmkg-bot/internal/model"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	RSSUrlID = "https://www.bmkg.go.id/alerts/nowcast/id/rss.xml"
	RSSUrlEN = "https://www.bmkg.go.id/alerts/nowcast/en/rss.xml"
)

type Client struct {
	http *resty.Client
}

func NewClient() *Client {
	client := resty.New()

	// Best Practice 1: Set Timeout
	client.SetTimeout(15 * time.Second)

	// Best Practice 2: Auto Retry
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	// Optional: Debug mode (prints requests/responses)
	// client.EnableTrace()

	return &Client{
		http: client,
	}
}

// FetchRSS retrieves and parses the RSS feed from BMKG
func (c *Client) FetchRSS(lang string) (*model.RSSFeed, error) {
	url := RSSUrlID
	if lang == "en" {
		url = RSSUrlEN
	}

	// Fetch data using Resty
	resp, err := c.http.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status())
	}

	// Parse XML
	var feed model.RSSFeed
	if err := xml.Unmarshal(resp.Body(), &feed); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &feed, nil
}

// FetchAlertDetail fetches the detailed CAP XML from the given URL
func (c *Client) FetchAlertDetail(url string) (*model.CAPAlert, error) {
	resp, err := c.http.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alert detail: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status())
	}

	var alert model.CAPAlert
	if err := xml.Unmarshal(resp.Body(), &alert); err != nil {
		return nil, fmt.Errorf("failed to parse alert XML: %w", err)
	}

	return &alert, nil
}

// DownloadImage fetches the image bytes from a URL
func (c *Client) DownloadImage(url string) ([]byte, error) {
	resp, err := c.http.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("HTTP error image download: %s", resp.Status())
	}
	return resp.Body(), nil
}
