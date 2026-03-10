package flow

import (
	"time"

	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/test"
)

// ConversationFlow represents a conversation path through a scenario
type ConversationFlow struct {
	ScenarioName  string
	TotalSteps    int
	CompletedSteps int
	DropOffPoint  int // -1 if completed, otherwise step number where it failed
	Duration      time.Duration
	Success       bool
	Steps         []StepResult
}

// StepResult represents the result of a single conversation step
type StepResult struct {
	StepNumber    int
	UserInput     string
	ExpectedIntent string
	ActualIntent  string
	Matched       bool
	Latency       time.Duration
	Error         string
}

// Analyzer analyzes conversation flows from test results
type Analyzer struct{}

// NewAnalyzer creates a new flow analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// Analyze converts test results into conversation flows
func (a *Analyzer) Analyze(result test.Result, scenario config.Scenario) ConversationFlow {
	flow := ConversationFlow{
		ScenarioName:   result.Scenario,
		TotalSteps:     result.Metrics.StepsTotal,
		CompletedSteps: result.Metrics.StepsCompleted,
		Duration:       result.Duration,
		Success:        result.Passed,
		DropOffPoint:   -1,
		Steps:          make([]StepResult, 0),
	}

	// If test failed, mark drop-off point
	if !result.Passed {
		flow.DropOffPoint = result.Metrics.StepsCompleted
	}

	// Build step results
	for i, step := range scenario.Steps {
		stepResult := StepResult{
			StepNumber:     i + 1,
			UserInput:      step.Say,
			ExpectedIntent: step.Expect,
			Matched:        i < result.Metrics.StepsCompleted,
			Latency:        result.Metrics.Latency,
		}

		if i == result.Metrics.StepsCompleted && !result.Passed {
			stepResult.Error = result.Error
		}

		flow.Steps = append(flow.Steps, stepResult)
	}

	return flow
}

// AnalyzeMultiple analyzes multiple test results
func (a *Analyzer) AnalyzeMultiple(results []test.Result, scenarios []config.Scenario) []ConversationFlow {
	flows := make([]ConversationFlow, 0)

	for i, result := range results {
		if i < len(scenarios) {
			flow := a.Analyze(result, scenarios[i])
			flows = append(flows, flow)
		}
	}

	return flows
}

// FlowInsights provides aggregate insights from multiple flows
type FlowInsights struct {
	TotalFlows       int
	SuccessfulFlows  int
	FailedFlows      int
	SuccessRate      float64
	AvgStepsCompleted float64
	CommonDropOffPoints map[int]int // step number -> count
	AvgDuration      time.Duration
}

// GenerateInsights creates insights from conversation flows
func (a *Analyzer) GenerateInsights(flows []ConversationFlow) FlowInsights {
	insights := FlowInsights{
		TotalFlows:          len(flows),
		CommonDropOffPoints: make(map[int]int),
	}

	if len(flows) == 0 {
		return insights
	}

	var totalSteps int
	var totalDuration time.Duration

	for _, flow := range flows {
		if flow.Success {
			insights.SuccessfulFlows++
		} else {
			insights.FailedFlows++
			if flow.DropOffPoint >= 0 {
				insights.CommonDropOffPoints[flow.DropOffPoint]++
			}
		}

		totalSteps += flow.CompletedSteps
		totalDuration += flow.Duration
	}

	insights.SuccessRate = float64(insights.SuccessfulFlows) / float64(insights.TotalFlows) * 100
	insights.AvgStepsCompleted = float64(totalSteps) / float64(insights.TotalFlows)
	insights.AvgDuration = totalDuration / time.Duration(insights.TotalFlows)

	return insights
}
