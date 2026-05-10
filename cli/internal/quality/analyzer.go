package quality

import (
	"github.com/ashczar77/mockingjay/internal/classifier"
	"github.com/ashczar77/mockingjay/internal/flow"
)

// QualityMetrics contains response quality measurements
type QualityMetrics struct {
	TotalResponses    int
	AvgResponseLength float64
	CompletenessScore float64
	SentimentScore    float64
	ConfidenceScore   float64
	EmptyResponses    int
}

// QualityAnalyzer analyzes response quality
type QualityAnalyzer struct {
	classifier *classifier.Classifier
}

// NewQualityAnalyzer creates a new quality analyzer
func NewQualityAnalyzer() *QualityAnalyzer {
	return &QualityAnalyzer{classifier: classifier.New()}
}

// Analyze evaluates response quality from conversation flows
func (q *QualityAnalyzer) Analyze(flows []flow.ConversationFlow) QualityMetrics {
	metrics := QualityMetrics{}
	if len(flows) == 0 {
		return metrics
	}

	var totalLength int
	var totalCompleteness, totalSentiment, totalConfidence float64
	var scoredResponses int

	for _, f := range flows {
		for _, step := range f.Steps {
			if step.ActualResponse == "" {
				metrics.EmptyResponses++
				continue
			}
			metrics.TotalResponses++
			totalLength += len(step.ActualResponse)

			if q.classifier != nil {
				scores, err := q.classifier.EvaluateQuality(step.ActualResponse)
				if err == nil {
					totalCompleteness += scores["completeness"]
					totalSentiment += scores["sentiment"]
					totalConfidence += scores["confidence"]
					scoredResponses++
				}
			}
		}
	}

	if metrics.TotalResponses > 0 {
		metrics.AvgResponseLength = float64(totalLength) / float64(metrics.TotalResponses)
	}
	if scoredResponses > 0 {
		metrics.CompletenessScore = totalCompleteness / float64(scoredResponses) * 100
		metrics.SentimentScore = totalSentiment / float64(scoredResponses) * 100
		metrics.ConfidenceScore = totalConfidence / float64(scoredResponses) * 100
	}

	return metrics
}
