package core

import "time"

type DatabaseType string

const (
	DatabaseTypeMySQL    DatabaseType = "mysql"
	DatabaseTypePostgres DatabaseType = "postgres"
	DatabaseTypeSQLite   DatabaseType = "sqlite3"
	DatabaseTypeMongoDB  DatabaseType = "mongodb"
)

type AnalysisRequest struct {
	DatabaseType DatabaseType `json:"databaseType"`
	ConnectionID string       `json:"connectionId"`
	Options      AnalysisOptions `json:"options"`
}

type AnalysisOptions struct {
	IncludeSchema     bool `json:"includeSchema"`
	IncludeData       bool `json:"includeData"`
	IncludeIndexes    bool `json:"includeIndexes"`
	IncludeRelations  bool `json:"includeRelations"`
	IncludePerformance bool `json:"includePerformance"`
	SampleSize        int  `json:"sampleSize"`
	MaxCollections    int  `json:"maxCollections"`
}

type AnalysisResult struct {
	ID            string          `json:"id"`
	DatabaseName  string          `json:"databaseName"`
	DatabaseType  DatabaseType    `json:"databaseType"`
	AnalysisTime  time.Time       `json:"analysisTime"`
	Summary       DatabaseSummary `json:"summary"`
	Tables        []TableInfo     `json:"tables"`
	Insights      []DatabaseInsight `json:"insights"`
	Recommendations []string       `json:"recommendations"`
	Performance   *PerformanceMetrics `json:"performance,omitempty"`
}

type DatabaseSummary struct {
	TotalTables     int     `json:"totalTables"`
	TotalColumns    int     `json:"totalColumns"`
	TotalRows       int64   `json:"totalRows"`
	TotalSize       string  `json:"totalSize"`
	HealthScore     float64 `json:"healthScore"`
	ComplexityScore float64 `json:"complexityScore"`
}

type TableInfo struct {
	Name         string        `json:"name"`
	RowCount     int64         `json:"rowCount"`
	Size         string        `json:"size"`
	LastModified time.Time     `json:"lastModified"`
	Columns      []ColumnInfo  `json:"columns"`
	Indexes      []IndexInfo   `json:"indexes"`
	Constraints  []Constraint  `json:"constraints"`
	Relationships []Relationship `json:"relationships"`
}

type ColumnInfo struct {
	Name         string      `json:"name"`
	DataType     string      `json:"dataType"`
	IsNullable   bool        `json:"isNullable"`
	IsPrimaryKey bool        `json:"isPrimaryKey"`
	DefaultValue string      `json:"defaultValue,omitempty"`
	MaxLength    int         `json:"maxLength,omitempty"`
	Precision    int         `json:"precision,omitempty"`
	Scale        int         `json:"scale,omitempty"`
	DataProfile  DataProfile `json:"dataProfile"`
	UniqueValues int64       `json:"uniqueValues"`
	NullCount    int64       `json:"nullCount"`
}

type DataProfile struct {
	SampleData []string  `json:"sampleData,omitempty"`
	Min        float64   `json:"min,omitempty"`
	Max        float64   `json:"max,omitempty"`
	Avg        float64   `json:"avg,omitempty"`
	Pattern    string    `json:"pattern,omitempty"`
	Quality    float64   `json:"quality"`
}

type IndexInfo struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	IsUnique bool     `json:"isUnique"`
	Columns  []string `json:"columns"`
}

type Constraint struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Columns    []string `json:"columns"`
	RefTable   string   `json:"refTable,omitempty"`
	RefColumns []string `json:"refColumns,omitempty"`
}

type Relationship struct {
	SourceColumn string `json:"sourceColumn"`
	TargetTable  string `json:"targetTable"`
	TargetColumn string `json:"targetColumn"`
	Type         string `json:"type"`
}

type DatabaseInsight struct {
	Type           string    `json:"type"`
	Severity       string    `json:"severity"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Suggestion     string    `json:"suggestion"`
	AffectedTables []string  `json:"affectedTables"`
	MetricValue    interface{} `json:"metricValue"`
}

type PerformanceMetrics struct {
	Connections ConnectionMetrics `json:"connections"`
	Operations  OperationMetrics  `json:"operations"`
	Memory      MemoryMetrics     `json:"memory"`
	Storage     StorageMetrics    `json:"storage"`
}

type ConnectionMetrics struct {
	Current      int `json:"current"`
	Available    int `json:"available"`
	TotalCreated int `json:"totalCreated"`
}

type OperationMetrics struct {
	Insert  int64 `json:"insert"`
	Query   int64 `json:"query"`
	Update  int64 `json:"update"`
	Delete  int64 `json:"delete"`
	GetMore int64 `json:"getMore"`
	Command int64 `json:"command"`
}

type MemoryMetrics struct {
	Resident int64 `json:"resident"`
	Virtual  int64 `json:"virtual"`
	Mapped   int64 `json:"mapped"`
}

type StorageMetrics struct {
	DataSize    int64 `json:"dataSize"`
	IndexSize   int64 `json:"indexSize"`
	StorageSize int64 `json:"storageSize"`
}

type MongoCollectionInfo struct {
	Name           string                 `json:"name"`
	DocumentCount  int64                  `json:"documentCount"`
	TotalSize      int64                  `json:"totalSize"`
	StorageSize    int64                  `json:"storageSize"`
	AvgDocSize     int64                  `json:"avgDocSize"`
	LastModified   time.Time              `json:"lastModified"`
	Fields         []MongoFieldInfo       `json:"fields"`
	Indexes        []MongoIndexInfo       `json:"indexes"`
	SampleDocument map[string]interface{} `json:"sampleDocument,omitempty"`
	IsSharded      bool                   `json:"isSharded"`
}

type MongoFieldInfo struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Frequency   float64     `json:"frequency"`
	SampleValue interface{} `json:"sampleValue,omitempty"`
}

type MongoIndexInfo struct {
	Name    string                 `json:"name"`
	Keys    map[string]interface{} `json:"keys"`
	IsUnique bool                  `json:"isUnique"`
	IsSparse bool                 `json:"isSparse"`
	IsPartial bool                `json:"isPartial"`
	UsageStats MongoIndexUsageStats `json:"usageStats"`
}

type MongoIndexUsageStats struct {
	Since time.Time `json:"since"`
}

type MongoDatabaseStats struct {
	Name        string `json:"name"`
	Collections int    `json:"collections"`
	Views       int    `json:"views"`
	Objects     int64  `json:"objects"`
	AvgObjSize  float64 `json:"avgObjSize"`
	DataSize    int64  `json:"dataSize"`
	StorageSize int64  `json:"storageSize"`
	IndexSize   int64  `json:"indexSize"`
	TotalSize   int64  `json:"totalSize"`
}
