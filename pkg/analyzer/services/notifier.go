package services

import (
	"context"
	"log"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type NotifierService struct{}

func NewNotifierService() *NotifierService {
	return &NotifierService{}
}

func (ns *NotifierService) NotifyAnalysisComplete(ctx context.Context, result *core.AnalysisResult) error {
	log.Printf("Analysis completed for database %s (%s) - Health Score: %.2f, Complexity: %.2f",
		result.DatabaseName,
		result.DatabaseType,
		result.Summary.HealthScore,
		result.Summary.ComplexityScore,
	)

	for _, insight := range result.Insights {
		if insight.Severity == "high" {
			ns.NotifyInsight(ctx, insight)
		}
	}

	return nil
}

func (ns *NotifierService) NotifyAnalysisError(ctx context.Context, err error) error {
	log.Printf("Analysis failed: %v", err)
	return nil
}

func (ns *NotifierService) NotifyInsight(ctx context.Context, insight core.DatabaseInsight) error {
	log.Printf("Database Insight: [%s] %s - %s", insight.Severity, insight.Title, insight.Description)
	return nil
}
