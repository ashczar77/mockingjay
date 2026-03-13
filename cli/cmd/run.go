package cmd

import (
	"fmt"
	"os"

	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/dialogue"
	"github.com/ashczar77/mockingjay/internal/dropoff"
	"github.com/ashczar77/mockingjay/internal/flow"
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

	// Print individual results
	for i, r := range results {
		fmt.Printf("  [%d/%d] %s", i+1, len(results), r.Scenario)
		if r.Passed {
			fmt.Printf(" ✓ PASS (latency: %dms)\n", r.Metrics.Latency.Milliseconds())
		} else {
			fmt.Printf(" ✗ FAIL (%s)\n", r.Error)
		}

		// Send to backend
		if reporterClient != nil {
			if err := reporterClient.Report(r); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to report result: %v\n", err)
			}
		}
	}

	// Calculate and display statistics
	stats := test.CalculateStats(results)

	fmt.Println()
	fmt.Println("📊 Results:")
	fmt.Printf("  Tests run: %d\n", stats.TotalTests)
	fmt.Printf("  Passed: %d\n", stats.PassedTests)
	fmt.Printf("  Failed: %d\n", stats.FailedTests)
	fmt.Printf("  Pass rate: %.1f%%\n", stats.PassRate)
	fmt.Println()
	fmt.Println("⚡ Performance:")
	fmt.Printf("  Avg latency: %dms\n", stats.AvgLatency)
	fmt.Printf("  P95 latency: %dms\n", stats.P95Latency)
	fmt.Printf("  P99 latency: %dms\n", stats.P99Latency)
	fmt.Printf("  Task completion: %.1f%%\n", stats.TaskCompletion)
	fmt.Println()

	// Analyze conversation flows
	analyzer := flow.NewAnalyzer()
	flows := analyzer.AnalyzeMultiple(results, scenarios)
	insights := analyzer.GenerateInsights(flows)

	// Analyze multi-turn dialogues
	dialogueAnalyzer := dialogue.NewDialogueAnalyzer()
	dialogueMetrics := dialogueAnalyzer.Analyze(flows)
	contextLoss := dialogueAnalyzer.DetectContextLoss(flows)

	// Analyze response quality
	qualityAnalyzer := quality.NewQualityAnalyzer()
	qualityMetrics := qualityAnalyzer.Analyze(flows)

	// Detect drop-off points
	dropoffDetector := dropoff.NewDetector()
	dropoffAnalysis := dropoffDetector.Analyze(flows)

	fmt.Println("💬 Conversation Intelligence:")
	fmt.Printf("  Success rate: %.1f%%\n", insights.SuccessRate)
	fmt.Printf("  Intent accuracy: %.1f%% (%d/%d correct)\n", insights.IntentAccuracy, insights.CorrectIntents, insights.TotalIntentChecks)
	fmt.Printf("  Avg steps completed: %.1f\n", insights.AvgStepsCompleted)
	fmt.Printf("  Avg conversation duration: %dms\n", insights.AvgDuration.Milliseconds())
	
	if len(insights.CommonDropOffPoints) > 0 {
		fmt.Println("  Common drop-off points:")
		for step, count := range insights.CommonDropOffPoints {
			fmt.Printf("    Step %d: %d failures\n", step, count)
		}
	}
	fmt.Println()

	fmt.Println("🔄 Multi-turn Dialogue:")
	fmt.Printf("  Multi-turn conversations: %d/%d\n", dialogueMetrics.MultiTurnCount, dialogueMetrics.TotalConversations)
	fmt.Printf("  Avg turns per conversation: %.1f\n", dialogueMetrics.AvgTurnsPerConv)
	fmt.Printf("  Max turns: %d\n", dialogueMetrics.MaxTurns)
	fmt.Printf("  Context retention: %.1f%%\n", dialogueMetrics.ContextRetention)
	fmt.Printf("  Coherence score: %.1f%%\n", dialogueMetrics.CoherenceScore)
	
	if len(contextLoss) > 0 {
		fmt.Println("  Context loss detected:")
		for _, loss := range contextLoss {
			fmt.Printf("    %s (step %d): expected '%s', got '%s'\n", 
				loss.ScenarioName, loss.StepNumber, loss.Expected, loss.Actual)
		}
	}
	fmt.Println()

	fmt.Println("✨ Response Quality:")
	fmt.Printf("  Avg response length: %.0f chars\n", qualityMetrics.AvgResponseLength)
	fmt.Printf("  Completeness: %.1f%%\n", qualityMetrics.CompletenessScore)
	fmt.Printf("  Positive sentiment: %.1f%%\n", qualityMetrics.SentimentScore)
	fmt.Printf("  Confidence: %.1f%%\n", qualityMetrics.ConfidenceScore)
	
	if qualityMetrics.VagueResponses > 0 {
		fmt.Printf("  ⚠️  Vague responses: %d\n", qualityMetrics.VagueResponses)
	}
	if qualityMetrics.EmptyResponses > 0 {
		fmt.Printf("  ⚠️  Empty responses: %d\n", qualityMetrics.EmptyResponses)
	}
	fmt.Println()

	if len(dropoffAnalysis.DropOffPoints) > 0 {
		fmt.Println("🚨 Drop-off Detection:")
		fmt.Printf("  Overall drop-off rate: %.1f%%\n", dropoffAnalysis.OverallDropOffRate)
		fmt.Printf("  Drop-off points found: %d\n", len(dropoffAnalysis.DropOffPoints))
		
		if len(dropoffAnalysis.CriticalPoints) > 0 {
			fmt.Printf("  ⚠️  Critical issues: %d\n", len(dropoffAnalysis.CriticalPoints))
			fmt.Println("\n  Critical drop-off points:")
			for _, point := range dropoffAnalysis.CriticalPoints {
				fmt.Printf("    Step %d: \"%s\" - %.1f%% drop-off (%s)\n", 
					point.StepNumber, point.StepInput, point.DropOffRate, point.Severity)
			}
		}
		fmt.Println()
	}

	if stats.FailedTests > 0 {
		fmt.Println("❌ Some tests failed")
		return fmt.Errorf("%d test(s) failed", stats.FailedTests)
	}

	fmt.Println("✨ All tests passed!")
	return nil
}
