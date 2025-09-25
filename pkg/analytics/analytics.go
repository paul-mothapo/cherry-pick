
package analytics

import (
	"github.com/cherry-pick/pkg/analytics/core"
	"github.com/cherry-pick/pkg/analytics/services"
	"github.com/cherry-pick/pkg/analytics/storage"
)

type Analytics struct {
	service core.AnalyticsService
}

func NewAnalytics() *Analytics {
	storage := storage.NewMemoryStorage()
	validator := services.NewValidatorService()
	calculator := services.NewCalculatorService(storage)
	aggregator := services.NewAggregatorService(storage)
	tracker := services.NewTrackerService(storage, validator)
	processor := services.NewProcessorService(storage, calculator, aggregator)
	reporter := services.NewReporterService(storage, processor, aggregator)
	notifier := services.NewNotifierService()
	service := services.NewAnalyticsService(tracker, processor, reporter, storage, validator, notifier)
	return &Analytics{
		service: service,
	}
}

func (a *Analytics) GetService() core.AnalyticsService {
	return a.service
}

func (a *Analytics) GetTracker() core.AnalyticsTracker {
	if service, ok := a.service.(*services.AnalyticsService); ok {
		return service.GetTracker()
	}
	return nil
}

func (a *Analytics) GetProcessor() core.AnalyticsProcessor {
	if service, ok := a.service.(*services.AnalyticsService); ok {
		return service.GetProcessor()
	}
	return nil
}

func (a *Analytics) GetReporter() core.AnalyticsReporter {
	if service, ok := a.service.(*services.AnalyticsService); ok {
		return service.GetReporter()
	}
	return nil
}

func (a *Analytics) GetStorage() core.AnalyticsStorage {
	if service, ok := a.service.(*services.AnalyticsService); ok {
		return service.GetStorage()
	}
	return nil
}
