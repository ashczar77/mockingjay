package dropoff

import (
	"sort"

	"github.com/ashczar77/mockingjay/internal/flow"
)

// Detector identifies conversation drop-off patterns
type Detector struct{}

// NewDetector creates a new drop-off detector
func NewDetector() *Detector {
	return &Detector{}
}

// DropOffPoint represents a point where users abandon conversations
type DropOffPoint struct {
	StepNumber  int
	StepInput   string
	Frequency   int
	DropOffRate float64
	Severity    string // "critical", "high", "medium", "low"
}

// DropOffAnalysis contains drop-off detection results
type DropOffAnalysis struct {
	TotalConversations int
	DropOffPoints      []DropOffPoint
	CriticalPoints     []DropOffPoint
	OverallDropOffRate float64
}

// Analyze detects drop-off points from conversation flows
func (d *Detector) Analyze(flows []flow.ConversationFlow) DropOffAnalysis {
	analysis := DropOffAnalysis{
		TotalConversations: len(flows),
	}

	if len(flows) == 0 {
		return analysis
	}

	// Track failures at each step
	stepFailures := make(map[int]map[string]int) // step -> input -> count
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

	// Calculate drop-off rates
	dropOffs := []DropOffPoint{}
	totalDropOffs := 0

	for stepNum, failures := range stepFailures {
		for input, count := range failures {
			if count > 0 {
				totalDropOffs += count
				rate := float64(count) / float64(stepTotals[stepNum]) * 100

				severity := d.calculateSeverity(rate)

				dropOffs = append(dropOffs, DropOffPoint{
					StepNumber:  stepNum,
					StepInput:   input,
					Frequency:   count,
					DropOffRate: rate,
					Severity:    severity,
				})
			}
		}
	}

	// Sort by frequency (most common first)
	sort.Slice(dropOffs, func(i, j int) bool {
		return dropOffs[i].Frequency > dropOffs[j].Frequency
	})

	analysis.DropOffPoints = dropOffs
	analysis.OverallDropOffRate = float64(totalDropOffs) / float64(analysis.TotalConversations) * 100

	// Extract critical points
	for _, point := range dropOffs {
		if point.Severity == "critical" || point.Severity == "high" {
			analysis.CriticalPoints = append(analysis.CriticalPoints, point)
		}
	}

	return analysis
}

// calculateSeverity determines severity based on drop-off rate
func (d *Detector) calculateSeverity(rate float64) string {
	if rate >= 50 {
		return "critical"
	} else if rate >= 25 {
		return "high"
	} else if rate >= 10 {
		return "medium"
	}
	return "low"
}
