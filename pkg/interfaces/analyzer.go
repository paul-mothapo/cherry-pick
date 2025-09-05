package interfaces

import (
	"database/sql"

	"github.com/cherry-pick/pkg/types"
)

type DatabaseAnalyzer interface {
	AnalyzeTables() ([]types.TableInfo, error)
	AnalyzeTable(tableName string) (types.TableInfo, error)
	AnalyzeColumns(tableName string) ([]types.ColumnInfo, error)
	AnalyzeColumnData(tableName, columnName, dataType string) (types.DataProfile, error)
	GetTableNames() ([]string, error)
	GetRowCount(tableName string) (int64, error)
	GetUniqueValueCount(tableName, columnName string) (int64, error)
	GetNullCount(tableName, columnName string) (int64, error)
	GetIndexes(tableName string) ([]types.IndexInfo, error)
	GetConstraints(tableName string) ([]types.Constraint, error)
	GetRelationships(tableName string) ([]types.Relationship, error)
	GetTableSize(tableName string) (string, error)
	CalculateDataQuality(tableName, columnName string) float64
}

type DatabaseConnector interface {
	Connect() error
	Close() error
	Ping() error
	GetDB() *sql.DB
	GetDatabaseName() (string, error)
	GetDatabaseType() string
}

type InsightGenerator interface {
	GenerateInsights(tables []types.TableInfo) []types.DatabaseInsight
	AnalyzeLargeTables(tables []types.TableInfo) []types.DatabaseInsight
	AnalyzeMissingIndexes(tables []types.TableInfo) []types.DatabaseInsight
	AnalyzeDataQuality(tables []types.TableInfo) []types.DatabaseInsight
	AnalyzeRelationships(tables []types.TableInfo) []types.DatabaseInsight
	AnalyzeUnusedColumns(tables []types.TableInfo) []types.DatabaseInsight
}

type ReportGenerator interface {
	GenerateSummary(tables []types.TableInfo) types.DatabaseSummary
	GenerateRecommendations(tables []types.TableInfo, insights []types.DatabaseInsight) []string
	CalculateHealthScore(tables []types.TableInfo) float64
	CalculateComplexityScore(tables []types.TableInfo) float64
	ExportReport(report *types.DatabaseReport, format string) ([]byte, error)
}

type PerformanceAnalyzer interface {
	AnalyzePerformance() types.PerformanceMetrics
}
