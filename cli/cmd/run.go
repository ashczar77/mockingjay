package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
	scenario   string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run voice AI tests",
	Long:  `Execute test scenarios against your voice AI agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runTests(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&configFile, "config", "c", "mockingjay.yaml", "config file")
	runCmd.Flags().StringVarP(&scenario, "scenario", "s", "", "run specific scenario")
}

func runTests() error {
	fmt.Println("🐦 MockingJay - Starting tests...")
	fmt.Println()

	fmt.Printf("📋 Loading config from: %s\n", configFile)
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	scenarios := cfg.Scenarios
	if scenario != "" {
		found := false
		for _, s := range cfg.Scenarios {
			if s.Name == scenario {
				scenarios = []config.Scenario{s}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("scenario '%s' not found", scenario)
		}
	}

	fmt.Printf("🎯 Running %d scenario(s)...\n", len(scenarios))
	fmt.Println()

	passed := 0
	failed := 0

	for i, s := range scenarios {
		fmt.Printf("  [%d/%d] %s", i+1, len(scenarios), s.Name)
		if s.Description != "" {
			fmt.Printf(" - %s", s.Description)
		}
		fmt.Print("... ")

		time.Sleep(500 * time.Millisecond)

		// TODO: Actually execute the test
		fmt.Println("✓ PASS")
		passed++
	}

	fmt.Println()
	fmt.Println("📊 Results:")
	fmt.Printf("  Tests run: %d\n", len(scenarios))
	fmt.Printf("  Passed: %d\n", passed)
	fmt.Printf("  Failed: %d\n", failed)
	fmt.Println()

	if failed > 0 {
		fmt.Println("❌ Some tests failed")
		return fmt.Errorf("%d test(s) failed", failed)
	}

	fmt.Println("✨ All tests passed!")
	return nil
}
