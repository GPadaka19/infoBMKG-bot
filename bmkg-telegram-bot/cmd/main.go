package main

import (
	"bmkg-bot/internal/bmkg"
	"bmkg-bot/internal/bot"
	"bmkg-bot/internal/model"
	"bmkg-bot/internal/storage"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Configuration variables
var (
	// Get from Environment Variables for security
	BotToken     string
	TargetChatID int64
	PollInterval = 3 * time.Minute
	FilteredProv []string // List of provinces to filter
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	BotToken = os.Getenv("TELEGRAM_TOKEN")
	chatIDStr := os.Getenv("TARGET_CHAT_ID")
	if chatIDStr != "" {
		TargetChatID, _ = strconv.ParseInt(chatIDStr, 10, 64)
	}

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

	log.Println("Starting BMKG Weather Alert Bot...")

	// 1. Initialize Storage (History)
	history := storage.NewHistory()
	log.Printf("Loaded history with %d items", len(history.Items))

	// 2. Initialize BMKG Client
	bmkgClient := bmkg.NewClient()

	// 3. Initialize Telegram Bot
	var telegramBot *bot.Bot
	var err error

	if BotToken != "" {
		telegramBot, err = bot.NewBot(BotToken)
		if err != nil {
			log.Fatalf("Failed to init Telegram Bot: %v", err)
		}
	} else {
		log.Println("WARNING: TELEGRAM_TOKEN is not set. Creating bot in DRY RUN mode (Console output only).")
	}

	// 4. Start Polling Loop
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	// Run immediately once at startup
	processUpdate(bmkgClient, history, telegramBot)

	for range ticker.C {
		processUpdate(bmkgClient, history, telegramBot)
	}
}

func processUpdate(client *bmkg.Client, history *storage.History, bot *bot.Bot) {
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
					// log.Printf("Matched filter: %s needed %s", item.Title, p)
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
			// Pass client to download image
			sendNotification(bot, client, item, photoURL)

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

func sendNotification(b *bot.Bot, client *bmkg.Client, item model.Item, photoURL string) {
	// Construct message
	body := strings.TrimSpace(item.Description)

	message := fmt.Sprintf(
		"⚠️ <b>PERINGATAN DINI CUACA</b> ⚠️\n\n"+
			"<b>%s</b>\n\n"+
			"%s\n\n"+
			"🕒 <i>%s</i>\n"+
			"🔗 <a href='%s'>Selengkapnya</a>\n\n"+
			"📡 <i>Sumber Data: BMKG (Badan Meteorologi, Klimatologi, dan Geofisika)</i>",
		item.Title,
		body,
		item.PubDate,
		item.Link,
	)

	// Truncate logic for Caption (Max 1024 chars in Telegram)
	caption := message
	if len(caption) > 1000 {
		// Cut body to fit
		allowedBody := 1000 - (len(message) - len(body)) - 50 // buffer
		if allowedBody > 0 && len(body) > allowedBody {
			truncatedBody := body[:allowedBody] + "..."
			caption = fmt.Sprintf(
				"⚠️ <b>PERINGATAN DINI CUACA</b> ⚠️\n\n"+
					"<b>%s</b>\n\n"+
					"%s\n\n"+
					"🕒 <i>%s</i>\n"+
					"🔗 <a href='%s'>Selengkapnya</a>\n\n"+
					"📡 <i>Sumber Data: BMKG (Badan Meteorologi, Klimatologi, dan Geofisika)</i>",
				item.Title,
				truncatedBody,
				item.PubDate,
				item.Link,
			)
		}
	}

	// Dry Run Check
	if b == nil || TargetChatID == 0 {
		fmt.Println("---------------------------------------------------")
		fmt.Println(" [DRY RUN] TELEGRAM MESSAGE PREVIEW:")
		if photoURL != "" {
			fmt.Printf(" [PHOTO] %s\n", photoURL)
		}
		// Print full message for debug
		fmt.Println(message)
		fmt.Println("---------------------------------------------------")
		return
	}

	// Logic: Try SendPhotoBytes first
	var err error
	sent := false

	if photoURL != "" {
		// Download image first
		imgBytes, errDown := client.DownloadImage(photoURL)
		if errDown == nil {
			err = b.SendPhotoBytes(TargetChatID, imgBytes, caption)
			if err == nil {
				sent = true
			} else {
				log.Printf("Failed to upload photo: %v. Falling back to text.", err)
			}
		} else {
			log.Printf("Failed to download photo content: %v. Falling back to text.", errDown)
		}
	}

	// Fallback to text if no photo or photo failed
	if !sent {
		err = b.SendMessage(TargetChatID, message)
		if err != nil {
			log.Printf("Failed to send telegram message: %v", err)
		}
	}
}
