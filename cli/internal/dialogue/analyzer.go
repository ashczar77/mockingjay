package dialogue

import (
	"github.com/ashczar77/mockingjay/internal/flow"
)

// DialogueAnalyzer analyzes multi-turn conversations
type DialogueAnalyzer struct{}

// NewDialogueAnalyzer creates a new dialogue analyzer
func NewDialogueAnalyzer() *DialogueAnalyzer {
	return &DialogueAnalyzer{}
}

// DialogueMetrics contains metrics for multi-turn conversations
type DialogueMetrics struct {
	TotalConversations   int
	MultiTurnCount       int     // Conversations with 2+ turns
	AvgTurnsPerConv      float64
	MaxTurns             int
	ContextRetention     float64 // % of conversations that maintained context
	CoherenceScore       float64 // Overall conversation coherence
}

// Analyze analyzes dialogue patterns from conversation flows
func (d *DialogueAnalyzer) Analyze(flows []flow.ConversationFlow) DialogueMetrics {
	metrics := DialogueMetrics{
		TotalConversations: len(flows),
	}

	if len(flows) == 0 {
		return metrics
	}

	var totalTurns int
	var contextRetained int

	for _, f := range flows {
		turns := len(f.Steps)
		totalTurns += turns

		if turns > metrics.MaxTurns {
			metrics.MaxTurns = turns
		}

		if turns >= 2 {
			metrics.MultiTurnCount++
			
			// Check context retention: did all steps succeed?
			if f.Success {
				contextRetained++
			}
		}
	}

	metrics.AvgTurnsPerConv = float64(totalTurns) / float64(len(flows))

	if metrics.MultiTurnCount > 0 {
		metrics.ContextRetention = float64(contextRetained) / float64(metrics.MultiTurnCount) * 100
	}

	// Calculate coherence score based on successful multi-turn conversations
	if metrics.MultiTurnCount > 0 {
		metrics.CoherenceScore = float64(contextRetained) / float64(metrics.MultiTurnCount) * 100
	} else {
		metrics.CoherenceScore = 100.0 // Single-turn conversations are always coherent
	}

	return metrics
}

// ContextLossPoint identifies where context was lost in a conversation
type ContextLossPoint struct {
	ScenarioName string
	StepNumber   int
	UserInput    string
	Expected     string
	Actual       string
	Reason       string
}

// DetectContextLoss finds points where conversation context was lost
func (d *DialogueAnalyzer) DetectContextLoss(flows []flow.ConversationFlow) []ContextLossPoint {
	lossPoints := make([]ContextLossPoint, 0)

	for _, f := range flows {
		if !f.Success && len(f.Steps) >= 2 {
			// Find the step where context was lost
			for i, step := range f.Steps {
				if !step.Matched && i > 0 {
					lossPoint := ContextLossPoint{
						ScenarioName: f.ScenarioName,
						StepNumber:   step.StepNumber,
						UserInput:    step.UserInput,
						Expected:     step.ExpectedIntent,
						Actual:       step.ActualIntent,
						Reason:       "Intent mismatch in multi-turn conversation",
					}
					lossPoints = append(lossPoints, lossPoint)
					break
				}
			}
		}
	}

	return lossPoints
}
