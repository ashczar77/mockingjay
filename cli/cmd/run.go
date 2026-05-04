package cmd

import (
	"fmt"
	"os"

	"github.com/ashczar77/mockingjay/internal/confusion"
	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/dialogue"
	"github.com/ashczar77/mockingjay/internal/dropoff"
	"github.com/ashczar77/mockingjay/internal/flow"
	"github.com/ashczar77/mockingjay/internal/printer"
	"github.com/ashczar77/mockingjay/internal/quality"
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
		scenarios, err = filterScenario(cfg.Scenarios, scenario)
		if err != nil {
			return err
		}
	}

	fmt.Printf("🎯 Running %d scenario(s)...\n\n", len(scenarios))

	executor := test.New(cfg)
	results := executor.RunAll(scenarios)

	p := printer.New()
	p.Results(results)

	if apiURL != "" {
		fmt.Printf("\n📤 Sending results to: %s\n", apiURL)
		r := reporter.NewClient(apiURL)
		for _, res := range results {
			if err := r.Report(res); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to report result: %v\n", err)
			}
		}
	}

	stats := test.CalculateStats(results)
	p.Stats(stats)

	flowAnalyzer := flow.NewAnalyzer()
	flows := flowAnalyzer.AnalyzeMultiple(results, scenarios)
	p.ConversationInsights(flowAnalyzer.GenerateInsights(flows))

	dialogueAnalyzer := dialogue.NewDialogueAnalyzer()
	p.DialogueMetrics(dialogueAnalyzer.Analyze(flows), dialogueAnalyzer.DetectContextLoss(flows))

	p.QualityMetrics(quality.NewQualityAnalyzer().Analyze(flows))
	p.DropOffAnalysis(dropoff.NewDetector().Analyze(flows))
	p.ConfusionAnalysis(confusion.NewAnalyzer().Analyze(flows))

	if stats.FailedTests > 0 {
		fmt.Println("\n❌ Some tests failed")
		return fmt.Errorf("%d test(s) failed", stats.FailedTests)
	}
	fmt.Println("\n✨ All tests passed!")
	return nil
}

func filterScenario(scenarios []config.Scenario, name string) ([]config.Scenario, error) {
	for _, s := range scenarios {
		if s.Name == name {
			return []config.Scenario{s}, nil
		}
	}
	return nil, fmt.Errorf("scenario '%s' not found", name)
}
