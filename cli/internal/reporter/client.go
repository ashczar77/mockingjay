package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ashczar77/mockingjay/internal/test"
)

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

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	resp, err := c.httpClient.Post(c.apiURL+"/api/results", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned status %d", resp.StatusCode)
	}

	return nil
}
