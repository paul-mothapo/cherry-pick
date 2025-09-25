package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type ReporterService struct{}

func NewReporterService() *ReporterService {
	return &ReporterService{}
}

func (rs *ReporterService) GenerateReport(ctx context.Context, result *core.AnalysisResult) (map[string]string, error) {
	reports := make(map[string]string)

	summary, err := rs.GenerateSummary(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	reports["summary"] = summary

	insights, err := rs.GenerateInsights(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate insights: %w", err)
	}

	insightsReport := ""
	for _, insight := range insights {
		insightsReport += fmt.Sprintf("- [%s] %s: %s\n", insight.Severity, insight.Title, insight.Description)
	}
	reports["insights"] = insightsReport

	recommendations, err := rs.GenerateRecommendations(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	recommendationsReport := ""
	for _, rec := range recommendations {
		recommendationsReport += fmt.Sprintf("- %s\n", rec)
	}
	reports["recommendations"] = recommendationsReport

	return reports, nil
}

func (rs *ReporterService) GenerateSummary(ctx context.Context, result *core.AnalysisResult) (string, error) {
	summary := fmt.Sprintf(`
Database Analysis Report
========================

Database: %s (%s)
Analysis Time: %s
Total Tables: %d
Total Columns: %d
Total Rows: %d
Total Size: %s
Health Score: %.2f/1.0
Complexity Score: %.2f

Performance Metrics:
`, 
		result.DatabaseName,
		result.DatabaseType,
		result.AnalysisTime.Format(time.RFC3339),
		result.Summary.TotalTables,
		result.Summary.TotalColumns,
		result.Summary.TotalRows,
		result.Summary.TotalSize,
		result.Summary.HealthScore,
		result.Summary.ComplexityScore,
	)

	if result.Performance != nil {
		summary += fmt.Sprintf(`
- Current Connections: %d
- Available Connections: %d
- Total Operations: %d
- Insert Operations: %d
- Query Operations: %d
- Update Operations: %d
- Delete Operations: %d
`,
			result.Performance.Connections.Current,
			result.Performance.Connections.Available,
			result.Performance.Operations.Insert+result.Performance.Operations.Query+result.Performance.Operations.Update+result.Performance.Operations.Delete,
			result.Performance.Operations.Insert,
			result.Performance.Operations.Query,
			result.Performance.Operations.Update,
			result.Performance.Operations.Delete,
		)
	}

	return summary, nil
}

func (rs *ReporterService) GenerateInsights(ctx context.Context, result *core.AnalysisResult) ([]core.DatabaseInsight, error) {
	return result.Insights, nil
}

func (rs *ReporterService) GenerateRecommendations(ctx context.Context, result *core.AnalysisResult) ([]string, error) {
	return result.Recommendations, nil
}
