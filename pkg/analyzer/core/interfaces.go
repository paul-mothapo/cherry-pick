package core

import (
	"context"
	"time"
)

type DatabaseAnalyzer interface {
	AnalyzeDatabase(ctx context.Context, request AnalysisRequest) (*AnalysisResult, error)
	AnalyzeTables(ctx context.Context, request AnalysisRequest) ([]TableInfo, error)
	AnalyzeTable(ctx context.Context, tableName string, request AnalysisRequest) (*TableInfo, error)
	GetTableNames(ctx context.Context, request AnalysisRequest) ([]string, error)
	GetTableStats(ctx context.Context, tableName string, request AnalysisRequest) (*TableInfo, error)
	GetPerformanceMetrics(ctx context.Context, request AnalysisRequest) (*PerformanceMetrics, error)
}

type MongoAnalyzer interface {
	AnalyzeDatabase(ctx context.Context, request AnalysisRequest) (*AnalysisResult, error)
	AnalyzeCollections(ctx context.Context, request AnalysisRequest) ([]MongoCollectionInfo, error)
	AnalyzeCollection(ctx context.Context, collectionName string, request AnalysisRequest) (*MongoCollectionInfo, error)
	GetCollectionNames(ctx context.Context, request AnalysisRequest) ([]string, error)
	GetCollectionStats(ctx context.Context, collectionName string, request AnalysisRequest) (*MongoCollectionInfo, error)
	AnalyzeSchema(ctx context.Context, collectionName string, request AnalysisRequest) ([]MongoFieldInfo, error)
	GetIndexes(ctx context.Context, collectionName string, request AnalysisRequest) ([]MongoIndexInfo, error)
	GetDatabaseStats(ctx context.Context, request AnalysisRequest) (*MongoDatabaseStats, error)
	GetPerformanceMetrics(ctx context.Context, request AnalysisRequest) (*PerformanceMetrics, error)
}

type AnalyzerService interface {
	AnalyzeDatabase(ctx context.Context, request AnalysisRequest) (*AnalysisResult, error)
	GetAnalysisHistory(ctx context.Context, limit int) ([]AnalysisResult, error)
	GetAnalysisByID(ctx context.Context, analysisID string) (*AnalysisResult, error)
	DeleteAnalysis(ctx context.Context, analysisID string) error
	GetSupportedDatabaseTypes() []DatabaseType
	GetAnalysisOptions() AnalysisOptions
}

type DatabaseConnector interface {
	Connect(ctx context.Context, connectionString string) error
	Disconnect(ctx context.Context) error
	IsConnected() bool
	GetDatabaseName() string
	GetDatabaseType() DatabaseType
	TestConnection(ctx context.Context) error
}

type MongoConnector interface {
	Connect(ctx context.Context, connectionString string) error
	Disconnect(ctx context.Context) error
	IsConnected() bool
	GetDatabaseName() string
	GetDatabase() interface{}
	TestConnection(ctx context.Context) error
}

type AnalysisStorage interface {
	SaveAnalysis(ctx context.Context, result *AnalysisResult) error
	GetAnalysis(ctx context.Context, analysisID string) (*AnalysisResult, error)
	GetAnalysisHistory(ctx context.Context, limit int) ([]AnalysisResult, error)
	DeleteAnalysis(ctx context.Context, analysisID string) error
	CleanupOldAnalyses(ctx context.Context, olderThan time.Duration) error
}

type AnalysisReporter interface {
	GenerateReport(ctx context.Context, result *AnalysisResult) (map[string]string, error)
	GenerateSummary(ctx context.Context, result *AnalysisResult) (string, error)
	GenerateInsights(ctx context.Context, result *AnalysisResult) ([]DatabaseInsight, error)
	GenerateRecommendations(ctx context.Context, result *AnalysisResult) ([]string, error)
}

type AnalysisValidator interface {
	ValidateRequest(request AnalysisRequest) error
	ValidateDatabaseType(dbType DatabaseType) error
	ValidateOptions(options AnalysisOptions) error
}

type AnalysisCalculator interface {
	CalculateHealthScore(tables []TableInfo) float64
	CalculateComplexityScore(tables []TableInfo) float64
	CalculateDataQuality(column ColumnInfo) float64
	CalculatePerformanceScore(metrics *PerformanceMetrics) float64
}

type AnalysisAggregator interface {
	AggregateTableStats(tables []TableInfo) DatabaseSummary
	AggregateInsights(insights []DatabaseInsight) []DatabaseInsight
	AggregateRecommendations(recommendations []string) []string
}

type AnalysisNotifier interface {
	NotifyAnalysisComplete(ctx context.Context, result *AnalysisResult) error
	NotifyAnalysisError(ctx context.Context, err error) error
	NotifyInsight(ctx context.Context, insight DatabaseInsight) error
}
