// Package monitoring provides database comparison functionality.
package monitoring

import (
	"fmt"

	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
	"github.com/cherry-pick/pkg/utils"
)

// ComparisonEngineImpl implements the ComparisonEngine interface.
type ComparisonEngineImpl struct{}

// NewComparisonEngine creates a new comparison engine instance.
func NewComparisonEngine() interfaces.ComparisonEngine {
	return &ComparisonEngineImpl{}
}

// CompareReports compares two database reports to identify changes.
func (ce *ComparisonEngineImpl) CompareReports(oldReport, newReport *types.DatabaseReport) *types.ComparisonReport {
	comparison := &types.ComparisonReport{
		OldAnalysisTime: oldReport.AnalysisTime,
		NewAnalysisTime: newReport.AnalysisTime,
		Changes:         []types.DatabaseChange{},
	}

	// Compare table counts
	if oldReport.Summary.TotalTables != newReport.Summary.TotalTables {
		change := types.DatabaseChange{
			Type:     "schema",
			Category: "table_count",
			Description: fmt.Sprintf("Table count changed from %d to %d",
				oldReport.Summary.TotalTables, newReport.Summary.TotalTables),
			Impact:   utils.CalculateImpact("table_count", oldReport.Summary.TotalTables, newReport.Summary.TotalTables),
			OldValue: oldReport.Summary.TotalTables,
			NewValue: newReport.Summary.TotalTables,
		}
		comparison.Changes = append(comparison.Changes, change)
	}

	// Compare row counts for existing tables
	oldTables := make(map[string]types.TableInfo)
	for _, table := range oldReport.Tables {
		oldTables[table.Name] = table
	}

	for _, newTable := range newReport.Tables {
		if oldTable, exists := oldTables[newTable.Name]; exists {
			if oldTable.RowCount != newTable.RowCount {
				change := types.DatabaseChange{
					Type:     "data",
					Category: "row_count",
					Description: fmt.Sprintf("Table '%s' row count changed from %d to %d",
						newTable.Name, oldTable.RowCount, newTable.RowCount),
					Impact:        utils.CalculateRowCountImpact(oldTable.RowCount, newTable.RowCount),
					AffectedTable: newTable.Name,
					OldValue:      oldTable.RowCount,
					NewValue:      newTable.RowCount,
				}
				comparison.Changes = append(comparison.Changes, change)
			}
		} else {
			// New table detected
			change := types.DatabaseChange{
				Type:     "schema",
				Category: "new_table",
				Description: fmt.Sprintf("New table '%s' added with %d rows",
					newTable.Name, newTable.RowCount),
				Impact:        "medium",
				AffectedTable: newTable.Name,
				NewValue:      newTable.RowCount,
			}
			comparison.Changes = append(comparison.Changes, change)
		}
	}

	// Check for removed tables
	newTables := make(map[string]bool)
	for _, table := range newReport.Tables {
		newTables[table.Name] = true
	}

	for _, oldTable := range oldReport.Tables {
		if !newTables[oldTable.Name] {
			change := types.DatabaseChange{
				Type:     "schema",
				Category: "removed_table",
				Description: fmt.Sprintf("Table '%s' was removed (had %d rows)",
					oldTable.Name, oldTable.RowCount),
				Impact:        "high",
				AffectedTable: oldTable.Name,
				OldValue:      oldTable.RowCount,
			}
			comparison.Changes = append(comparison.Changes, change)
		}
	}

	// Generate summary
	comparison.Summary = ce.generateChangeSummary(comparison.Changes)

	return comparison
}

// generateChangeSummary creates a summary of changes.
func (ce *ComparisonEngineImpl) generateChangeSummary(changes []types.DatabaseChange) types.ChangeSummary {
	summary := types.ChangeSummary{
		TotalChanges: len(changes),
	}

	for _, change := range changes {
		switch change.Type {
		case "schema":
			summary.SchemaChanges++
		case "data":
			summary.DataChanges++
		}

		switch change.Impact {
		case "high":
			summary.HighImpact++
		case "medium":
			summary.MediumImpact++
		case "low":
			summary.LowImpact++
		}
	}

	return summary
}
