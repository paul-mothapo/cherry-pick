package insights

import (
	"fmt"
	"sort"

	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
)

type InsightGeneratorImpl struct{}

func NewInsightGenerator() interfaces.InsightGenerator {
	return &InsightGeneratorImpl{}
}

func (ig *InsightGeneratorImpl) GenerateInsights(tables []types.TableInfo) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	insights = append(insights, ig.AnalyzeLargeTables(tables)...)
	insights = append(insights, ig.AnalyzeMissingIndexes(tables)...)
	insights = append(insights, ig.AnalyzeDataQuality(tables)...)
	insights = append(insights, ig.AnalyzeRelationships(tables)...)
	insights = append(insights, ig.AnalyzeUnusedColumns(tables)...)

	return insights
}

func (ig *InsightGeneratorImpl) AnalyzeLargeTables(tables []types.TableInfo) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	sortedTables := make([]types.TableInfo, len(tables))
	copy(sortedTables, tables)
	sort.Slice(sortedTables, func(i, j int) bool {
		return sortedTables[i].RowCount > sortedTables[j].RowCount
	})

	for _, table := range sortedTables {
		if table.RowCount > 1000000 {
			insight := types.DatabaseInsight{
				Type:     "performance",
				Severity: "medium",
				Title:    "Large Table Detected",
				Description: fmt.Sprintf("Table '%s' has %d rows, which may impact performance",
					table.Name, table.RowCount),
				Suggestion:     "Consider partitioning, archiving old data, or optimizing queries",
				AffectedTables: []string{table.Name},
				MetricValue:    table.RowCount,
			}
			insights = append(insights, insight)
		}
	}

	return insights
}

func (ig *InsightGeneratorImpl) AnalyzeMissingIndexes(tables []types.TableInfo) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	for _, table := range tables {
		if table.RowCount > 10000 && len(table.Indexes) <= 1 {
			insight := types.DatabaseInsight{
				Type:     "performance",
				Severity: "high",
				Title:    "Potential Missing Indexes",
				Description: fmt.Sprintf("Table '%s' has %d rows but only %d indexes",
					table.Name, table.RowCount, len(table.Indexes)),
				Suggestion:     "Consider adding indexes on frequently queried columns",
				AffectedTables: []string{table.Name},
				MetricValue:    len(table.Indexes),
			}
			insights = append(insights, insight)
		}
	}

	return insights
}

func (ig *InsightGeneratorImpl) AnalyzeDataQuality(tables []types.TableInfo) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	for _, table := range tables {
		for _, column := range table.Columns {
			if column.DataProfile.Quality < 0.7 {
				insight := types.DatabaseInsight{
					Type:     "quality",
					Severity: "medium",
					Title:    "Poor Data Quality",
					Description: fmt.Sprintf("Column '%s.%s' has quality score of %.2f",
						table.Name, column.Name, column.DataProfile.Quality),
					Suggestion:     "Review data validation rules and consider data cleanup",
					AffectedTables: []string{table.Name},
					MetricValue:    column.DataProfile.Quality,
				}
				insights = append(insights, insight)
			}
		}
	}

	return insights
}

func (ig *InsightGeneratorImpl) AnalyzeRelationships(tables []types.TableInfo) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	orphanTables := 0
	for _, table := range tables {
		if len(table.Relationships) == 0 && table.RowCount > 100 {
			orphanTables++
		}
	}

	if orphanTables > len(tables)/4 {
		insight := types.DatabaseInsight{
			Type:     "design",
			Severity: "medium",
			Title:    "Many Isolated Tables",
			Description: fmt.Sprintf("%d tables have no relationships, which may indicate poor normalization",
				orphanTables),
			Suggestion:  "Review database design and consider establishing proper relationships",
			MetricValue: orphanTables,
		}
		insights = append(insights, insight)
	}

	return insights
}

func (ig *InsightGeneratorImpl) AnalyzeUnusedColumns(tables []types.TableInfo) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	for _, table := range tables {
		if table.RowCount == 0 {
			continue
		}

		for _, column := range table.Columns {
			nullRatio := float64(column.NullCount) / float64(table.RowCount)
			if nullRatio > 0.95 {
				insight := types.DatabaseInsight{
					Type:     "optimization",
					Severity: "low",
					Title:    "Potentially Unused Column",
					Description: fmt.Sprintf("Column '%s.%s' has %.1f%% null values",
						table.Name, column.Name, nullRatio*100),
					Suggestion:     "Consider removing this column if it's truly unused",
					AffectedTables: []string{table.Name},
					MetricValue:    nullRatio,
				}
				insights = append(insights, insight)
			}
		}
	}

	return insights
}
