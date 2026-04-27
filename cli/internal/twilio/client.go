package twilio

import (
	"fmt"
	"time"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/twilio/twilio-go"
)

// Config holds Twilio credentials and call settings
type Config struct {
	AccountSID string
	AuthToken  string
	FromNumber string // Twilio phone number
	ToNumber   string // target phone number to call
	WebhookURL string // TwiML webhook URL for call instructions
}

// CallResult holds the outcome of a Twilio call
type CallResult struct {
	CallSID    string
	Status     string
	Duration   time.Duration
	RecordURL  string
	Error      string
}

// Client wraps the Twilio API for voice call testing
type Client struct {
	client *twilio.RestClient
	cfg    Config
}

// NewClient creates a new Twilio client
func NewClient(cfg Config) *Client {
	c := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})
	return &Client{client: c, cfg: cfg}
}

// MakeCall initiates an outbound call and waits for completion
func (c *Client) MakeCall(record bool) (*CallResult, error) {
	params := &twilioApi.CreateCallParams{}
	params.SetTo(c.cfg.ToNumber)
	params.SetFrom(c.cfg.FromNumber)
	params.SetUrl(c.cfg.WebhookURL)
	if record {
		params.SetRecord(true)
	}

	resp, err := c.client.Api.CreateCall(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create call: %w", err)
	}

	result := &CallResult{CallSID: *resp.Sid, Status: *resp.Status}

	// Poll until call completes (max 5 minutes)
	result, err = c.waitForCompletion(result.CallSID)
	if err != nil {
		return result, err
	}

	return result, nil
}

// waitForCompletion polls the call status until it reaches a terminal state
func (c *Client) waitForCompletion(callSID string) (*CallResult, error) {
	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		time.Sleep(3 * time.Second)

		resp, err := c.client.Api.FetchCall(callSID, &twilioApi.FetchCallParams{})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch call status: %w", err)
		}

		status := ""
		if resp.Status != nil {
			status = *resp.Status
		}

		switch status {
		case "completed", "failed", "busy", "no-answer", "canceled":
			result := &CallResult{CallSID: callSID, Status: status}
			if resp.Duration != nil {
				secs, _ := time.ParseDuration(*resp.Duration + "s")
				result.Duration = secs
			}
			if resp.SubresourceUris != nil {
				if rec, ok := (*resp.SubresourceUris)["recordings"]; ok {
					result.RecordURL = fmt.Sprintf("%v", rec)
				}
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf("call timed out after 5 minutes (SID: %s)", callSID)
}
