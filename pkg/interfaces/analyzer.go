// Package interfaces defines all the interfaces used throughout the database intelligence system.
// This promotes loose coupling and makes the system more testable and extensible.
package interfaces

import (
	"database/sql"
	"github.com/intelligent-algorithm/pkg/types"
)

// DatabaseAnalyzer defines the interface for database analysis operations.
type DatabaseAnalyzer interface {
	// AnalyzeTables analyzes all tables in the database.
	AnalyzeTables() ([]types.TableInfo, error)
	
	// AnalyzeTable performs detailed analysis of a single table.
	AnalyzeTable(tableName string) (types.TableInfo, error)
	
	// AnalyzeColumns analyzes all columns in a table.
	AnalyzeColumns(tableName string) ([]types.ColumnInfo, error)
	
	// AnalyzeColumnData creates a data profile for a column.
	AnalyzeColumnData(tableName, columnName, dataType string) (types.DataProfile, error)
	
	// GetTableNames retrieves all table names from the database.
	GetTableNames() ([]string, error)
	
	// GetRowCount returns the number of rows in a table.
	GetRowCount(tableName string) (int64, error)
	
	// GetUniqueValueCount returns the number of unique values in a column.
	GetUniqueValueCount(tableName, columnName string) (int64, error)
	
	// GetNullCount returns the number of null values in a column.
	GetNullCount(tableName, columnName string) (int64, error)
	
	// GetIndexes retrieves index information for a table.
	GetIndexes(tableName string) ([]types.IndexInfo, error)
	
	// GetConstraints retrieves constraint information for a table.
	GetConstraints(tableName string) ([]types.Constraint, error)
	
	// GetRelationships retrieves relationship information for a table.
	GetRelationships(tableName string) ([]types.Relationship, error)
	
	// GetTableSize returns the size of a table.
	GetTableSize(tableName string) (string, error)
	
	// CalculateDataQuality calculates a quality score for a column.
	CalculateDataQuality(tableName, columnName string) float64
}

// DatabaseConnector defines the interface for database connection operations.
type DatabaseConnector interface {
	// Connect establishes a connection to the database.
	Connect() error
	
	// Close closes the database connection.
	Close() error
	
	// Ping tests the database connection.
	Ping() error
	
	// GetDB returns the underlying database connection.
	GetDB() *sql.DB
	
	// GetDatabaseName returns the name of the database.
	GetDatabaseName() (string, error)
	
	// GetDatabaseType returns the type of the database (mysql, postgres, sqlite3).
	GetDatabaseType() string
}

// InsightGenerator defines the interface for generating database insights.
type InsightGenerator interface {
	// GenerateInsights creates intelligent insights about the database.
	GenerateInsights(tables []types.TableInfo) []types.DatabaseInsight
	
	// AnalyzeLargeTables identifies tables that may need attention due to size.
	AnalyzeLargeTables(tables []types.TableInfo) []types.DatabaseInsight
	
	// AnalyzeMissingIndexes identifies tables that might benefit from indexes.
	AnalyzeMissingIndexes(tables []types.TableInfo) []types.DatabaseInsight
	
	// AnalyzeDataQuality identifies data quality issues.
	AnalyzeDataQuality(tables []types.TableInfo) []types.DatabaseInsight
	
	// AnalyzeRelationships analyzes table relationships.
	AnalyzeRelationships(tables []types.TableInfo) []types.DatabaseInsight
	
	// AnalyzeUnusedColumns identifies potentially unused columns.
	AnalyzeUnusedColumns(tables []types.TableInfo) []types.DatabaseInsight
}

// ReportGenerator defines the interface for generating database reports.
type ReportGenerator interface {
	// GenerateSummary creates a high-level summary of the database.
	GenerateSummary(tables []types.TableInfo) types.DatabaseSummary
	
	// GenerateRecommendations creates actionable recommendations.
	GenerateRecommendations(tables []types.TableInfo, insights []types.DatabaseInsight) []string
	
	// CalculateHealthScore calculates an overall health score for the database.
	CalculateHealthScore(tables []types.TableInfo) float64
	
	// CalculateComplexityScore calculates a complexity score for the database.
	CalculateComplexityScore(tables []types.TableInfo) float64
	
	// ExportReport exports the analysis report to various formats.
	ExportReport(report *types.DatabaseReport, format string) ([]byte, error)
}

// PerformanceAnalyzer defines the interface for performance analysis.
type PerformanceAnalyzer interface {
	// AnalyzePerformance analyzes database performance metrics.
	AnalyzePerformance() types.PerformanceMetrics
}
