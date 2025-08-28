package interfaces

import (
	"context"

	"github.com/cherry-pick/pkg/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoConnector interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
	GetClient() *mongo.Client
	GetDatabase(name string) *mongo.Database
	GetDatabaseName() string
	GetConnectionString() string
}

type MongoAnalyzer interface {
	AnalyzeDatabase(ctx context.Context) (*types.DatabaseReport, error)
	AnalyzeCollections(ctx context.Context) ([]types.MongoCollectionInfo, error)
	AnalyzeCollection(ctx context.Context, collectionName string) (*types.MongoCollectionInfo, error)
	GetCollectionNames(ctx context.Context) ([]string, error)
	GetCollectionStats(ctx context.Context, collectionName string) (*types.MongoCollectionInfo, error)
	AnalyzeSchema(ctx context.Context, collectionName string, sampleSize int) ([]types.MongoFieldInfo, error)
	GetIndexes(ctx context.Context, collectionName string) ([]types.MongoIndexInfo, error)
	GetDatabaseStats(ctx context.Context) (*types.MongoDatabaseStats, error)
	GetPerformanceMetrics(ctx context.Context) (*types.MongoPerformanceMetrics, error)
}

type MongoInsightGenerator interface {
	GenerateInsights(collections []types.MongoCollectionInfo, stats *types.MongoDatabaseStats) []types.DatabaseInsight
	AnalyzeLargeCollections(collections []types.MongoCollectionInfo) []types.DatabaseInsight
	AnalyzeMissingIndexes(collections []types.MongoCollectionInfo) []types.DatabaseInsight
	AnalyzeSchemaIssues(collections []types.MongoCollectionInfo) []types.DatabaseInsight
	AnalyzePerformanceIssues(metrics *types.MongoPerformanceMetrics) []types.DatabaseInsight
}
