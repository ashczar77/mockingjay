package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ashczar77/mockingjay/internal/config"
	"github.com/ashczar77/mockingjay/internal/reporter"
	"github.com/ashczar77/mockingjay/internal/test"
	"github.com/spf13/cobra"
)

var (
	monitorInterval  int
	monitorThreshold float64
	monitorWebhook   string
	monitorAPIURL    string
	monitorConfig    string
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Continuously run tests and alert on failures",
	Long:  `Runs test scenarios on a schedule and alerts when pass rate drops below a threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMonitor(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	monitorCmd.Flags().IntVar(&monitorInterval, "interval", 60, "Interval between test runs in seconds")
	monitorCmd.Flags().Float64Var(&monitorThreshold, "threshold", 80.0, "Pass rate threshold for alerting (percent)")
	monitorCmd.Flags().StringVar(&monitorWebhook, "alert-webhook", "", "Webhook URL to POST alerts to (e.g. Slack)")
	monitorCmd.Flags().StringVar(&monitorAPIURL, "api-url", "", "Backend URL to report results")
	monitorCmd.Flags().StringVar(&monitorConfig, "config", "mockingjay.yaml", "Config file")
}

func runMonitor() error {
	cfg, err := config.Load(monitorConfig)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("🐦 MockingJay - Production Monitor")
	fmt.Printf("   Config:    %s\n", monitorConfig)
	fmt.Printf("   Interval:  %ds\n", monitorInterval)
	fmt.Printf("   Threshold: %.0f%% pass rate\n", monitorThreshold)
	if monitorWebhook != "" {
		fmt.Printf("   Alerts:    %s\n", monitorWebhook)
	}
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	var reporterClient *reporter.Client
	if monitorAPIURL != "" {
		reporterClient = reporter.NewClient(monitorAPIURL)
	}

	run := 0
	for {
		run++
		now := time.Now().Format("15:04:05")
		fmt.Printf("[%s] Run #%d — running %d scenario(s)...\n", now, run, len(cfg.Scenarios))

		executor := test.New(cfg)
		results := executor.RunAll(cfg.Scenarios)
		stats := test.CalculateStats(results)

		status := "✓ PASS"
		if stats.PassRate < monitorThreshold {
			status = "✗ FAIL"
		}

		fmt.Printf("[%s] Run #%d — %s  pass rate: %.1f%%  avg latency: %dms\n",
			now, run, status, stats.PassRate, stats.AvgLatency)

		// Report to backend
		if reporterClient != nil {
			for _, r := range results {
				reporterClient.Report(r)
			}
		}

		// Alert if below threshold
		if stats.PassRate < monitorThreshold {
			msg := fmt.Sprintf("🚨 MockingJay alert: pass rate %.1f%% is below threshold %.0f%% (run #%d at %s)",
				stats.PassRate, monitorThreshold, run, now)
			fmt.Println(msg)
			if monitorWebhook != "" {
				sendWebhookAlert(monitorWebhook, msg, stats.PassRate, monitorThreshold)
			}
		}

		time.Sleep(time.Duration(monitorInterval) * time.Second)
	}
}

func sendWebhookAlert(webhookURL, message string, passRate, threshold float64) {
	payload, _ := json.Marshal(map[string]interface{}{
		"text": message,
		"attachments": []map[string]interface{}{
			{
				"color": "danger",
				"fields": []map[string]string{
					{"title": "Pass Rate", "value": fmt.Sprintf("%.1f%%", passRate), "short": "true"},
					{"title": "Threshold", "value": fmt.Sprintf("%.0f%%", threshold), "short": "true"},
				},
			},
		},
	})
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to send alert: %v\n", err)
		return
	}
	defer resp.Body.Close()
}
