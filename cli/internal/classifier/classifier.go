package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Result holds the classification output
type Result struct {
	Intent     string
	Confidence float64 // 0.0 - 1.0
	Reasoning  string
}

// Provider is the interface for any LLM backend
type Provider interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

// Classifier uses a Provider to classify voice AI responses
type Classifier struct {
	provider Provider
}

// New creates a Classifier from environment variables.
// Checks LLM_PROVIDER (openai|anthropic|ollama), defaults to openai.
// Falls back gracefully if no API key is set.
func New() *Classifier {
	providerName := strings.ToLower(os.Getenv("LLM_PROVIDER"))
	if providerName == "" {
		providerName = "openai"
	}

	switch providerName {
	case "anthropic":
		key := os.Getenv("ANTHROPIC_API_KEY")
		if key == "" {
			return nil
		}
		return &Classifier{provider: newAnthropicProvider(key)}
	case "ollama":
		model := os.Getenv("OLLAMA_MODEL")
		if model == "" {
			model = "llama3"
		}
		host := os.Getenv("OLLAMA_HOST")
		if host == "" {
			host = "http://localhost:11434"
		}
		return &Classifier{provider: newOllamaProvider(host, model)}
	default: // openai
		key := os.Getenv("OPENAI_API_KEY")
		if key == "" {
			return nil
		}
		return &Classifier{provider: newOpenAIProvider(key)}
	}
}

// NewWithProvider creates a Classifier with an explicit provider (useful for testing)
func NewWithProvider(p Provider) *Classifier {
	return &Classifier{provider: p}
}

// Classify determines whether an agent's text response matches an expected intent
func (c *Classifier) Classify(agentResponse, expectedIntent string) (*Result, error) {
	prompt := fmt.Sprintf(`You are evaluating a voice AI agent's response.

Agent response: "%s"
Expected intent: "%s"

Does the agent's response match the expected intent?
Reply with JSON only, no markdown:
{"intent": "<the actual intent of the response>", "matches": <true|false>, "confidence": <0.0-1.0>, "reasoning": "<one sentence>"}`,
		agentResponse, expectedIntent)

	raw, err := c.provider.Complete(context.Background(), prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM classification failed: %w", err)
	}
	return parseResult(strings.TrimSpace(raw), expectedIntent)
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

	raw, err := c.provider.Complete(context.Background(), prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM classification failed: %w", err)
	}
	return parseResult(strings.TrimSpace(raw), expectedIntent)
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

	raw, err := c.provider.Complete(context.Background(), prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM quality evaluation failed: %w", err)
	}
	return parseQuality(strings.TrimSpace(raw))
}

// --- OpenAI Provider ---

type openAIProvider struct {
	client openai.Client
}

func newOpenAIProvider(apiKey string) *openAIProvider {
	return &openAIProvider{client: openai.NewClient(option.WithAPIKey(apiKey))}
}

func (p *openAIProvider) Complete(ctx context.Context, prompt string) (string, error) {
	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       openai.ChatModelGPT4oMini,
		Messages:    []openai.ChatCompletionMessageParamUnion{openai.UserMessage(prompt)},
		Temperature: openai.Float(0.0),
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// --- Anthropic Provider ---

type anthropicProvider struct {
	client anthropic.Client
}

func newAnthropicProvider(apiKey string) *anthropicProvider {
	return &anthropicProvider{client: anthropic.NewClient(anthropicoption.WithAPIKey(apiKey))}
}

func (p *anthropicProvider) Complete(ctx context.Context, prompt string) (string, error) {
	resp, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 256,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic")
	}
	return resp.Content[0].Text, nil
}

// --- Ollama Provider (local models) ---

type ollamaProvider struct {
	host  string
	model string
}

func newOllamaProvider(host, model string) *ollamaProvider {
	return &ollamaProvider{host: host, model: model}
}

func (p *ollamaProvider) Complete(ctx context.Context, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model":  p.model,
		"prompt": prompt,
		"stream": false,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", p.host+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to parse ollama response: %w", err)
	}
	return result.Response, nil
}

// --- Parsing helpers ---

func parseResult(raw, expectedIntent string) (*Result, error) {
	intent := extractString(raw, "intent")
	if intent == "" {
		intent = expectedIntent
	}
	matches := strings.Contains(raw, `"matches": true`) || strings.Contains(raw, `"matches":true`)
	confidence := extractFloat(raw, "confidence")
	reasoning := extractString(raw, "reasoning")

	result := &Result{Intent: intent, Confidence: confidence, Reasoning: reasoning}
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

func extractString(j, key string) string {
	for _, sep := range []string{`": "`, `":"`} {
		search := fmt.Sprintf(`"%s%s`, key, sep)
		if idx := strings.Index(j, search); idx != -1 {
			start := idx + len(search)
			if end := strings.Index(j[start:], `"`); end != -1 {
				return j[start : start+end]
			}
		}
	}
	return ""
}

func extractFloat(j, key string) float64 {
	for _, sep := range []string{`": `, `":`} {
		search := fmt.Sprintf(`"%s%s`, key, sep)
		if idx := strings.Index(j, search); idx != -1 {
			var val float64
			fmt.Sscanf(j[idx+len(search):], "%f", &val)
			return val
		}
	}
	return 0
}
