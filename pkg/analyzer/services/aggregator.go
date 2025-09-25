package services

import (
	"fmt"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type AggregatorService struct{}

func NewAggregatorService() *AggregatorService {
	return &AggregatorService{}
}

func (as *AggregatorService) AggregateTableStats(tables []core.TableInfo) core.DatabaseSummary {
	var totalRows int64
	var totalColumns int

	for _, table := range tables {
		totalRows += table.RowCount
		totalColumns += len(table.Columns)
	}

	healthScore := as.calculateHealthScore(tables)
	complexityScore := as.calculateComplexityScore(tables)

	return core.DatabaseSummary{
		TotalTables:     len(tables),
		TotalColumns:    totalColumns,
		TotalRows:       totalRows,
		TotalSize:       as.calculateTotalSize(tables),
		HealthScore:     healthScore,
		ComplexityScore: complexityScore,
	}
}

func (as *AggregatorService) AggregateInsights(insights []core.DatabaseInsight) []core.DatabaseInsight {
	insightMap := make(map[string]core.DatabaseInsight)

	for _, insight := range insights {
		key := fmt.Sprintf("%s_%s", insight.Type, insight.Title)
		if existing, exists := insightMap[key]; exists {
			if insight.Severity == "high" && existing.Severity != "high" {
				insightMap[key] = insight
			}
		} else {
			insightMap[key] = insight
		}
	}

	var aggregated []core.DatabaseInsight
	for _, insight := range insightMap {
		aggregated = append(aggregated, insight)
	}

	return aggregated
}

func (as *AggregatorService) AggregateRecommendations(recommendations []string) []string {
	recommendationMap := make(map[string]bool)
	var unique []string

	for _, rec := range recommendations {
		if !recommendationMap[rec] {
			recommendationMap[rec] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

func (as *AggregatorService) calculateHealthScore(tables []core.TableInfo) float64 {
	if len(tables) == 0 {
		return 0.0
	}

	var totalScore float64
	for _, table := range tables {
		tableScore := 1.0

		if len(table.Indexes) <= 1 && table.RowCount > 1000 {
			tableScore -= 0.2
		}

		if table.RowCount > 10000000 {
			tableScore -= 0.3
		}

		if tableScore < 0 {
			tableScore = 0
		}

		totalScore += tableScore
	}

	return totalScore / float64(len(tables))
}

func (as *AggregatorService) calculateComplexityScore(tables []core.TableInfo) float64 {
	complexity := float64(len(tables)) * 0.1

	for _, table := range tables {
		complexity += float64(len(table.Columns)) * 0.05
		complexity += float64(len(table.Indexes)) * 0.1
		complexity += float64(len(table.Constraints)) * 0.05
		complexity += float64(len(table.Relationships)) * 0.05
	}

	return complexity
}

func (as *AggregatorService) calculateTotalSize(tables []core.TableInfo) string {
	var totalSize int64
	for _, table := range tables {
		if size, err := parseSize(table.Size); err == nil {
			totalSize += size
		}
	}

	if totalSize < 1024 {
		return fmt.Sprintf("%d bytes", totalSize)
	} else if totalSize < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(totalSize)/1024)
	} else if totalSize < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(totalSize)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(totalSize)/(1024*1024*1024))
	}
}

func parseSize(sizeStr string) (int64, error) {
	if sizeStr == "Unknown" {
		return 0, fmt.Errorf("unknown size")
	}

	var size int64
	var unit string
	_, err := fmt.Sscanf(sizeStr, "%d%s", &size, &unit)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "bytes":
		return size, nil
	case "KB":
		return size * 1024, nil
	case "MB":
		return size * 1024 * 1024, nil
	case "GB":
		return size * 1024 * 1024 * 1024, nil
	default:
		return size, nil
	}
}
