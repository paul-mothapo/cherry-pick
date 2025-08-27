package types

import "time"

type ComparisonReport struct {
	OldAnalysisTime time.Time     `json:"old_analysis_time"`
	NewAnalysisTime time.Time     `json:"new_analysis_time"`
	Changes         []DatabaseChange `json:"changes"`
	Summary         ChangeSummary `json:"summary"`
}

type DatabaseChange struct {
	Type          string      `json:"type"`
	Category      string      `json:"category"`
	Description   string      `json:"description"`
	Impact        string      `json:"impact"`
	AffectedTable string      `json:"affected_table,omitempty"`
	OldValue      interface{} `json:"old_value,omitempty"`
	NewValue      interface{} `json:"new_value,omitempty"`
}

type ChangeSummary struct {
	TotalChanges  int `json:"total_changes"`
	SchemaChanges int `json:"schema_changes"`
	DataChanges   int `json:"data_changes"`
	HighImpact    int `json:"high_impact_changes"`
	MediumImpact  int `json:"medium_impact_changes"`
	LowImpact     int `json:"low_impact_changes"`
}

type MonitoringAlert struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Condition   string    `json:"condition"`
	Threshold   float64   `json:"threshold"`
	Severity    string    `json:"severity"`
	Triggered   bool      `json:"triggered"`
	LastTrigger time.Time `json:"last_trigger"`
	Message     string    `json:"message"`
}

type DataLineage struct {
	TableName      string              `json:"table_name"`
	ColumnName     string              `json:"column_name,omitempty"`
	UpstreamDeps   []LineageDependency `json:"upstream_dependencies"`
	DownstreamDeps []LineageDependency `json:"downstream_dependencies"`
	LastUpdated    time.Time           `json:"last_updated"`
}

type LineageDependency struct {
	TableName  string `json:"table_name"`
	ColumnName string `json:"column_name,omitempty"`
	Type       string `json:"type"`
}
