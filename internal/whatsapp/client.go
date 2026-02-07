package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Client represents a WhatsApp gateway client
type Client struct {
	gatewayURL  string
	targetPhone string
	user        string
	password    string
	httpClient  *http.Client
}

// TextPayload for sending text messages
type TextPayload struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// ImagePayload for sending images with caption via URL
type ImagePayload struct {
	Phone    string `json:"phone"`
	Caption  string `json:"caption"`
	ImageURL string `json:"image_url"`
	Compress bool   `json:"compress"`
}

// NewClient creates a new WhatsApp client
func NewClient(gatewayURL, targetPhone, user, password string) *Client {
	// Sanitize phone number
	targetPhone = formatPhoneNumber(targetPhone)
	log.Printf("📱 WhatsApp Target: %s", targetPhone)

	return &Client{
		gatewayURL:  gatewayURL,
		targetPhone: targetPhone,
		user:        user,
		password:    password,
		httpClient:  &http.Client{Timeout: 20 * time.Second},
	}
}

// formatPhoneNumber formats phone/group ID for WhatsApp gateway
// Formats supported:
// - Group: 120363xxx@g.us (left as-is)
// - Individual: 6281xxx@s.whatsapp.net or just 6281xxx
func formatPhoneNumber(phone string) string {
	// If already contains @g.us (group), return as-is
	if strings.Contains(phone, "@g.us") {
		return phone
	}

	// If already contains @s.whatsapp.net (individual), return as-is
	if strings.Contains(phone, "@s.whatsapp.net") {
		return phone
	}

	// Sanitize phone number for individual
	reg, _ := regexp.Compile("[^0-9]+")
	clean := reg.ReplaceAllString(phone, "")
	if strings.HasPrefix(clean, "08") {
		return "62" + clean[1:]
	}
	if strings.HasPrefix(clean, "8") {
		return "62" + clean
	}
	return clean
}

// CheckConnectivity tests connection to the WhatsApp gateway
func (c *Client) CheckConnectivity() error {
	u, _ := url.Parse(c.gatewayURL)
	parts := strings.Split(u.Path, "/")
	basePath := "/"
	if len(parts) > 1 {
		basePath = "/" + parts[1]
	}

	testURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, basePath)
	req, _ := http.NewRequest("GET", testURL, nil)
	if c.user != "" {
		req.SetBasicAuth(c.user, c.password)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("⛔ Basic Auth Error: Invalid credentials")
	}
	return nil
}

// SendMessage sends a text message with auto-healing
func (c *Client) SendMessage(text string) error {
	err := c.sendText(text)
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "AUTHENTICATION_ERROR") || strings.Contains(errMsg, "please reconnect") {
		log.Println("🚑 [HEALING] Reconnecting...")
		c.triggerReconnect()
		time.Sleep(5 * time.Second)
		return c.sendText(text)
	}
	return err
}

// SendPhoto sends an image with caption via URL (auto-healing enabled)
func (c *Client) SendPhoto(imageURL string, caption string) error {
	err := c.sendImage(imageURL, caption)
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "AUTHENTICATION_ERROR") || strings.Contains(errMsg, "please reconnect") {
		log.Println("🚑 [HEALING] Reconnecting...")
		c.triggerReconnect()
		time.Sleep(5 * time.Second)
		return c.sendImage(imageURL, caption)
	}
	return err
}

// sendText sends a plain text message
func (c *Client) sendText(text string) error {
	payload := TextPayload{Phone: c.targetPhone, Message: text}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", c.gatewayURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		req.SetBasicAuth(c.user, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API %d: %s", resp.StatusCode, string(body))
	}
	log.Printf("✅ [SENT] Text message sent successfully")
	return nil
}

// sendImage sends an image via URL with caption
func (c *Client) sendImage(imageURL string, caption string) error {
	// Build image endpoint URL (replace /send/message with /send/image)
	endpointURL := strings.Replace(c.gatewayURL, "/send/message", "/send/image", 1)

	payload := ImagePayload{
		Phone:    c.targetPhone,
		Caption:  caption,
		ImageURL: imageURL,
		Compress: true,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", endpointURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		req.SetBasicAuth(c.user, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API %d: %s", resp.StatusCode, string(body))
	}
	log.Printf("✅ [SENT] Image sent successfully")
	return nil
}

// triggerReconnect attempts to reconnect the WhatsApp session
func (c *Client) triggerReconnect() {
	u, _ := url.Parse(c.gatewayURL)
	u.Path = strings.Replace(u.Path, "/send/message", "/app/reconnect", 1)
	req, _ := http.NewRequest("GET", u.String(), nil)
	if c.user != "" {
		req.SetBasicAuth(c.user, c.password)
	}
	(&http.Client{Timeout: 10 * time.Second}).Do(req)
}

// GetTargetPhone returns the formatted target phone number
func (c *Client) GetTargetPhone() string {
	return c.targetPhone
}
