package config

// Config represents the MockingJay configuration
type Config struct {
	Version   int       `yaml:"version"`
	Agent     Agent     `yaml:"agent"`
	Scenarios []Scenario `yaml:"scenarios"`
	Metrics   []string  `yaml:"metrics"`
	Thresholds Thresholds `yaml:"thresholds"`
}

// Agent configuration
type Agent struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	Phone    string `yaml:"phone,omitempty"`
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
	LatencyP95       int `yaml:"latency_p95"`
	TaskCompletion   int `yaml:"task_completion"`
	WordErrorRate    int `yaml:"word_error_rate"`
}
