package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rashad-j/ai-subtitle/pkg/configs"
	"github.com/rashad-j/ai-subtitle/pkg/transcriber"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "transcriber",
		Usage: "A CLI application for generating subtitles from video or audio files using AI from AWS or GCP.",
		Commands: []*cli.Command{
			downloadCommand(),
			transcribeCommand(),
			transdownloadCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func downloadCommand() *cli.Command {
	return &cli.Command{
		Name:    "download",
		Aliases: []string{"d"},
		Usage:   "Download a file from a URL.",
		Action: func(c *cli.Context) error {
			path := c.String("path")
			name := c.String("name")
			url := c.String("url")

			// Start the loading spinner
			s := startSpinner("Downloading...")
			defer s.Stop()

			// Ensure the filename has the .srt extension
			name = ensureSRTExtension(name)

			if err := transcriber.Download(path, name, url); err != nil {
				fmt.Println("Error:", err)
				return err
			}

			fmt.Println("Downloaded successfully.")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Destination directory for downloaded file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the downloaded file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "URL to download the file from",
				Required: true,
			},
		},
	}
}

func transcribeCommand() *cli.Command {
	return &cli.Command{
		Name:    "transcribe",
		Aliases: []string{"t"},
		Usage:   "Transcribe an audio file.",
		Action: func(c *cli.Context) error {
			file := c.String("file")
			aiService := c.String("service")

			cfg, err := configs.LoadConfigs()
			if err != nil {
				fmt.Println("Error loading config:", err)
				return err
			}

			// Start the loading spinner
			s := startSpinner("Transcribing...")
			defer s.Stop()

			transcriptionURL, err := transcriber.Subtitle(cfg, file, aiService)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}

			fmt.Println("Transcription URL:", transcriptionURL)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Audio file to transcribe",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "AI service to use (aws or gcp), default aws",
				Value:   "aws", // Default value is AWS
			},
		},
	}
}

func transdownloadCommand() *cli.Command {
	return &cli.Command{
		Name:    "transdownload",
		Aliases: []string{"td"},
		Usage:   "Transcribe an audio file and download the transcription.",
		Action: func(c *cli.Context) error {
			file := c.String("file")
			transcriptionPath := c.String("path")
			transcriptionName := c.String("name")
			aiService := c.String("service")

			cfg, err := configs.LoadConfigs()
			if err != nil {
				fmt.Println("Error loading config:", err)
				return err
			}

			// Start the loading spinner
			s := startSpinner("Transcribing & Downloading...")
			defer s.Stop()

			transcriptionURL, err := transcriber.Subtitle(cfg, file, aiService)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}

			// Ensure the filename has the .srt extension
			transcriptionName = ensureSRTExtension(transcriptionName)

			// Download the transcription
			if err := transcriber.Download(transcriptionPath, transcriptionName, transcriptionURL); err != nil {
				fmt.Println("Error downloading transcription:", err)
				return err
			}

			fmt.Println("Transcription downloaded successfully.")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Audio file to transcribe",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the downloaded file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Destination directory for downloaded transcription",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "AI service to use (aws or gcp), default aws",
				Value:   "aws", // Default value is AWS
			},
		},
	}
}

func startSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = suffix
	s.Start()
	return s
}

func ensureSRTExtension(name string) string {
	// Check if the filename has the .srt extension
	if !strings.HasSuffix(name, ".srt") {
		// Remove any existing extension and add .srt
		name = strings.TrimSuffix(name, filepath.Ext(name)) + ".srt"
	}
	return name
}
