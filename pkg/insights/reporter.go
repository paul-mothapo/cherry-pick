package insights

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
)

type ReportGeneratorImpl struct{}

func NewReportGenerator() interfaces.ReportGenerator {
	return &ReportGeneratorImpl{}
}

func (rg *ReportGeneratorImpl) GenerateSummary(tables []types.TableInfo) types.DatabaseSummary {
	var totalRows int64
	var totalColumns int

	for _, table := range tables {
		totalRows += table.RowCount
		totalColumns += len(table.Columns)
	}

	return types.DatabaseSummary{
		TotalTables:     len(tables),
		TotalColumns:    totalColumns,
		TotalRows:       totalRows,
		TotalSize:       "Unknown",
		HealthScore:     rg.CalculateHealthScore(tables),
		ComplexityScore: rg.CalculateComplexityScore(tables),
	}
}

func (rg *ReportGeneratorImpl) GenerateRecommendations(tables []types.TableInfo, insights []types.DatabaseInsight) []string {
	var recommendations []string

	if len(tables) > 50 {
		recommendations = append(recommendations,
			"Consider database normalization to reduce the number of tables")
	}

	highPriorityCount := 0
	for _, insight := range insights {
		if insight.Severity == "high" {
			highPriorityCount++
			recommendations = append(recommendations,
				fmt.Sprintf("Priority: %s", insight.Suggestion))
		}
	}

	if highPriorityCount == 0 {
		recommendations = append(recommendations,
			"Database appears to be in good condition with no critical issues")
	}

	return recommendations
}

func (rg *ReportGeneratorImpl) CalculateHealthScore(tables []types.TableInfo) float64 {
	if len(tables) == 0 {
		return 0.0
	}

	var totalScore float64
	for _, table := range tables {
		tableScore := 1.0

		if len(table.Indexes) == 0 && table.RowCount > 1000 {
			tableScore -= 0.2
		}

		for _, column := range table.Columns {
			if column.DataProfile.Quality < 0.7 {
				tableScore -= 0.1
			}
		}

		if tableScore < 0 {
			tableScore = 0
		}

		totalScore += tableScore
	}

	return totalScore / float64(len(tables))
}

func (rg *ReportGeneratorImpl) CalculateComplexityScore(tables []types.TableInfo) float64 {
	complexity := float64(len(tables)) * 0.1

	for _, table := range tables {
		complexity += float64(len(table.Columns)) * 0.05
		complexity += float64(len(table.Relationships)) * 0.1
	}

	return complexity
}

func (rg *ReportGeneratorImpl) ExportReport(report *types.DatabaseReport, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(report, "", "  ")
	case "summary":
		return rg.generateTextSummary(report), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

func (rg *ReportGeneratorImpl) generateTextSummary(report *types.DatabaseReport) []byte {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Database Analysis Report - %s\n", report.DatabaseName))
	summary.WriteString(fmt.Sprintf("Analysis Date: %s\n\n", report.AnalysisTime.Format("2006-01-02 15:04:05")))

	summary.WriteString("SUMMARY\n")
	summary.WriteString("=======\n")
	summary.WriteString(fmt.Sprintf("Tables: %d\n", report.Summary.TotalTables))
	summary.WriteString(fmt.Sprintf("Total Columns: %d\n", report.Summary.TotalColumns))
	summary.WriteString(fmt.Sprintf("Total Rows: %d\n", report.Summary.TotalRows))
	summary.WriteString(fmt.Sprintf("Health Score: %.2f/1.0\n", report.Summary.HealthScore))
	summary.WriteString(fmt.Sprintf("Complexity Score: %.2f\n\n", report.Summary.ComplexityScore))

	summary.WriteString("KEY INSIGHTS\n")
	summary.WriteString("============\n")
	for _, insight := range report.Insights {
		summary.WriteString(fmt.Sprintf("• [%s] %s: %s\n",
			strings.ToUpper(insight.Severity), insight.Title, insight.Description))
	}

	summary.WriteString("\nRECOMMENDATIONS\n")
	summary.WriteString("===============\n")
	for _, rec := range report.Recommendations {
		summary.WriteString(fmt.Sprintf("• %s\n", rec))
	}

	return []byte(summary.String())
}
