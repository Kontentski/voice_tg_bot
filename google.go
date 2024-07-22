package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
)

func setupGoogleCloudClient(ctx context.Context) (*speech.Client, error) {
	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return client, nil
}

func transcribeVoice(ctx context.Context, client *speech.Client, voiceFilePath string) (string, error) {
	audioData, err := os.ReadFile(voiceFilePath)
	if err != nil {
		return "", err
	}

	if len(audioData) == 0 {
		return "", fmt.Errorf("audio data is empty")
	}

	// Log the file size and type for debugging
	log.Printf("Voice file size: %d bytes", len(audioData))
	log.Printf("Voice file path: %s", voiceFilePath)

	// Call Google Cloud Speech-to-Text API
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_OGG_OPUS,
			SampleRateHertz:            48000,
			LanguageCode:               "uk-UA",
			AlternativeLanguageCodes:   []string{"en-US"},
			EnableAutomaticPunctuation: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: audioData,
			},
		},
	})
	if err != nil {
		return "", err
	}

	// Log the full API response for debugging
	log.Printf("API response: %v", resp)

	if len(resp.Results) == 0 {
		return "", fmt.Errorf("no transcription results")
	}

	var sb strings.Builder
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			sb.WriteString(alt.Transcript)
		}
	}

	transcription := sb.String()
	log.Printf("Transcription: %s", transcription)

	if transcription == "" {
		return "", fmt.Errorf("transcription is empty")
	}

	return transcription, nil
}
