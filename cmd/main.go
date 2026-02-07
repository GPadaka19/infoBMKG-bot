package main

import (
	"bmkg-bot/internal/bmkg"
	"bmkg-bot/internal/model"
	"bmkg-bot/internal/storage"
	"bmkg-bot/internal/whatsapp"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Configuration variables
var (
	// WhatsApp Configuration
	WaGatewayURL  string
	WaTargetPhone string
	WaUser        string
	WaPassword    string

	PollInterval = 3 * time.Minute
	FilteredProv []string // List of provinces to filter
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// WhatsApp Config
	WaGatewayURL = os.Getenv("WA_GATEWAY_URL")
	WaTargetPhone = os.Getenv("WA_TARGET_PHONE")
	WaUser = os.Getenv("WA_USER")
	WaPassword = os.Getenv("WA_PASSWORD")

	// Parse Province Filter
	provStr := os.Getenv("FILTER_PROVINCE")
	if provStr != "" {
		rawProvs := strings.Split(provStr, ",")
		for _, p := range rawProvs {
			p = strings.TrimSpace(p)
			if p != "" {
				FilteredProv = append(FilteredProv, p)
			}
		}
		log.Printf("Active Filter: %v", FilteredProv)
	}

	log.Println("Starting BMKG Weather Alert Bot (WhatsApp Edition)...")

	// 1. Initialize Storage (History)
	history := storage.NewHistory()
	log.Printf("Loaded history with %d items", len(history.Items))

	// 2. Initialize BMKG Client
	bmkgClient := bmkg.NewClient()

	// 3. Initialize WhatsApp Client
	var waClient *whatsapp.Client

	if WaGatewayURL != "" && WaTargetPhone != "" {
		waClient = whatsapp.NewClient(WaGatewayURL, WaTargetPhone, WaUser, WaPassword)

		// Check connectivity
		if err := waClient.CheckConnectivity(); err != nil {
			log.Printf("⚠️ [WARN] WhatsApp Gateway connectivity issue: %v", err)
		} else {
			log.Println("✅ WhatsApp Gateway connected")
		}
	} else {
		log.Println("WARNING: WA_GATEWAY_URL or WA_TARGET_PHONE is not set. Creating bot in DRY RUN mode (Console output only).")
	}

	// 4. Start Polling Loop
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	// Run immediately once at startup
	processUpdate(bmkgClient, history, waClient)

	for range ticker.C {
		processUpdate(bmkgClient, history, waClient)
	}
}

func processUpdate(client *bmkg.Client, history *storage.History, waClient *whatsapp.Client) {
	log.Println("Checking for updates...")

	// Fetch Indonesian feed
	feed, err := client.FetchRSS("id")
	if err != nil {
		log.Printf("Error fetching RSS: %v", err)
		return
	}

	newCount := 0

	for _, item := range feed.Channel.Items {
		guid := item.GUID.Value
		if guid == "" {
			guid = item.Link
		}

		// SKIP if filtered
		if len(FilteredProv) > 0 {
			matched := false
			for _, p := range FilteredProv {
				// Check Title ONLY (usually "... di <Provinsi>") to avoid false positives in description (e.g. "bali" in "kembali")
				if strings.Contains(strings.ToLower(item.Title), strings.ToLower(p)) {
					matched = true
					break
				}
			}
			if !matched {
				continue // Skip this item
			}
		}

		if history.IsNew(guid) {
			newCount++
			log.Printf("New Alert Found: %s", item.Title)

			// Fetch Detail to get Infographic
			var photoURL string
			if item.Link != "" {
				detail, err := client.FetchAlertDetail(item.Link)
				if err == nil {
					photoURL = detail.Info.Web
				} else {
					log.Printf("Failed to fetch detail for %s: %v", item.Title, err)
				}
			}

			// 1. Send Notification
			sendNotification(waClient, client, item, photoURL)

			// 2. Add to History (Mark as seen)
			if err := history.Add(guid); err != nil {
				log.Printf("Error processing history: %v", err)
			}
		}
	}

	if newCount == 0 {
		log.Println("No new alerts.")
	} else {
		log.Printf("Processed %d new alerts.", newCount)
	}
}

func sendNotification(waClient *whatsapp.Client, client *bmkg.Client, item model.Item, photoURL string) {
	// Construct message - WhatsApp format (plain text with emoji, no HTML)
	body := strings.TrimSpace(item.Description)

	message := fmt.Sprintf(
		"⚠️ *PERINGATAN DINI CUACA* ⚠️\n\n"+
			"*%s*\n\n"+
			"%s\n\n"+
			"🕒 _%s_\n"+
			"🔗 %s\n\n"+
			"📡 _Sumber Data: BMKG (Badan Meteorologi, Klimatologi, dan Geofisika)_",
		item.Title,
		body,
		item.PubDate,
		item.Link,
	)

	// Caption for photo (shorter version if needed)
	caption := message
	if len(caption) > 1000 {
		// Cut body to fit
		allowedBody := 1000 - (len(message) - len(body)) - 50 // buffer
		if allowedBody > 0 && len(body) > allowedBody {
			truncatedBody := body[:allowedBody] + "..."
			caption = fmt.Sprintf(
				"⚠️ *PERINGATAN DINI CUACA* ⚠️\n\n"+
					"*%s*\n\n"+
					"%s\n\n"+
					"🕒 _%s_\n"+
					"🔗 %s\n\n"+
					"📡 _Sumber Data: BMKG (Badan Meteorologi, Klimatologi, dan Geofisika)_",
				item.Title,
				truncatedBody,
				item.PubDate,
				item.Link,
			)
		}
	}

	// Dry Run Check
	if waClient == nil {
		fmt.Println("---------------------------------------------------")
		fmt.Println(" [DRY RUN] WHATSAPP MESSAGE PREVIEW:")
		if photoURL != "" {
			fmt.Printf(" [PHOTO] %s\n", photoURL)
		}
		// Print full message for debug
		fmt.Println(message)
		fmt.Println("---------------------------------------------------")
		return
	}

	// Logic: Try SendPhoto first if we have an image
	var err error
	sent := false

	if photoURL != "" {
		// Download image first
		imgBytes, errDown := client.DownloadImage(photoURL)
		if errDown == nil {
			err = waClient.SendPhoto(imgBytes, caption)
			if err == nil {
				sent = true
			} else {
				log.Printf("Failed to send photo via WhatsApp: %v. Falling back to text.", err)
			}
		} else {
			log.Printf("Failed to download photo content: %v. Falling back to text.", errDown)
		}
	}

	// Fallback to text if no photo or photo failed
	if !sent {
		err = waClient.SendMessage(message)
		if err != nil {
			log.Printf("Failed to send WhatsApp message: %v", err)
		}
	}
}
