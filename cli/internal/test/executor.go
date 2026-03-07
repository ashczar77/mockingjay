package test

import (
	"fmt"
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
	}

	// Execute each step
	for i, step := range scenario.Steps {
		if e.client != nil {
			// Actually call the agent
			resp, latency, err := e.client.Call(step.Say)
			if err != nil {
				result.Passed = false
				result.Error = fmt.Sprintf("step %d failed: %v", i+1, err)
				break
			}

			if !resp.Success {
				result.Passed = false
				result.Error = fmt.Sprintf("step %d: %s", i+1, resp.Error)
				break
			}

			if latency > result.Metrics.Latency {
				result.Metrics.Latency = latency
			}

			// TODO: Validate response matches expectation
			_ = step.Expect
		} else {
			// Simulate call if no endpoint configured
			time.Sleep(100 * time.Millisecond)
			result.Metrics.Latency = 100 * time.Millisecond
		}

		result.Metrics.StepsCompleted = i + 1
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
