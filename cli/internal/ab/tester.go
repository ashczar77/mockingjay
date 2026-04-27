package ab

import (
	"fmt"
	"sort"

	"github.com/ashczar77/mockingjay/internal/test"
)

// Variant represents one version of an agent being tested
type Variant struct {
	Name     string
	Endpoint string
	Results  []test.Result
}

// Comparison holds the A/B test comparison results
type Comparison struct {
	VariantA    *Variant
	VariantB    *Variant
	Winner      string // name of winning variant, or "tie"
	Differences []Difference
}

// Difference describes a metric difference between variants
type Difference struct {
	Metric    string
	ValueA    float64
	ValueB    float64
	Delta     float64 // B - A
	Improved  bool    // true if B is better
}

// Compare runs a statistical comparison between two variants
func Compare(a, b *Variant) Comparison {
	statsA := test.CalculateStats(a.Results)
	statsB := test.CalculateStats(b.Results)

	diffs := []Difference{
		metric("Pass Rate (%)", statsA.PassRate, statsB.PassRate, true),
		metric("Avg Latency (ms)", float64(statsA.AvgLatency), float64(statsB.AvgLatency), false),
		metric("P95 Latency (ms)", float64(statsA.P95Latency), float64(statsB.P95Latency), false),
		metric("Task Completion (%)", statsA.TaskCompletion, statsB.TaskCompletion, true),
	}

	// Determine winner by counting improvements
	aWins, bWins := 0, 0
	for _, d := range diffs {
		if d.Delta > 0 && d.Improved {
			bWins++
		} else if d.Delta < 0 && d.Improved {
			aWins++
		} else if d.Delta < 0 && !d.Improved {
			bWins++
		} else if d.Delta > 0 && !d.Improved {
			aWins++
		}
	}

	winner := "tie"
	if aWins > bWins {
		winner = a.Name
	} else if bWins > aWins {
		winner = b.Name
	}

	return Comparison{
		VariantA:    a,
		VariantB:    b,
		Winner:      winner,
		Differences: diffs,
	}
}

func metric(name string, a, b float64, higherIsBetter bool) Difference {
	delta := b - a
	improved := (delta > 0) == higherIsBetter
	return Difference{
		Metric:   name,
		ValueA:   a,
		ValueB:   b,
		Delta:    delta,
		Improved: improved,
	}
}

// PrintComparison prints a formatted A/B test comparison
func PrintComparison(c Comparison) {
	fmt.Printf("🔬 A/B Test: %s vs %s\n", c.VariantA.Name, c.VariantB.Name)
	fmt.Println()

	// Header
	fmt.Printf("  %-25s %12s %12s %12s\n", "Metric", c.VariantA.Name, c.VariantB.Name, "Delta")
	fmt.Printf("  %-25s %12s %12s %12s\n", "------", "--------", "--------", "-----")

	// Sort diffs for consistent output
	diffs := make([]Difference, len(c.Differences))
	copy(diffs, c.Differences)
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Metric < diffs[j].Metric
	})

	for _, d := range diffs {
		indicator := "  "
		if d.Delta != 0 {
			if d.Improved {
				indicator = "✓ "
			} else {
				indicator = "✗ "
			}
		}
		fmt.Printf("  %-25s %12.1f %12.1f %s%+.1f\n",
			d.Metric, d.ValueA, d.ValueB, indicator, d.Delta)
	}

	fmt.Println()
	if c.Winner == "tie" {
		fmt.Println("  🤝 Result: Tie - no significant difference")
	} else {
		fmt.Printf("  🏆 Winner: %s\n", c.Winner)
	}
}
