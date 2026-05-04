package models

import "time"

type TestResult struct {
	ID        int64     `json:"id"`
	Scenario  string    `json:"scenario"`
	Passed    bool      `json:"passed"`
	Latency   int64     `json:"latency_ms"`
	Error     string    `json:"error,omitempty"`
	Variant   string    `json:"variant,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ABTestResult struct {
	ID        int64     `json:"id"`
	TestName  string    `json:"test_name"`
	VariantA  string    `json:"variant_a"`
	VariantB  string    `json:"variant_b"`
	Winner    string    `json:"winner"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

type Transcription struct {
	ID         int64     `json:"id"`
	CallSID    string    `json:"call_sid,omitempty"`
	AudioPath  string    `json:"audio_path"`
	Text       string    `json:"text"`
	Confidence float64   `json:"confidence"`
	Duration   float64   `json:"duration_seconds"`
	CreatedAt  time.Time `json:"created_at"`
}

type ConversationMetrics struct {
	SuccessRate       float64 `json:"success_rate"`
	IntentAccuracy    float64 `json:"intent_accuracy"`
	AvgStepsCompleted float64 `json:"avg_steps_completed"`
	MultiTurnCount    int     `json:"multi_turn_count"`
	ContextRetention  float64 `json:"context_retention"`
	CoherenceScore    float64 `json:"coherence_score"`
	CompletenessScore float64 `json:"completeness_score"`
	SentimentScore    float64 `json:"sentiment_score"`
	ConfidenceScore   float64 `json:"confidence_score"`
	AvgResponseLength float64 `json:"avg_response_length"`
	TotalTests        int     `json:"total_tests"`
	PassedTests       int     `json:"passed_tests"`
	AvgLatency        float64 `json:"avg_latency_ms"`
}

type HealthStatus struct {
	Status     string       `json:"status"`
	PassRate24h float64     `json:"pass_rate_24h"`
	Total24h   int          `json:"total_24h"`
	Passed24h  int          `json:"passed_24h"`
	RecentRuns []RecentRun  `json:"recent_runs"`
}

type RecentRun struct {
	Scenario  string `json:"scenario"`
	Passed    bool   `json:"passed"`
	LatencyMs int64  `json:"latency_ms"`
	CreatedAt string `json:"created_at"`
}
