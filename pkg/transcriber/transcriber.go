package transcriber

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rashad-j/ai-subtitle/pkg/configs"
	logger "github.com/rs/zerolog/log"
	transloadit "github.com/transloadit/go-sdk"
)

var ErrAssembly error = errors.New("failed to complete the assembly")

func Subtitle(cfg configs.Config, file, provider string) (string, error) {
	options := transloadit.DefaultConfig
	options.AuthKey = cfg.TransloaditAuthKey
	options.AuthSecret = cfg.TransloaditAuthSecret

	client := transloadit.NewClient(options)

	// Initialize a new assembly
	assembly := transloadit.NewAssembly()

	if err := assembly.AddFile("audio", file); err != nil {
		return "", err
	}

	assembly.AddStep("transcribed", map[string]interface{}{
		"use":      ":original",
		"robot":    "/speech/transcribe",
		"provider": provider,
		"format":   "srt",
		"result":   true,
	})

	info, err := client.StartAssembly(context.Background(), assembly)
	if err != nil {
		return "", err
	}

	info, err = client.WaitForAssembly(context.Background(), info)
	if err != nil {
		return "", err
	}

	if info.Ok != "ASSEMBLY_COMPLETED" {
		logger.Err(errors.New(info.Error)).Str("transloadit", info.Message).Msg("failed to upload")
		return "", ErrAssembly
	}

	result, ok := info.Results["transcribed"]
	if !ok || len(result) == 0 || result[0].SSLURL == "" {
		return "", errors.New("failed to get transcription URL")
	}

	url := result[0].SSLURL
	return url, nil
}

func Download(filepath, name, url string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
		return err
	}

	// Create the file with the specified name
	file := filepath + "/" + name
	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Make the HTTP request
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check for non-200 status codes
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	// Copy the response body to the file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return err
	}

	// logger.Info().Str("path", filepath).Str("name", name).Msg("file downloaded")

	return nil
}
