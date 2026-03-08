package cmd

import (
	"fmt"
	"os"

	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/reporter"
	"github.com/ashczar77/mockingjay/internal/test"
	"github.com/spf13/cobra"
)

var (
	configFile string
	scenario   string
	apiURL     string
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
	runCmd.Flags().StringVar(&apiURL, "api-url", "", "backend API URL (optional)")
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

	executor := test.New(cfg)
	results := executor.RunAll(scenarios)

	// Send results to backend if API URL provided
	var reporterClient *reporter.Client
	if apiURL != "" {
		reporterClient = reporter.NewClient(apiURL)
		fmt.Printf("📤 Sending results to: %s\n", apiURL)
	}

	passed := 0
	failed := 0
	var totalLatency int64

	for i, r := range results {
		fmt.Printf("  [%d/%d] %s", i+1, len(results), r.Scenario)
		if r.Passed {
			fmt.Printf(" ✓ PASS (latency: %dms)\n", r.Metrics.Latency.Milliseconds())
			passed++
			totalLatency += r.Metrics.Latency.Milliseconds()
		} else {
			fmt.Printf(" ✗ FAIL (%s)\n", r.Error)
			failed++
		}

		// Send to backend
		if reporterClient != nil {
			if err := reporterClient.Report(r); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to report result: %v\n", err)
			}
		}
	}

	fmt.Println()
	fmt.Println("📊 Results:")
	fmt.Printf("  Tests run: %d\n", len(results))
	fmt.Printf("  Passed: %d\n", passed)
	fmt.Printf("  Failed: %d\n", failed)
	if passed > 0 {
		fmt.Printf("  Avg latency: %dms\n", totalLatency/int64(passed))
	}
	fmt.Println()

	if failed > 0 {
		fmt.Println("❌ Some tests failed")
		return fmt.Errorf("%d test(s) failed", failed)
	}

	fmt.Println("✨ All tests passed!")
	return nil
}
