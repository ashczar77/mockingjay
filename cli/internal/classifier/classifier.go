package classifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Result holds the classification output
type Result struct {
	Intent     string
	Confidence float64 // 0.0 - 1.0
	Reasoning  string
}

// Classifier classifies the intent of a text response using an LLM
type Classifier struct {
	client openai.Client
}

// New creates a new Classifier using the provided OpenAI API key
func New(apiKey string) *Classifier {
	return &Classifier{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

// Classify determines whether an agent's text response matches an expected intent.
// It returns the inferred intent and a confidence score.
func (c *Classifier) Classify(agentResponse, expectedIntent string) (*Result, error) {
	prompt := fmt.Sprintf(`You are evaluating a voice AI agent's response.

Agent response: "%s"
Expected intent: "%s"

Does the agent's response match the expected intent? 
Reply with JSON only, no markdown:
{"intent": "<the actual intent of the response>", "matches": <true|false>, "confidence": <0.0-1.0>, "reasoning": "<one sentence>"}`,
		agentResponse, expectedIntent)

	resp, err := c.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.0),
	})
	if err != nil {
		return nil, fmt.Errorf("LLM classification failed: %w", err)
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	return parseResult(raw, expectedIntent)
}

// ClassifyTranscript determines the intent of a call transcript against an expected intent
func (c *Classifier) ClassifyTranscript(transcript, expectedIntent string) (*Result, error) {
	prompt := fmt.Sprintf(`You are evaluating a voice AI call transcript.

Transcript: "%s"
Expected intent: "%s"

Does the transcript indicate the agent fulfilled the expected intent?
Reply with JSON only, no markdown:
{"intent": "<the actual intent expressed>", "matches": <true|false>, "confidence": <0.0-1.0>, "reasoning": "<one sentence>"}`,
		transcript, expectedIntent)

	resp, err := c.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.0),
	})
	if err != nil {
		return nil, fmt.Errorf("LLM classification failed: %w", err)
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	return parseResult(raw, expectedIntent)
}

// EvaluateQuality scores the quality of an agent response
func (c *Classifier) EvaluateQuality(agentResponse string) (map[string]float64, error) {
	prompt := fmt.Sprintf(`Rate this voice AI agent response on three dimensions.

Response: "%s"

Reply with JSON only, no markdown:
{"completeness": <0.0-1.0>, "sentiment": <0.0-1.0>, "confidence": <0.0-1.0>}

- completeness: does the response fully address what was asked?
- sentiment: how positive/helpful is the tone? (0=negative, 1=positive)
- confidence: does the agent sound certain and authoritative?`, agentResponse)

	resp, err := c.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.0),
	})
	if err != nil {
		return nil, fmt.Errorf("LLM quality evaluation failed: %w", err)
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	return parseQuality(raw)
}

func parseResult(raw, expectedIntent string) (*Result, error) {
	// Simple JSON field extraction without importing encoding/json to keep it lean
	intent := extractString(raw, "intent")
	if intent == "" {
		intent = expectedIntent
	}
	matches := strings.Contains(raw, `"matches": true`) || strings.Contains(raw, `"matches":true`)
	confidence := extractFloat(raw, "confidence")
	reasoning := extractString(raw, "reasoning")

	result := &Result{
		Intent:     intent,
		Confidence: confidence,
		Reasoning:  reasoning,
	}
	if !matches {
		result.Confidence = 1.0 - confidence
	}
	return result, nil
}

func parseQuality(raw string) (map[string]float64, error) {
	return map[string]float64{
		"completeness": extractFloat(raw, "completeness"),
		"sentiment":    extractFloat(raw, "sentiment"),
		"confidence":   extractFloat(raw, "confidence"),
	}, nil
}

func extractString(json, key string) string {
	search := fmt.Sprintf(`"%s": "`, key)
	idx := strings.Index(json, search)
	if idx == -1 {
		search = fmt.Sprintf(`"%s":"`, key)
		idx = strings.Index(json, search)
		if idx == -1 {
			return ""
		}
	}
	start := idx + len(search)
	end := strings.Index(json[start:], `"`)
	if end == -1 {
		return ""
	}
	return json[start : start+end]
}

func extractFloat(json, key string) float64 {
	search := fmt.Sprintf(`"%s": `, key)
	idx := strings.Index(json, search)
	if idx == -1 {
		search = fmt.Sprintf(`"%s":`, key)
		idx = strings.Index(json, search)
		if idx == -1 {
			return 0
		}
	}
	start := idx + len(search)
	var val float64
	fmt.Sscanf(json[start:], "%f", &val)
	return val
}
