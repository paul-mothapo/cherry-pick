package services

import (
	"fmt"
	"math"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (cs *CalculatorService) CalculateHealthScore(tables []core.TableInfo) float64 {
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

func (cs *CalculatorService) CalculateComplexityScore(tables []core.TableInfo) float64 {
	complexity := float64(len(tables)) * 0.1

	for _, table := range tables {
		complexity += float64(len(table.Columns)) * 0.05
		complexity += float64(len(table.Indexes)) * 0.1
		complexity += float64(len(table.Constraints)) * 0.05
		complexity += float64(len(table.Relationships)) * 0.05
	}

	return complexity
}

func (cs *CalculatorService) CalculateDataQuality(column core.ColumnInfo) float64 {
	if column.DataProfile.Quality > 0 {
		return column.DataProfile.Quality
	}

	totalRows := column.NullCount + column.UniqueValues
	if totalRows == 0 {
		return 1.0
	}

	nullRatio := float64(column.NullCount) / float64(totalRows)
	uniquenessRatio := float64(column.UniqueValues) / float64(totalRows)

	quality := 1.0 - nullRatio
	if uniquenessRatio > 0.9 {
		quality *= 0.8
	}

	return math.Max(0.0, math.Min(1.0, quality))
}

func (cs *CalculatorService) CalculatePerformanceScore(metrics *core.PerformanceMetrics) float64 {
	if metrics == nil {
		return 0.5
	}

	score := 1.0

	if metrics.Connections.Current > 0 && metrics.Connections.Available > 0 {
		connectionRatio := float64(metrics.Connections.Current) / float64(metrics.Connections.Current+metrics.Connections.Available)
		if connectionRatio > 0.8 {
			score -= 0.3
		}
	}

	totalOps := metrics.Operations.Insert + metrics.Operations.Query + metrics.Operations.Update + metrics.Operations.Delete
	if totalOps > 0 {
		errorRate := float64(metrics.Operations.Delete) / float64(totalOps)
		if errorRate > 0.1 {
			score -= 0.2
		}
	}

	return math.Max(0.0, math.Min(1.0, score))
}
