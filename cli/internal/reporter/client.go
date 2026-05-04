package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ashczar77/mockingjay/internal/test"
)

// Reporter is the interface for sending test results to a backend
type Reporter interface {
	Report(result test.Result) error
}

// Transcription holds transcription data to report
type Transcription struct {
	CallSID    string
	AudioPath  string
	Text       string
	Confidence float64
	Duration   float64
}

// RawResult holds a raw pass/fail result to report
type RawResult struct {
	Scenario  string
	Passed    bool
	LatencyMs int64
	Error     string
}

// Client sends test results to the backend API
type Client struct {
	httpClient *http.Client
	apiURL     string
}

// NewClient creates a new reporter client
func NewClient(apiURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiURL: apiURL,
	}
}

// Report sends a test result to the backend
func (c *Client) Report(result test.Result) error {
	payload := map[string]interface{}{
		"scenario":   result.Scenario,
		"passed":     result.Passed,
		"latency_ms": result.Metrics.Latency.Milliseconds(),
		"error":      result.Error,
	}
	return c.post("/api/results", payload)
}

// ReportRaw sends a raw pass/fail result to the backend
func (c *Client) ReportRaw(r RawResult) error {
	return c.post("/api/results", map[string]interface{}{
		"scenario":   r.Scenario,
		"passed":     r.Passed,
		"latency_ms": r.LatencyMs,
		"error":      r.Error,
	})
}

// ReportTranscription sends a transcription to the backend
func (c *Client) ReportTranscription(t Transcription) error {
	return c.post("/api/transcriptions", map[string]interface{}{
		"call_sid":         t.CallSID,
		"audio_path":       t.AudioPath,
		"text":             t.Text,
		"confidence":       t.Confidence,
		"duration_seconds": t.Duration,
	})
}

func (c *Client) post(path string, payload map[string]interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.apiURL+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send to backend: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned status %d", resp.StatusCode)
	}
	return nil
}
