package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new MockingJay project",
	Long:  `Creates a mockingjay.yaml configuration file in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initProject(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func initProject() error {
	configTemplate := `# MockingJay Configuration
version: 1

# Voice AI agent to test
agent:
  endpoint: "https://your-voice-agent.com/call"
  # or use phone number for telephony testing
  # phone: "+1234567890"

# Test scenarios
scenarios:
  - name: "basic-greeting"
    description: "Test basic greeting flow"
    steps:
      - say: "Hello"
        expect: "greeting"
    
  - name: "appointment-booking"
    description: "Test appointment booking"
    steps:
      - say: "I want to book an appointment"
        expect: "booking_intent"
      - say: "Tomorrow at 2pm"
        expect: "confirmation"

# Metrics to track
metrics:
  - latency
  - task_completion
  - word_error_rate

# Thresholds
thresholds:
  latency_p95: 5000  # milliseconds
  task_completion: 85  # percentage
  word_error_rate: 10  # percentage
`

	filename := "mockingjay.yaml"
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("mockingjay.yaml already exists")
	}

	if err := os.WriteFile(filename, []byte(configTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Println("✓ Created mockingjay.yaml")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit mockingjay.yaml with your agent details")
	fmt.Println("  2. Run: mockingjay run")
	
	return nil
}
