package audio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Transcriber handles audio transcription via ASR APIs
type Transcriber interface {
	Transcribe(audioPath string) (*Transcript, error)
}

// Transcript holds the transcription result
type Transcript struct {
	Text       string
	Confidence float64
	Duration   float64 // seconds
	Words      []Word
}

// Word represents a transcribed word with timing
type Word struct {
	Text       string
	Start      float64
	End        float64
	Confidence float64
}

// DeepgramTranscriber uses Deepgram API for transcription
type DeepgramTranscriber struct {
	apiKey string
	client *http.Client
}

// NewDeepgramTranscriber creates a Deepgram transcriber
func NewDeepgramTranscriber(apiKey string) *DeepgramTranscriber {
	return &DeepgramTranscriber{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Transcribe sends audio to Deepgram and returns transcript
func (d *DeepgramTranscriber) Transcribe(audioPath string) (*Transcript, error) {
	data, err := os.ReadFile(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.deepgram.com/v1/listen?model=general&punctuate=true", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+d.apiKey)
	req.Header.Set("Content-Type", "audio/wav")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deepgram API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deepgram API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Results struct {
			Channels []struct {
				Alternatives []struct {
					Transcript string  `json:"transcript"`
					Confidence float64 `json:"confidence"`
					Words      []struct {
						Word       string  `json:"word"`
						Start      float64 `json:"start"`
						End        float64 `json:"end"`
						Confidence float64 `json:"confidence"`
					} `json:"words"`
				} `json:"alternatives"`
			} `json:"channels"`
		} `json:"results"`
		Metadata struct {
			Duration float64 `json:"duration"`
		} `json:"metadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse deepgram response: %w", err)
	}

	if len(result.Results.Channels) == 0 || len(result.Results.Channels[0].Alternatives) == 0 {
		return nil, fmt.Errorf("no transcription results")
	}

	alt := result.Results.Channels[0].Alternatives[0]
	transcript := &Transcript{
		Text:       alt.Transcript,
		Confidence: alt.Confidence,
		Duration:   result.Metadata.Duration,
		Words:      make([]Word, len(alt.Words)),
	}

	for i, w := range alt.Words {
		transcript.Words[i] = Word{
			Text:       w.Word,
			Start:      w.Start,
			End:        w.End,
			Confidence: w.Confidence,
		}
	}

	return transcript, nil
}

// SaveAudio downloads audio from URL and saves to local file
func SaveAudio(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save audio: %w", err)
	}

	return nil
}
