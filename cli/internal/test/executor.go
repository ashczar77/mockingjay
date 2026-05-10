package test

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ashczar77/mockingjay/internal/classifier"
	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/voice"
)

// Result represents a test execution result
type Result struct {
	Scenario string
	Passed   bool
	Duration time.Duration
	Error    string
	Metrics  Metrics
	Steps    []StepDetail
}

// StepDetail contains detailed information about a test step
type StepDetail struct {
	Input          string
	ExpectedIntent string
	ActualIntent   string
	Response       string
	Latency        time.Duration
	Confidence     float64
	Reasoning      string
	Success        bool
}

// Metrics captured during test execution
type Metrics struct {
	Latency        time.Duration
	StepsCompleted int
	StepsTotal     int
}

// Executor runs test scenarios
type Executor struct {
	config     *config.Config
	client     voice.Caller
	classifier *classifier.Classifier
}

// New creates a new test executor
func New(cfg *config.Config) *Executor {
	if cfg.Agent.Endpoint == "" {
		fmt.Fprintln(os.Stderr, "Error: agent.endpoint must be set in mockingjay.yaml")
		os.Exit(1)
	}

	var c *classifier.Classifier
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		c = classifier.New(key)
	}

	return &Executor{
		config:     cfg,
		client:     voice.NewClient(cfg.Agent.Endpoint),
		classifier: c,
	}
}

// Run executes a single scenario
func (e *Executor) Run(scenario config.Scenario) Result {
	start := time.Now()

	result := Result{
		Scenario: scenario.Name,
		Passed:   true,
		Metrics:  Metrics{StepsTotal: len(scenario.Steps)},
		Steps:    make([]StepDetail, 0),
	}

	for i, step := range scenario.Steps {
		stepDetail := StepDetail{
			Input:          step.Say,
			ExpectedIntent: step.Expect,
		}

		resp, latency, err := e.client.Call(step.Say)
		stepDetail.Latency = latency

		if err != nil {
			result.Passed = false
			result.Error = fmt.Sprintf("step %d failed: %v", i+1, err)
			result.Steps = append(result.Steps, stepDetail)
			break
		}

		if !resp.Success {
			result.Passed = false
			result.Error = fmt.Sprintf("step %d: %s", i+1, resp.Error)
			result.Steps = append(result.Steps, stepDetail)
			break
		}

		stepDetail.Response = resp.Text

		// Use LLM classifier if available, otherwise fall back to intent field
		if e.classifier != nil && step.Expect != "" {
			classification, err := e.classifier.Classify(resp.Text, step.Expect)
			if err == nil {
				stepDetail.ActualIntent = classification.Intent
				stepDetail.Confidence = classification.Confidence
				stepDetail.Reasoning = classification.Reasoning
				stepDetail.Success = classification.Confidence >= 0.7
			} else {
				// LLM failed — fall back to intent field comparison
				stepDetail.ActualIntent = resp.Intent
				stepDetail.Success = resp.Intent == step.Expect
			}
		} else {
			// No classifier and no OPENAI_API_KEY — use intent field
			stepDetail.ActualIntent = resp.Intent
			stepDetail.Success = resp.Intent == step.Expect
		}

		if !stepDetail.Success {
			result.Passed = false
			result.Error = fmt.Sprintf("step %d: expected intent '%s', got '%s' (confidence: %.0f%%)",
				i+1, step.Expect, stepDetail.ActualIntent, stepDetail.Confidence*100)
		}

		if latency > result.Metrics.Latency {
			result.Metrics.Latency = latency
		}

		result.Steps = append(result.Steps, stepDetail)
		result.Metrics.StepsCompleted = i + 1

		if !stepDetail.Success {
			break
		}
	}

	result.Duration = time.Since(start)
	return result
}

// RunAll executes multiple scenarios in parallel
func (e *Executor) RunAll(scenarios []config.Scenario) []Result {
	results := make([]Result, len(scenarios))
	var wg sync.WaitGroup

	for i, scenario := range scenarios {
		wg.Add(1)
		go func(idx int, s config.Scenario) {
			defer wg.Done()
			results[idx] = e.Run(s)
		}(i, scenario)
	}

	wg.Wait()
	return results
}

// Stats calculates aggregate statistics from results
type Stats struct {
	TotalTests     int
	PassedTests    int
	FailedTests    int
	PassRate       float64
	AvgLatency     int64
	P95Latency     int64
	P99Latency     int64
	TaskCompletion float64
	AvgConfidence  float64
}

// CalculateStats computes statistics from test results
func CalculateStats(results []Result) Stats {
	stats := Stats{TotalTests: len(results)}
	if len(results) == 0 {
		return stats
	}

	var latencies []int64
	var totalLatency, totalConfidence int64
	var confidenceCount int
	var totalSteps, completedSteps int

	for _, r := range results {
		if r.Passed {
			stats.PassedTests++
		} else {
			stats.FailedTests++
		}

		latencyMs := r.Metrics.Latency.Milliseconds()
		latencies = append(latencies, latencyMs)
		totalLatency += latencyMs
		totalSteps += r.Metrics.StepsTotal
		completedSteps += r.Metrics.StepsCompleted

		for _, s := range r.Steps {
			if s.Confidence > 0 {
				totalConfidence += int64(s.Confidence * 100)
				confidenceCount++
			}
		}
	}

	stats.PassRate = float64(stats.PassedTests) / float64(stats.TotalTests) * 100
	stats.AvgLatency = totalLatency / int64(len(results))

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	p95 := int(float64(len(latencies)) * 0.95)
	if p95 >= len(latencies) {
		p95 = len(latencies) - 1
	}
	stats.P95Latency = latencies[p95]

	p99 := int(float64(len(latencies)) * 0.99)
	if p99 >= len(latencies) {
		p99 = len(latencies) - 1
	}
	stats.P99Latency = latencies[p99]

	if totalSteps > 0 {
		stats.TaskCompletion = float64(completedSteps) / float64(totalSteps) * 100
	}
	if confidenceCount > 0 {
		stats.AvgConfidence = float64(totalConfidence) / float64(confidenceCount)
	}

	return stats
}
