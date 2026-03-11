package quality

import (
	"strings"

	"github.com/ashczar77/mockingjay/internal/flow"
)

// QualityAnalyzer analyzes response quality
type QualityAnalyzer struct{}

// NewQualityAnalyzer creates a new quality analyzer
func NewQualityAnalyzer() *QualityAnalyzer {
	return &QualityAnalyzer{}
}

// QualityMetrics contains response quality measurements
type QualityMetrics struct {
	TotalResponses      int
	AvgResponseLength   float64
	CompletenessScore   float64 // % of responses that seem complete
	SentimentScore      float64 // Positive sentiment (0-100)
	ConfidenceScore     float64 // Response confidence (0-100)
	VagueResponses      int     // Responses with uncertain language
	EmptyResponses      int
}

// Analyze analyzes response quality from conversation flows
func (q *QualityAnalyzer) Analyze(flows []flow.ConversationFlow) QualityMetrics {
	metrics := QualityMetrics{}

	if len(flows) == 0 {
		return metrics
	}

	var totalLength int
	var completeResponses int
	var positiveResponses int
	var confidentResponses int

	for _, f := range flows {
		for _, step := range f.Steps {
			if step.ActualResponse == "" {
				metrics.EmptyResponses++
				continue
			}

			metrics.TotalResponses++
			responseLen := len(step.ActualResponse)
			totalLength += responseLen

			// Check completeness (has punctuation, reasonable length)
			if q.isComplete(step.ActualResponse) {
				completeResponses++
			}

			// Check sentiment (positive words)
			if q.isPositive(step.ActualResponse) {
				positiveResponses++
			}

			// Check confidence (no uncertain language)
			if !q.isVague(step.ActualResponse) {
				confidentResponses++
			} else {
				metrics.VagueResponses++
			}
		}
	}

	if metrics.TotalResponses > 0 {
		metrics.AvgResponseLength = float64(totalLength) / float64(metrics.TotalResponses)
		metrics.CompletenessScore = float64(completeResponses) / float64(metrics.TotalResponses) * 100
		metrics.SentimentScore = float64(positiveResponses) / float64(metrics.TotalResponses) * 100
		metrics.ConfidenceScore = float64(confidentResponses) / float64(metrics.TotalResponses) * 100
	}

	return metrics
}

// isComplete checks if a response seems complete
func (q *QualityAnalyzer) isComplete(response string) bool {
	// Has ending punctuation
	if !strings.HasSuffix(response, ".") && 
	   !strings.HasSuffix(response, "!") && 
	   !strings.HasSuffix(response, "?") {
		return false
	}

	// Reasonable length (at least 10 characters)
	if len(response) < 10 {
		return false
	}

	return true
}

// isPositive checks if response has positive sentiment
func (q *QualityAnalyzer) isPositive(response string) bool {
	lower := strings.ToLower(response)
	
	positiveWords := []string{
		"great", "happy", "help", "sure", "yes", "perfect",
		"excellent", "wonderful", "glad", "pleasure", "welcome",
	}

	for _, word := range positiveWords {
		if strings.Contains(lower, word) {
			return true
		}
	}

	return false
}

// isVague checks if response contains uncertain language
func (q *QualityAnalyzer) isVague(response string) bool {
	lower := strings.ToLower(response)
	
	vagueWords := []string{
		"maybe", "perhaps", "might", "possibly", "not sure",
		"i think", "i guess", "probably", "unclear", "unsure",
	}

	for _, word := range vagueWords {
		if strings.Contains(lower, word) {
			return true
		}
	}

	return false
}
