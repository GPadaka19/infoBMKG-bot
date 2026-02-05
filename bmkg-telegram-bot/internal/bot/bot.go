package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// Print bot info
	fmt.Printf("Authorized on account %s\n", api.Self.UserName)

	return &Bot{api: api}, nil
}

// SendMessage sends a text message to a specific chat ID
func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML" // Enable HTML parsing for bold/links
	msg.DisableWebPagePreview = false

	_, err := b.api.Send(msg)
	return err
}

// GetAPI returns the raw API instance if needed
func (b *Bot) GetAPI() *tgbotapi.BotAPI {
	return b.api
}

// SendPhoto sends a photo with a caption using URL (Not Recommended for stability)
func (b *Bot) SendPhoto(chatID int64, photoURL string, caption string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(photoURL))
	photo.Caption = caption
	photo.ParseMode = "HTML"

	_, err := b.api.Send(photo)
	return err
}

// SendPhotoBytes uploads a photo from memory bytes (Recommended)
func (b *Bot) SendPhotoBytes(chatID int64, photoData []byte, caption string) error {
	fileBytes := tgbotapi.FileBytes{
		Name:  "image.jpg",
		Bytes: photoData,
	}
	photo := tgbotapi.NewPhoto(chatID, fileBytes)
	photo.Caption = caption
	photo.ParseMode = "HTML"

	_, err := b.api.Send(photo)
	return err
}
