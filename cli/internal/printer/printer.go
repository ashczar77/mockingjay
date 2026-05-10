package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/ashczar77/mockingjay/internal/confusion"
	"github.com/ashczar77/mockingjay/internal/dialogue"
	"github.com/ashczar77/mockingjay/internal/dropoff"
	"github.com/ashczar77/mockingjay/internal/flow"
	"github.com/ashczar77/mockingjay/internal/quality"
	"github.com/ashczar77/mockingjay/internal/test"
)

// Printer handles all CLI output formatting
type Printer struct {
	out io.Writer
}

// New creates a Printer writing to stdout
func New() *Printer {
	return &Printer{out: os.Stdout}
}

func (p *Printer) printf(format string, args ...any) {
	fmt.Fprintf(p.out, format, args...)
}

// Results prints individual test results
func (p *Printer) Results(results []test.Result) {
	for i, r := range results {
		p.printf("  [%d/%d] %s", i+1, len(results), r.Scenario)
		if r.Passed {
			p.printf(" ✓ PASS (latency: %dms)\n", r.Metrics.Latency.Milliseconds())
		} else {
			p.printf(" ✗ FAIL (%s)\n", r.Error)
		}
	}
}

// Stats prints aggregate statistics
func (p *Printer) Stats(stats test.Stats) {
	p.printf("\n📊 Results:\n")
	p.printf("  Tests run: %d\n", stats.TotalTests)
	p.printf("  Passed: %d\n", stats.PassedTests)
	p.printf("  Failed: %d\n", stats.FailedTests)
	p.printf("  Pass rate: %.1f%%\n", stats.PassRate)
	p.printf("\n⚡ Performance:\n")
	p.printf("  Avg latency: %dms\n", stats.AvgLatency)
	p.printf("  P95 latency: %dms\n", stats.P95Latency)
	p.printf("  P99 latency: %dms\n", stats.P99Latency)
	p.printf("  Task completion: %.1f%%\n", stats.TaskCompletion)
}

// ConversationInsights prints flow analysis
func (p *Printer) ConversationInsights(insights flow.FlowInsights) {
	p.printf("\n💬 Conversation Intelligence:\n")
	p.printf("  Success rate: %.1f%%\n", insights.SuccessRate)
	p.printf("  Intent accuracy: %.1f%% (%d/%d correct)\n", insights.IntentAccuracy, insights.CorrectIntents, insights.TotalIntentChecks)
	p.printf("  Avg steps completed: %.1f\n", insights.AvgStepsCompleted)
	p.printf("  Avg conversation duration: %dms\n", insights.AvgDuration.Milliseconds())
	if len(insights.CommonDropOffPoints) > 0 {
		p.printf("  Common drop-off points:\n")
		for step, count := range insights.CommonDropOffPoints {
			p.printf("    Step %d: %d failures\n", step, count)
		}
	}
}

// DialogueMetrics prints multi-turn dialogue analysis
func (p *Printer) DialogueMetrics(metrics dialogue.DialogueMetrics, contextLoss []dialogue.ContextLossPoint) {
	p.printf("\n🔄 Multi-turn Dialogue:\n")
	p.printf("  Multi-turn conversations: %d/%d\n", metrics.MultiTurnCount, metrics.TotalConversations)
	p.printf("  Avg turns per conversation: %.1f\n", metrics.AvgTurnsPerConv)
	p.printf("  Max turns: %d\n", metrics.MaxTurns)
	p.printf("  Context retention: %.1f%%\n", metrics.ContextRetention)
	if len(contextLoss) > 0 {
		p.printf("  Context loss detected:\n")
		for _, loss := range contextLoss {
			p.printf("    %s (step %d): expected '%s', got '%s'\n", loss.ScenarioName, loss.StepNumber, loss.Expected, loss.Actual)
		}
	}
}

// QualityMetrics prints response quality analysis
func (p *Printer) QualityMetrics(metrics quality.QualityMetrics) {
	p.printf("\n✨ Response Quality:\n")
	p.printf("  Avg response length: %.0f chars\n", metrics.AvgResponseLength)
	p.printf("  Completeness: %.1f%%\n", metrics.CompletenessScore)
	p.printf("  Positive sentiment: %.1f%%\n", metrics.SentimentScore)
	p.printf("  Confidence: %.1f%%\n", metrics.ConfidenceScore)
	if metrics.EmptyResponses > 0 {
		p.printf("  ⚠️  Empty responses: %d\n", metrics.EmptyResponses)
	}
}

// FailurePoints prints failure point detection results
func (p *Printer) FailurePoints(analysis dropoff.FailureAnalysis) {
	if len(analysis.FailurePoints) == 0 {
		return
	}
	p.printf("\n🚨 Failure Points:\n")
	p.printf("  Overall failure rate: %.1f%%\n", analysis.OverallFailureRate)
	p.printf("  Failure points found: %d\n", len(analysis.FailurePoints))
	if len(analysis.CriticalPoints) > 0 {
		p.printf("  ⚠️  Critical issues: %d\n", len(analysis.CriticalPoints))
		p.printf("\n  Critical failure points:\n")
		for _, point := range analysis.CriticalPoints {
			p.printf("    Step %d: \"%s\" - %.1f%% failure rate (%s)\n", point.StepNumber, point.StepInput, point.FailureRate, point.Severity)
		}
	}
}

// ConfusionAnalysis prints confusion pattern results
func (p *Printer) ConfusionAnalysis(analysis confusion.ConfusionAnalysis) {
	if len(analysis.Patterns) == 0 {
		return
	}
	p.printf("\n🤔 Confusion Patterns:\n")
	p.printf("  Confusion rate: %.1f%% (%d/%d intents)\n", analysis.ConfusionRate, analysis.ConfusedIntents, analysis.TotalIntentChecks)
	p.printf("  Patterns detected: %d\n", len(analysis.Patterns))
	if len(analysis.MostConfusedInputs) > 0 {
		p.printf("\n  Most confused inputs:\n")
		for _, pattern := range analysis.MostConfusedInputs {
			p.printf("    \"%s\"\n", pattern.UserInput)
			p.printf("      Expected: %s → Got: %s (%.1f%% confusion)\n", pattern.ExpectedIntent, pattern.ActualIntent, pattern.ConfusionRate)
		}
	}
}
