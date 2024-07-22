package main

import (
	"context"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Set up Google Cloud
	ctx := context.Background()
	client, err := setupGoogleCloudClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Initialize Telegram
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	// Set up webhook configuration
	webhookURL := os.Getenv("WEBHOOK_URL")
	webhookURL = webhookURL + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}
	log.Printf("Webhook server started. Listening for Telegram updates at %s", webhookURL)

	// Process incoming updates
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":8080", nil)

	handleUpdates(ctx, bot, updates, client)
}
