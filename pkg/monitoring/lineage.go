// Package monitoring provides data lineage tracking functionality.
package monitoring

import (
	"fmt"
	"time"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
)

// DataLineageTrackerImpl implements the DataLineageTracker interface.
type DataLineageTrackerImpl struct {
	analyzer     interfaces.DatabaseAnalyzer
	dependencies map[string][]string // table -> dependent tables
}

// NewDataLineageTracker creates a new data lineage tracker.
func NewDataLineageTracker(analyzer interfaces.DatabaseAnalyzer) interfaces.DataLineageTracker {
	return &DataLineageTrackerImpl{
		analyzer:     analyzer,
		dependencies: make(map[string][]string),
	}
}

// TrackLineage builds the data lineage for all tables.
func (dlt *DataLineageTrackerImpl) TrackLineage() (map[string]types.DataLineage, error) {
	lineage := make(map[string]types.DataLineage)

	// Get all tables
	tables, err := dlt.analyzer.AnalyzeTables()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze tables: %w", err)
	}

	// Build lineage based on foreign key relationships
	for _, table := range tables {
		tableLineage := types.DataLineage{
			TableName:   table.Name,
			LastUpdated: time.Now(),
		}

		// Find upstream dependencies (tables this table references)
		for _, relationship := range table.Relationships {
			if relationship.Type == "foreign_key" {
				dep := types.LineageDependency{
					TableName:  relationship.TargetTable,
					ColumnName: relationship.TargetColumn,
					Type:       "foreign_key",
				}
				tableLineage.UpstreamDeps = append(tableLineage.UpstreamDeps, dep)
			}
		}

		lineage[table.Name] = tableLineage
	}

	// Build downstream dependencies
	for tableName, tableLineage := range lineage {
		for _, upstream := range tableLineage.UpstreamDeps {
			if upstreamLineage, exists := lineage[upstream.TableName]; exists {
				downstreamDep := types.LineageDependency{
					TableName: tableName,
					Type:      upstream.Type,
				}
				upstreamLineage.DownstreamDeps = append(upstreamLineage.DownstreamDeps, downstreamDep)
				lineage[upstream.TableName] = upstreamLineage
			}
		}
	}

	return lineage, nil
}

// GetLineageForTable returns lineage information for a specific table.
func (dlt *DataLineageTrackerImpl) GetLineageForTable(tableName string) (*types.DataLineage, error) {
	lineage, err := dlt.TrackLineage()
	if err != nil {
		return nil, fmt.Errorf("failed to track lineage: %w", err)
	}

	if tableLineage, exists := lineage[tableName]; exists {
		return &tableLineage, nil
	}

	return nil, fmt.Errorf("table %s not found", tableName)
}
