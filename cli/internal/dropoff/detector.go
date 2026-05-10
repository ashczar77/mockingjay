package dropoff

import (
	"sort"

	"github.com/ashczar77/mockingjay/internal/flow"
)

// Detector identifies conversation failure points
type Detector struct{}

// NewDetector creates a new failure point detector
func NewDetector() *Detector {
	return &Detector{}
}

// FailurePoint represents a step where conversations consistently fail
type FailurePoint struct {
	StepNumber  int
	StepInput   string
	Frequency   int
	FailureRate float64
	Severity    string // "critical", "high", "medium", "low"
}

// DropOffPoint is an alias for FailurePoint for backwards compatibility
type DropOffPoint = FailurePoint

// FailureAnalysis contains failure point detection results
type FailureAnalysis struct {
	TotalConversations int
	FailurePoints      []FailurePoint
	CriticalPoints     []FailurePoint
	OverallFailureRate float64
}

// DropOffAnalysis is an alias for FailureAnalysis for backwards compatibility
type DropOffAnalysis = FailureAnalysis

// Analyze detects failure points from conversation flows
func (d *Detector) Analyze(flows []flow.ConversationFlow) FailureAnalysis {
	analysis := FailureAnalysis{
		TotalConversations: len(flows),
	}

	if len(flows) == 0 {
		return analysis
	}

	stepFailures := make(map[int]map[string]int)
	stepTotals := make(map[int]int)

	for _, f := range flows {
		for _, step := range f.Steps {
			stepNum := step.StepNumber
			if stepFailures[stepNum] == nil {
				stepFailures[stepNum] = make(map[string]int)
			}
			stepTotals[stepNum]++
			if !step.Matched {
				stepFailures[stepNum][step.UserInput]++
			}
		}
	}

	var failures []FailurePoint
	totalFailures := 0

	for stepNum, stepMap := range stepFailures {
		for input, count := range stepMap {
			if count > 0 {
				totalFailures += count
				rate := float64(count) / float64(stepTotals[stepNum]) * 100
				failures = append(failures, FailurePoint{
					StepNumber:  stepNum,
					StepInput:   input,
					Frequency:   count,
					FailureRate: rate,
					Severity:    d.calculateSeverity(rate),
				})
			}
		}
	}

	sort.Slice(failures, func(i, j int) bool {
		return failures[i].Frequency > failures[j].Frequency
	})

	analysis.FailurePoints = failures
	analysis.OverallFailureRate = float64(totalFailures) / float64(analysis.TotalConversations) * 100

	for _, p := range failures {
		if p.Severity == "critical" || p.Severity == "high" {
			analysis.CriticalPoints = append(analysis.CriticalPoints, p)
		}
	}

	return analysis
}

func (d *Detector) calculateSeverity(rate float64) string {
	switch {
	case rate >= 50:
		return "critical"
	case rate >= 25:
		return "high"
	case rate >= 10:
		return "medium"
	default:
		return "low"
	}
}
