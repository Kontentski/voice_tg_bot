package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func extractAudio(videoFilePath, audioFilePath string) error {
	if err := os.MkdirAll("voices", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	err := ffmpeg.Input(videoFilePath).
		Output(audioFilePath, ffmpeg.KwArgs{"vn": "", "acodec": "libopus"}).
		OverWriteOutput().
		Run()

	if err != nil {
		return fmt.Errorf("ffmpeg command failed: %w", err)
	}

	log.Printf("Successfully extracted audio to: %s", audioFilePath)
	return nil
}

func downloadFile(url, fileName string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	voiceDir := "./voices"
	if err := os.MkdirAll(voiceDir, os.ModePerm); err != nil {
		return "", err
	}

	filePath := filepath.Join(voiceDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
