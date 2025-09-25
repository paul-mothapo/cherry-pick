package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type AnalyticsService struct {
	tracker   core.AnalyticsTracker
	processor core.AnalyticsProcessor
	reporter  core.AnalyticsReporter
	storage   core.AnalyticsStorage
	validator core.AnalyticsValidator
	notifier  core.AnalyticsNotifier
}

func (as *AnalyticsService) GetTracker() core.AnalyticsTracker {
	return as.tracker
}

func (as *AnalyticsService) GetProcessor() core.AnalyticsProcessor {
	return as.processor
}

func (as *AnalyticsService) GetReporter() core.AnalyticsReporter {
	return as.reporter
}

func (as *AnalyticsService) GetStorage() core.AnalyticsStorage {
	return as.storage
}

func NewAnalyticsService(
	tracker core.AnalyticsTracker,
	processor core.AnalyticsProcessor,
	reporter core.AnalyticsReporter,
	storage core.AnalyticsStorage,
	validator core.AnalyticsValidator,
	notifier core.AnalyticsNotifier,
) *AnalyticsService {
	return &AnalyticsService{
		tracker:   tracker,
		processor: processor,
		reporter:  reporter,
		storage:   storage,
		validator: validator,
		notifier:  notifier,
	}
}

func (as *AnalyticsService) TrackPageView(event core.PageViewEvent) error {
	return as.tracker.TrackPageView(event)
}

func (as *AnalyticsService) TrackBehavioralPattern(event core.BehavioralEvent) error {
	return as.tracker.TrackBehavioralPattern(event)
}

func (as *AnalyticsService) TrackPerformance(event core.PerformanceEvent) error {
	return as.tracker.TrackPerformance(event)
}

func (as *AnalyticsService) TrackCustomEvent(event core.AnalyticsEvent) error {
	return as.tracker.TrackCustomEvent(event)
}

func (as *AnalyticsService) CreateSession(session core.UserSession) error {
	return as.tracker.CreateSession(session)
}

func (as *AnalyticsService) GetSession(sessionID string) (*core.UserSession, error) {
	return as.tracker.GetSession(sessionID)
}

func (as *AnalyticsService) UpdateSession(session core.UserSession) error {
	return as.tracker.UpdateSession(session)
}

func (as *AnalyticsService) EndSession(sessionID string) error {
	return as.tracker.EndSession(sessionID)
}

func (as *AnalyticsService) GetUserJourney(sessionID string) (*core.UserJourney, error) {
	return as.processor.ProcessUserJourney(sessionID)
}

func (as *AnalyticsService) GetFunnelAnalysis(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error) {
	return as.processor.ProcessFunnelAnalysis(funnelID, startTime, endTime)
}

func (as *AnalyticsService) GetRealTimeMetrics() (*core.RealTimeMetrics, error) {
	return as.processor.ProcessRealTimeMetrics()
}

func (as *AnalyticsService) GetInsights(sessionID string) ([]core.AnalyticsInsight, error) {
	return as.processor.ProcessInsights(sessionID)
}

func (as *AnalyticsService) GetAlerts() ([]core.AnalyticsAlert, error) {
	return as.processor.ProcessAlerts()
}

func (as *AnalyticsService) GetHeatmapData(pagePath string, startTime, endTime time.Time) ([]core.HeatmapPoint, error) {
	return as.processor.ProcessHeatmapData(pagePath, startTime, endTime)
}

func (as *AnalyticsService) GenerateReport(request core.AnalyticsRequest) (*core.AnalyticsReport, error) {
	return as.reporter.GenerateReport(request)
}

func (as *AnalyticsService) GenerateSummary(startTime, endTime time.Time) (*core.AnalyticsSummary, error) {
	return as.reporter.GenerateSummary(startTime, endTime)
}

func (as *AnalyticsService) GenerateInsights(startTime, endTime time.Time) ([]core.AnalyticsInsight, error) {
	return as.reporter.GenerateInsights(startTime, endTime)
}

func (as *AnalyticsService) GenerateFunnelReport(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error) {
	return as.reporter.GenerateFunnelReport(funnelID, startTime, endTime)
}

func (as *AnalyticsService) GeneratePerformanceReport(startTime, endTime time.Time) ([]core.PerformanceEvent, error) {
	return as.reporter.GeneratePerformanceReport(startTime, endTime)
}

func (as *AnalyticsService) GenerateBehavioralReport(startTime, endTime time.Time) ([]core.BehavioralEvent, error) {
	return as.reporter.GenerateBehavioralReport(startTime, endTime)
}

func (as *AnalyticsService) SubscribeToRealTimeMetrics(ctx context.Context) (<-chan core.RealTimeMetrics, error) {
	metricsChan := make(chan core.RealTimeMetrics, 10)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		defer close(metricsChan)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics, err := as.GetRealTimeMetrics()
				if err != nil {
					continue
				}
				select {
				case metricsChan <- *metrics:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return metricsChan, nil
}

func (as *AnalyticsService) UnsubscribeFromRealTimeMetrics(subscriberID string) error {
	return nil
}

func (as *AnalyticsService) CleanupOldData(olderThan time.Time) error {
	return as.storage.CleanupOldData(olderThan)
}

func (as *AnalyticsService) GetStats() (map[string]interface{}, error) {
	return as.storage.GetStats()
}
