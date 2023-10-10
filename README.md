# AI Transcription CLI

This CLI tool uses [transloadit](https://transloadit.com). Transloadit is an online service that can easily modify files using api calls. It enables easy access to AI from AWS or GCP.

This CLI tool is a command-line tool that allows you to generate subtitles from video or audio files using AI services from AWS or GCP. This tool provides several commands to download files, transcribe audio, and transcribe audio while downloading the transcription.

## Installation

To use the AI Subtitle CLI, you'll need to build it from source. Make sure you have Go installed on your system.

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/ai-subtitle-cli.git
   cd ai-subtitle-cli
   ```
2. Build the CLI:

    ```bash
    go build -o ai-subtitle
    ```

3. Run:
    ```bash
    ./ai-subtitle --help
    ```

## Usage
This CLI tool provides the following commands:

### Transcribe
Transcribe an audio file using an AI service (AWS or GCP). This command will return a URL to download the SRT file. Use the `Download` command to download from a URL.

```bash
./ai-subtitle transcribe -f /path/to/audio/file.wav -s aws
```
```
-f,--file: Audio file to transcribe
-s,--service: AI provider to use was or gcp, default is aws
```

### Download
Download a transcription from a URL
```bash
./ai-subtitle download -p /path/to/destination -n filename.ext -u https://example.com/file.ext
```

### Transdownload
Transcribe an audio file and download the transcription.
```bash
./ai-subtitle transdownload -f /path/to/audio/file.wav -n transcription.srt -p /path/to/destination -s gcp
```
```
-f, --file: Audio file to transcribe.
-n, --name: Name of the downloaded transcription file.
-p, --path: Destination directory for the downloaded transcription.
-s, --service: AI service to use (aws or gcp), default is AWS.
```

## Configuration
The CLI uses environment variables for configuration. You can set these variables directly or create a `.env` file in the project directory, or simple in your environment variables.

Ensure that you have created an account at [transloadit](https://transloadit.com). Then you you can create an auth key and an auth secret:

```bash
TRANSLOADIT_AUTH_KEY="your-key"
TRANSLOADIT_AUTH_SECRET="your-secret"
```

## Contribution
If you'd like to contribute to the AI Subtitle CLI, please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Create a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
