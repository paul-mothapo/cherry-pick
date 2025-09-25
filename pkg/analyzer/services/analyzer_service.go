package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type AnalyzerService struct {
	databaseAnalyzer core.DatabaseAnalyzer
	mongoAnalyzer    core.MongoAnalyzer
	storage          core.AnalysisStorage
	reporter         core.AnalysisReporter
	validator        core.AnalysisValidator
	notifier         core.AnalysisNotifier
}

func NewAnalyzerService(
	databaseAnalyzer core.DatabaseAnalyzer,
	mongoAnalyzer core.MongoAnalyzer,
	storage core.AnalysisStorage,
	reporter core.AnalysisReporter,
	validator core.AnalysisValidator,
	notifier core.AnalysisNotifier,
) *AnalyzerService {
	return &AnalyzerService{
		databaseAnalyzer: databaseAnalyzer,
		mongoAnalyzer:    mongoAnalyzer,
		storage:          storage,
		reporter:         reporter,
		validator:        validator,
		notifier:         notifier,
	}
}

func (as *AnalyzerService) AnalyzeDatabase(ctx context.Context, request core.AnalysisRequest) (*core.AnalysisResult, error) {
	if err := as.validator.ValidateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid analysis request: %w", err)
	}

	var result *core.AnalysisResult
	var err error

	switch request.DatabaseType {
	case core.DatabaseTypeMongoDB:
		result, err = as.mongoAnalyzer.AnalyzeDatabase(ctx, request)
	default:
		result, err = as.databaseAnalyzer.AnalyzeDatabase(ctx, request)
	}

	if err != nil {
		as.notifier.NotifyAnalysisError(ctx, err)
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	if err := as.storage.SaveAnalysis(ctx, result); err != nil {
		return result, fmt.Errorf("failed to save analysis: %w", err)
	}

	as.notifier.NotifyAnalysisComplete(ctx, result)

	return result, nil
}

func (as *AnalyzerService) GetAnalysisHistory(ctx context.Context, limit int) ([]core.AnalysisResult, error) {
	return as.storage.GetAnalysisHistory(ctx, limit)
}

func (as *AnalyzerService) GetAnalysisByID(ctx context.Context, analysisID string) (*core.AnalysisResult, error) {
	return as.storage.GetAnalysis(ctx, analysisID)
}

func (as *AnalyzerService) DeleteAnalysis(ctx context.Context, analysisID string) error {
	return as.storage.DeleteAnalysis(ctx, analysisID)
}

func (as *AnalyzerService) GetSupportedDatabaseTypes() []core.DatabaseType {
	return []core.DatabaseType{
		core.DatabaseTypeMySQL,
		core.DatabaseTypePostgres,
		core.DatabaseTypeSQLite,
		core.DatabaseTypeMongoDB,
	}
}

func (as *AnalyzerService) GetAnalysisOptions() core.AnalysisOptions {
	return core.AnalysisOptions{
		IncludeSchema:     true,
		IncludeData:       true,
		IncludeIndexes:    true,
		IncludeRelations:  true,
		IncludePerformance: true,
		SampleSize:        100,
		MaxCollections:    50,
	}
}
