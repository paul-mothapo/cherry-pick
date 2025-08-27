package types

import "time"

type DatabaseReport struct {
	DatabaseName       string             `json:"database_name"`
	DatabaseType       string             `json:"database_type"`
	AnalysisTime       time.Time          `json:"analysis_time"`
	Summary            DatabaseSummary    `json:"summary"`
	Tables             []TableInfo        `json:"tables"`
	Insights           []DatabaseInsight  `json:"insights"`
	Recommendations    []string           `json:"recommendations"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
}

type DatabaseSummary struct {
	TotalTables     int     `json:"total_tables"`
	TotalColumns    int     `json:"total_columns"`
	TotalRows       int64   `json:"total_rows"`
	TotalSize       string  `json:"total_size"`
	HealthScore     float64 `json:"health_score"`
	ComplexityScore float64 `json:"complexity_score"`
}

type TableInfo struct {
	Name          string         `json:"name"`
	RowCount      int64          `json:"row_count"`
	Columns       []ColumnInfo   `json:"columns"`
	Indexes       []IndexInfo    `json:"indexes"`
	Constraints   []Constraint   `json:"constraints"`
	Size          string         `json:"size"`
	LastModified  time.Time      `json:"last_modified"`
	Relationships []Relationship `json:"relationships"`
}

type ColumnInfo struct {
	Name         string      `json:"name"`
	DataType     string      `json:"data_type"`
	IsNullable   bool        `json:"is_nullable"`
	IsPrimaryKey bool        `json:"is_primary_key"`
	DefaultValue interface{} `json:"default_value"`
	MaxLength    int         `json:"max_length"`
	Precision    int         `json:"precision"`
	Scale        int         `json:"scale"`
	UniqueValues int64       `json:"unique_values"`
	NullCount    int64       `json:"null_count"`
	DataProfile  DataProfile `json:"data_profile"`
}

type DataProfile struct {
	Min        interface{} `json:"min"`
	Max        interface{} `json:"max"`
	Avg        interface{} `json:"avg"`
	SampleData []string    `json:"sample_data"`
	Pattern    string      `json:"pattern"`
	Quality    float64     `json:"quality_score"`
}

type IndexInfo struct {
	Name     string   `json:"name"`
	Columns  []string `json:"columns"`
	IsUnique bool     `json:"is_unique"`
	Type     string   `json:"type"`
}

type Constraint struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Columns    []string `json:"columns"`
	RefTable   string   `json:"ref_table,omitempty"`
	RefColumns []string `json:"ref_columns,omitempty"`
}

type Relationship struct {
	Type         string `json:"type"`
	TargetTable  string `json:"target_table"`
	SourceColumn string `json:"source_column"`
	TargetColumn string `json:"target_column"`
}

type DatabaseInsight struct {
	Type           string      `json:"type"`
	Severity       string      `json:"severity"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Suggestion     string      `json:"suggestion"`
	AffectedTables []string    `json:"affected_tables"`
	MetricValue    interface{} `json:"metric_value"`
}

type PerformanceMetrics struct {
	SlowQueries     []SlowQuery `json:"slow_queries"`
	IndexUsage      float64     `json:"index_usage_ratio"`
	TableScanRatio  float64     `json:"table_scan_ratio"`
	ConnectionCount int         `json:"connection_count"`
	BufferHitRatio  float64     `json:"buffer_hit_ratio"`
}

type SlowQuery struct {
	Query     string        `json:"query"`
	Duration  time.Duration `json:"duration"`
	Frequency int           `json:"frequency"`
	Tables    []string      `json:"tables"`
}
