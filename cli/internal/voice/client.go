package voice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with voice AI agents
type Client struct {
	httpClient *http.Client
	endpoint   string
}

// NewClient creates a new voice AI client
func NewClient(endpoint string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint: endpoint,
	}
}

// CallRequest represents a request to the voice AI
type CallRequest struct {
	Text string `json:"text"`
}

// CallResponse represents the voice AI response
type CallResponse struct {
	Text     string `json:"text"`
	Intent   string `json:"intent,omitempty"`
	Latency  int64  `json:"latency_ms,omitempty"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// Call sends text to the voice AI and gets a response
func (c *Client) Call(text string) (*CallResponse, time.Duration, error) {
	start := time.Now()

	reqBody := CallRequest{Text: text}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call endpoint: %w", err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, latency, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &CallResponse{
			Success: false,
			Error:   fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
		}, latency, nil
	}

	var callResp CallResponse
	if err := json.Unmarshal(body, &callResp); err != nil {
		return nil, latency, fmt.Errorf("failed to parse response: %w", err)
	}

	callResp.Success = true
	return &callResp, latency, nil
}
