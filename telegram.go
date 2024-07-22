package main

import (
	"context"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const stickerFileID = "CAACAgIAAxkBAAEMgeJmmAk7yiqL45fsxmRbOzWXaHbytwAC_TYAAqMPoEvi4Nvo6UXPwTUE"

func handleUpdates(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, client *speech.Client) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Handle voice messages
		if update.Message.Voice != nil {
			log.Printf("Received voice message from user %s", update.Message.From.UserName)

			voiceFileURL, err := bot.GetFileDirectURL(update.Message.Voice.FileID)
			if err != nil {
				log.Printf("Failed to get voice file URL: %v", err)
				continue
			}

			voiceFilePath, err := downloadFile(voiceFileURL, "voice.oga")
			if err != nil {
				log.Printf("Failed to download voice file: %v", err)
				continue
			}

			response, err := transcribeVoice(ctx, client, voiceFilePath)
			if err != nil {
				log.Printf("Error transcribing voice: %v", err)
				continue
			}

			sendResponse(bot, update.Message.Chat.ID, response)
		}

		// Handle video messages (video circles)
		if update.Message.VideoNote != nil {
			log.Printf("Received video message from user %s", update.Message.From.UserName)

			videoFileURL, err := bot.GetFileDirectURL(update.Message.VideoNote.FileID)
			if err != nil {
				log.Printf("Failed to get video file URL: %v", err)
				continue
			}

			videoFilePath, err := downloadFile(videoFileURL, "video.mp4")
			if err != nil {
				log.Printf("Failed to download video file: %v", err)
				continue
			}

			audioFilePath := "voices/audio.oga"
			err = extractAudio(videoFilePath, audioFilePath)
			if err != nil {
				log.Printf("Failed to extract audio from video file: %v", err)
				continue
			}

			response, err := transcribeVoice(ctx, client, audioFilePath)
			if err != nil {
				log.Printf("Error transcribing audio: %v", err)
				continue
			}

			sendResponse(bot, update.Message.Chat.ID, response)
		}
	}
}

func sendResponse(bot *tgbotapi.BotAPI, chatID int64, response string) {
	sticker := tgbotapi.NewStickerShare(chatID, stickerFileID)
	_, err := bot.Send(sticker)
	if err != nil {
		log.Printf("Failed to send sticker: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, response)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
