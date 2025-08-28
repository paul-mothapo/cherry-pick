// Package interfaces defines monitoring and alerting interfaces.
package interfaces

import (
	"time"

	"github.com/cherry-pick/pkg/types"
)

// AlertManager defines the interface for handling database monitoring alerts.
type AlertManager interface {
	// CheckAlerts evaluates alerts against the current database report.
	CheckAlerts(report *types.DatabaseReport) []types.MonitoringAlert

	// AddAlert adds a new alert condition.
	AddAlert(alert types.MonitoringAlert) error

	// RemoveAlert removes an alert condition.
	RemoveAlert(alertID string) error

	// GetAlerts returns all configured alerts.
	GetAlerts() []types.MonitoringAlert
}

// ComparisonEngine defines the interface for comparing database reports.
type ComparisonEngine interface {
	// CompareReports compares two database reports to identify changes.
	CompareReports(oldReport, newReport *types.DatabaseReport) *types.ComparisonReport
}

// DataLineageTracker defines the interface for tracking data lineage and dependencies.
type DataLineageTracker interface {
	// TrackLineage builds the data lineage for all tables.
	TrackLineage() (map[string]types.DataLineage, error)

	// GetLineageForTable returns lineage information for a specific table.
	GetLineageForTable(tableName string) (*types.DataLineage, error)
}

// Scheduler defines the interface for automated periodic analysis.
type Scheduler interface {
	// ScheduleAnalysis allows for automated periodic analysis.
	ScheduleAnalysis(interval time.Duration, callback func(*types.DatabaseReport)) error

	// Stop stops the scheduled analysis.
	Stop() error
}
