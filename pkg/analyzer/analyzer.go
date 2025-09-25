package analyzer

import (
	"github.com/cherry-pick/pkg/analyzer/core"
	"github.com/cherry-pick/pkg/analyzer/services"
	"github.com/cherry-pick/pkg/analyzer/storage"
)

type Analyzer struct {
	service core.AnalyzerService
}

func NewAnalyzer() *Analyzer {
	storage := storage.NewMemoryStorage()
	validator := services.NewValidatorService()
	calculator := services.NewCalculatorService()
	aggregator := services.NewAggregatorService()
	reporter := services.NewReporterService()
	notifier := services.NewNotifierService()

	databaseAnalyzer := services.NewDatabaseAnalyzerService(nil, calculator, aggregator, validator)
	mongoAnalyzer := services.NewMongoAnalyzerService(nil, calculator, aggregator, validator)

	service := services.NewAnalyzerService(
		databaseAnalyzer,
		mongoAnalyzer,
		storage,
		reporter,
		validator,
		notifier,
	)

	return &Analyzer{
		service: service,
	}
}

func (a *Analyzer) GetService() core.AnalyzerService {
	return a.service
}
