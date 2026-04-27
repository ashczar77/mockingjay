package cmd

import (
	"fmt"
	"os"

	"github.com/ashczar77/mockingjay/internal/ab"
	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/test"
	"github.com/spf13/cobra"
)

var abCmd = &cobra.Command{
	Use:   "ab",
	Short: "Run A/B test between two agent variants",
	Long:  `Compare two versions of your voice AI agent side-by-side.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runABTest(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	abCmd.Flags().StringVarP(&configFile, "config", "c", "mockingjay.yaml", "config file")
}

func runABTest() error {
	fmt.Println("🐦 MockingJay - A/B Test")
	fmt.Println()

	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.ABTest == nil {
		return fmt.Errorf("no ab_test configuration found in %s", configFile)
	}

	varA := cfg.ABTest.VariantA
	varB := cfg.ABTest.VariantB

	fmt.Printf("📋 Variants: %s vs %s\n", varA.Name, varB.Name)
	fmt.Printf("🎯 Running %d scenario(s) per variant...\n\n", len(cfg.Scenarios))

	execA := test.New(&config.Config{Version: cfg.Version, Agent: config.Agent{Endpoint: varA.Endpoint}, Scenarios: cfg.Scenarios})
	fmt.Printf("  Running %s...\n", varA.Name)
	resultsA := execA.RunAll(cfg.Scenarios)

	execB := test.New(&config.Config{Version: cfg.Version, Agent: config.Agent{Endpoint: varB.Endpoint}, Scenarios: cfg.Scenarios})
	fmt.Printf("  Running %s...\n", varB.Name)
	resultsB := execB.RunAll(cfg.Scenarios)

	fmt.Println()

	comparison := ab.Compare(
		&ab.Variant{Name: varA.Name, Endpoint: varA.Endpoint, Results: resultsA},
		&ab.Variant{Name: varB.Name, Endpoint: varB.Endpoint, Results: resultsB},
	)
	ab.PrintComparison(comparison)

	return nil
}
