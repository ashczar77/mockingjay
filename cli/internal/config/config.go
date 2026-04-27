package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the MockingJay configuration
type Config struct {
	Version    int        `yaml:"version"`
	Agent      Agent      `yaml:"agent"`
	Scenarios  []Scenario `yaml:"scenarios"`
	Metrics    []string   `yaml:"metrics"`
	Thresholds Thresholds `yaml:"thresholds"`
	ABTest     *ABTest    `yaml:"ab_test,omitempty"`
}

// Agent configuration
type Agent struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	Phone    string `yaml:"phone,omitempty"`
}

// ABTest configuration for A/B testing two agent variants
type ABTest struct {
	VariantA AgentVariant `yaml:"variant_a"`
	VariantB AgentVariant `yaml:"variant_b"`
}

// AgentVariant is a named agent endpoint for A/B testing
type AgentVariant struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

// Scenario represents a test scenario
type Scenario struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Steps       []Step `yaml:"steps"`
}

// Step in a test scenario
type Step struct {
	Say    string `yaml:"say"`
	Expect string `yaml:"expect"`
}

// Thresholds for test metrics
type Thresholds struct {
	LatencyP95     int `yaml:"latency_p95"`
	TaskCompletion int `yaml:"task_completion"`
	WordErrorRate  int `yaml:"word_error_rate"`
}

// Load reads and parses a config file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// Validate checks if the config is valid
func (c *Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version: %d (expected 1)", c.Version)
	}

	hasAgent := c.Agent.Endpoint != "" || c.Agent.Phone != ""
	hasABTest := c.ABTest != nil && c.ABTest.VariantA.Endpoint != "" && c.ABTest.VariantB.Endpoint != ""
	if !hasAgent && !hasABTest {
		return fmt.Errorf("agent must have either endpoint or phone (or configure ab_test with two variants)")
	}

	if len(c.Scenarios) == 0 {
		return fmt.Errorf("at least one scenario is required")
	}

	for i, scenario := range c.Scenarios {
		if scenario.Name == "" {
			return fmt.Errorf("scenario %d: name is required", i)
		}
		if len(scenario.Steps) == 0 {
			return fmt.Errorf("scenario %s: at least one step is required", scenario.Name)
		}
	}

	return nil
}
