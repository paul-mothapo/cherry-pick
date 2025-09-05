package interfaces

import (
	"time"

	"github.com/cherry-pick/pkg/types"
)

type AlertManager interface {
	CheckAlerts(report *types.DatabaseReport) []types.MonitoringAlert

	AddAlert(alert types.MonitoringAlert) error
	RemoveAlert(alertID string) error
	GetAlerts() []types.MonitoringAlert
}

type ComparisonEngine interface {
	CompareReports(oldReport, newReport *types.DatabaseReport) *types.ComparisonReport
}

type DataLineageTracker interface {
	TrackLineage() (map[string]types.DataLineage, error)

	GetLineageForTable(tableName string) (*types.DataLineage, error)
}

type Scheduler interface {
	ScheduleAnalysis(interval time.Duration, callback func(*types.DatabaseReport)) error

	Stop() error
}
