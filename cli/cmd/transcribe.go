package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashczar77/mockingjay/internal/audio"
	"github.com/spf13/cobra"
)

var (
	deepgramKey  string
	audioURL     string
	audioFile    string
	outputDir    string
)

var transcribeCmd = &cobra.Command{
	Use:   "transcribe",
	Short: "Transcribe audio from a file or URL",
	Long:  `Transcribe a call recording using Deepgram ASR.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runTranscribe(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	transcribeCmd.Flags().StringVar(&deepgramKey, "deepgram-key", os.Getenv("DEEPGRAM_API_KEY"), "Deepgram API key")
	transcribeCmd.Flags().StringVar(&audioURL, "url", "", "URL of audio file to transcribe")
	transcribeCmd.Flags().StringVar(&audioFile, "file", "", "Local audio file to transcribe")
	transcribeCmd.Flags().StringVar(&outputDir, "output-dir", ".", "Directory to save downloaded audio")
}

func runTranscribe() error {
	if deepgramKey == "" {
		return fmt.Errorf("Deepgram API key required: --deepgram-key or DEEPGRAM_API_KEY env var")
	}
	if audioURL == "" && audioFile == "" {
		return fmt.Errorf("provide --url or --file")
	}

	fmt.Println("🐦 MockingJay - Transcribe")
	fmt.Println()

	localPath := audioFile

	// Download if URL provided
	if audioURL != "" {
		localPath = filepath.Join(outputDir, "recording.wav")
		fmt.Printf("⬇️  Downloading audio from %s...\n", audioURL)
		if err := audio.SaveAudio(audioURL, localPath); err != nil {
			return err
		}
		fmt.Printf("   Saved to %s\n\n", localPath)
	}

	fmt.Printf("🎙️  Transcribing %s...\n\n", localPath)

	t := audio.NewDeepgramTranscriber(deepgramKey)
	transcript, err := t.Transcribe(localPath)
	if err != nil {
		return err
	}

	fmt.Printf("📝 Transcript:\n  %s\n\n", transcript.Text)
	fmt.Printf("📊 Stats:\n")
	fmt.Printf("  Confidence: %.1f%%\n", transcript.Confidence*100)
	fmt.Printf("  Duration:   %.1fs\n", transcript.Duration)
	fmt.Printf("  Words:      %d\n", len(transcript.Words))

	return nil
}
