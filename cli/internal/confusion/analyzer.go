package confusion

import (
	"sort"

	"github.com/ashczar77/mockingjay/internal/flow"
)

// Analyzer identifies confusion patterns in conversations
type Analyzer struct{}

// NewAnalyzer creates a new confusion analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// ConfusionPattern represents a misunderstanding pattern
type ConfusionPattern struct {
	UserInput      string
	ExpectedIntent string
	ActualIntent   string
	Frequency      int
	ConfusionRate  float64
}

// ConfusionAnalysis contains confusion detection results
type ConfusionAnalysis struct {
	TotalIntentChecks  int
	ConfusedIntents    int
	ConfusionRate      float64
	Patterns           []ConfusionPattern
	MostConfusedInputs []ConfusionPattern
}

// Analyze detects confusion patterns from conversation flows
func (a *Analyzer) Analyze(flows []flow.ConversationFlow) ConfusionAnalysis {
	analysis := ConfusionAnalysis{}

	if len(flows) == 0 {
		return analysis
	}

	// Track intent mismatches
	confusionMap := make(map[string]map[string]map[string]int) // input -> expected -> actual -> count
	inputTotals := make(map[string]int)

	for _, f := range flows {
		for _, step := range f.Steps {
			if step.ExpectedIntent == "" {
				continue
			}

			analysis.TotalIntentChecks++
			inputTotals[step.UserInput]++

			if !step.Matched {
				analysis.ConfusedIntents++

				if confusionMap[step.UserInput] == nil {
					confusionMap[step.UserInput] = make(map[string]map[string]int)
				}
				if confusionMap[step.UserInput][step.ExpectedIntent] == nil {
					confusionMap[step.UserInput][step.ExpectedIntent] = make(map[string]int)
				}
				confusionMap[step.UserInput][step.ExpectedIntent][step.ActualIntent]++
			}
		}
	}

	// Calculate confusion rate
	if analysis.TotalIntentChecks > 0 {
		analysis.ConfusionRate = float64(analysis.ConfusedIntents) / float64(analysis.TotalIntentChecks) * 100
	}

	// Build patterns
	patterns := []ConfusionPattern{}
	for input, expectedMap := range confusionMap {
		for expected, actualMap := range expectedMap {
			for actual, count := range actualMap {
				rate := float64(count) / float64(inputTotals[input]) * 100

				patterns = append(patterns, ConfusionPattern{
					UserInput:      input,
					ExpectedIntent: expected,
					ActualIntent:   actual,
					Frequency:      count,
					ConfusionRate:  rate,
				})
			}
		}
	}

	// Sort by frequency
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency > patterns[j].Frequency
	})

	analysis.Patterns = patterns

	// Extract top 5 most confused inputs
	if len(patterns) > 5 {
		analysis.MostConfusedInputs = patterns[:5]
	} else {
		analysis.MostConfusedInputs = patterns
	}

	return analysis
}
