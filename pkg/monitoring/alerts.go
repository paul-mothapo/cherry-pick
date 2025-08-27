// Package monitoring provides database monitoring and alerting functionality.
package monitoring

import (
	"fmt"
	"time"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
)

// AlertManagerImpl implements the AlertManager interface.
type AlertManagerImpl struct {
	alerts []types.MonitoringAlert
}

// NewAlertManager creates a new alert manager with default alerts.
func NewAlertManager() interfaces.AlertManager {
	defaultAlerts := []types.MonitoringAlert{
		{
			ID:        "large_table_growth",
			Name:      "Large Table Growth",
			Condition: "table_growth_rate > threshold",
			Threshold: 0.5, // 50% growth
			Severity:  "high",
		},
		{
			ID:        "data_quality_degradation",
			Name:      "Data Quality Degradation",
			Condition: "quality_score < threshold",
			Threshold: 0.7,
			Severity:  "medium",
		},
		{
			ID:        "missing_indexes",
			Name:      "Missing Indexes on Large Tables",
			Condition: "table_size > 10000 AND index_count <= 1",
			Threshold: 10000,
			Severity:  "high",
		},
	}

	return &AlertManagerImpl{alerts: defaultAlerts}
}

// CheckAlerts evaluates alerts against the current database report.
func (am *AlertManagerImpl) CheckAlerts(report *types.DatabaseReport) []types.MonitoringAlert {
	var triggeredAlerts []types.MonitoringAlert

	for i := range am.alerts {
		alert := &am.alerts[i]

		switch alert.ID {
		case "data_quality_degradation":
			if am.checkDataQualityAlert(alert, report) {
				triggeredAlerts = append(triggeredAlerts, *alert)
			}
		case "missing_indexes":
			if am.checkMissingIndexesAlert(alert, report) {
				triggeredAlerts = append(triggeredAlerts, *alert)
			}
		}
	}

	return triggeredAlerts
}

// checkDataQualityAlert checks for data quality degradation.
func (am *AlertManagerImpl) checkDataQualityAlert(alert *types.MonitoringAlert, report *types.DatabaseReport) bool {
	for _, table := range report.Tables {
		for _, column := range table.Columns {
			if column.DataProfile.Quality < alert.Threshold {
				alert.Triggered = true
				alert.LastTrigger = time.Now()
				alert.Message = fmt.Sprintf("Column %s.%s has quality score %.2f",
					table.Name, column.Name, column.DataProfile.Quality)
				return true
			}
		}
	}
	return false
}

// checkMissingIndexesAlert checks for missing indexes on large tables.
func (am *AlertManagerImpl) checkMissingIndexesAlert(alert *types.MonitoringAlert, report *types.DatabaseReport) bool {
	for _, table := range report.Tables {
		if float64(table.RowCount) > alert.Threshold && len(table.Indexes) <= 1 {
			alert.Triggered = true
			alert.LastTrigger = time.Now()
			alert.Message = fmt.Sprintf("Table %s has %d rows but only %d indexes",
				table.Name, table.RowCount, len(table.Indexes))
			return true
		}
	}
	return false
}

// AddAlert adds a new alert condition.
func (am *AlertManagerImpl) AddAlert(alert types.MonitoringAlert) error {
	// Check if alert with same ID already exists
	for _, existingAlert := range am.alerts {
		if existingAlert.ID == alert.ID {
			return fmt.Errorf("alert with ID %s already exists", alert.ID)
		}
	}

	am.alerts = append(am.alerts, alert)
	return nil
}

// RemoveAlert removes an alert condition.
func (am *AlertManagerImpl) RemoveAlert(alertID string) error {
	for i, alert := range am.alerts {
		if alert.ID == alertID {
			am.alerts = append(am.alerts[:i], am.alerts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("alert with ID %s not found", alertID)
}

// GetAlerts returns all configured alerts.
func (am *AlertManagerImpl) GetAlerts() []types.MonitoringAlert {
	// Return a copy to prevent external modification
	alerts := make([]types.MonitoringAlert, len(am.alerts))
	copy(alerts, am.alerts)
	return alerts
}
