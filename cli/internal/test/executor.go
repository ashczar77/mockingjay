package test

import (
	"fmt"
	"sort"
	"sync"
	"time"

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
	config *config.Config
	client *voice.Client
}

// New creates a new test executor
func New(cfg *config.Config) *Executor {
	var client *voice.Client
	if cfg.Agent.Endpoint != "" {
		client = voice.NewClient(cfg.Agent.Endpoint)
	}
	
	return &Executor{
		config: cfg,
		client: client,
	}
}

// Run executes a single scenario
func (e *Executor) Run(scenario config.Scenario) Result {
	start := time.Now()
	
	result := Result{
		Scenario: scenario.Name,
		Passed:   true,
		Metrics: Metrics{
			StepsTotal: len(scenario.Steps),
		},
		Steps: make([]StepDetail, 0),
	}

	// Execute each step
	for i, step := range scenario.Steps {
		stepDetail := StepDetail{
			Input:          step.Say,
			ExpectedIntent: step.Expect,
			Success:        false,
		}

		if e.client != nil {
			// Actually call the agent
			resp, latency, err := e.client.Call(step.Say)
			if err != nil {
				result.Passed = false
				result.Error = fmt.Sprintf("step %d failed: %v", i+1, err)
				stepDetail.Latency = latency
				result.Steps = append(result.Steps, stepDetail)
				break
			}

			if !resp.Success {
				result.Passed = false
				result.Error = fmt.Sprintf("step %d: %s", i+1, resp.Error)
				stepDetail.Latency = latency
				result.Steps = append(result.Steps, stepDetail)
				break
			}

			stepDetail.ActualIntent = resp.Intent
			stepDetail.Response = resp.Text
			stepDetail.Latency = latency
			stepDetail.Success = true

			// Check if intent matches expectation
			if resp.Intent != step.Expect {
				result.Passed = false
				result.Error = fmt.Sprintf("step %d: expected intent '%s', got '%s'", i+1, step.Expect, resp.Intent)
				stepDetail.Success = false
			}

			if latency > result.Metrics.Latency {
				result.Metrics.Latency = latency
			}
		} else {
			// Simulate call if no endpoint configured
			time.Sleep(100 * time.Millisecond)
			result.Metrics.Latency = 100 * time.Millisecond
			stepDetail.ActualIntent = step.Expect
			stepDetail.Response = "Simulated response"
			stepDetail.Latency = 100 * time.Millisecond
			stepDetail.Success = true
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
	TotalTests      int
	PassedTests     int
	FailedTests     int
	PassRate        float64
	AvgLatency      int64
	P95Latency      int64
	P99Latency      int64
	TaskCompletion  float64
}

// CalculateStats computes statistics from test results
func CalculateStats(results []Result) Stats {
	stats := Stats{
		TotalTests: len(results),
	}

	if len(results) == 0 {
		return stats
	}

	var latencies []int64
	var totalLatency int64
	var totalSteps int
	var completedSteps int

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
	}

	// Pass rate
	stats.PassRate = float64(stats.PassedTests) / float64(stats.TotalTests) * 100

	// Average latency
	stats.AvgLatency = totalLatency / int64(len(results))

	// Percentiles
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	p95Index := int(float64(len(latencies)) * 0.95)
	if p95Index >= len(latencies) {
		p95Index = len(latencies) - 1
	}
	stats.P95Latency = latencies[p95Index]

	p99Index := int(float64(len(latencies)) * 0.99)
	if p99Index >= len(latencies) {
		p99Index = len(latencies) - 1
	}
	stats.P99Latency = latencies[p99Index]

	// Task completion
	if totalSteps > 0 {
		stats.TaskCompletion = float64(completedSteps) / float64(totalSteps) * 100
	}

	return stats
}
