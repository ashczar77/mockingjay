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
	// Append ngrok browser-warning bypass so Twilio receives TwiML directly
	webhookURL := c.cfg.WebhookURL
	if len(webhookURL) > 0 {
		sep := "?"
		for _, ch := range webhookURL {
			if ch == '?' {
				sep = "&"
				break
			}
		}
		webhookURL += sep + "ngrok-skip-browser-warning=true"
	}
	params.SetUrl(webhookURL)
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

// GetRecordingURL fetches the first recording URL for a completed call
func GetRecordingURL(accountSID, authToken, callSID string) (string, error) {
	c := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	resp, err := c.Api.ListRecording(&twilioApi.ListRecordingParams{
		CallSid: &callSID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list recordings: %w", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("no recordings found for call %s", callSID)
	}
	sid := *resp[0].Sid
	return fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Recordings/%s.wav", accountSID, sid), nil
}
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
