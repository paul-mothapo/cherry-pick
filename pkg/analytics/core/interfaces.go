package core

import (
	"context"
	"time"
)

type AnalyticsTracker interface {
	TrackPageView(event PageViewEvent) error
	TrackBehavioralPattern(event BehavioralEvent) error
	TrackPerformance(event PerformanceEvent) error
	TrackCustomEvent(event AnalyticsEvent) error
	GetSession(sessionID string) (*UserSession, error)
	CreateSession(session UserSession) error
	UpdateSession(session UserSession) error
	EndSession(sessionID string) error
}

type AnalyticsProcessor interface {
	ProcessUserJourney(sessionID string) (*UserJourney, error)
	ProcessFunnelAnalysis(funnelID string, startTime, endTime time.Time) (*FunnelAnalysis, error)
	ProcessRealTimeMetrics() (*RealTimeMetrics, error)
	ProcessInsights(sessionID string) ([]AnalyticsInsight, error)
	ProcessAlerts() ([]AnalyticsAlert, error)
	ProcessHeatmapData(pagePath string, startTime, endTime time.Time) ([]HeatmapPoint, error)
}

type AnalyticsReporter interface {
	GenerateReport(request AnalyticsRequest) (*AnalyticsReport, error)
	GenerateSummary(startTime, endTime time.Time) (*AnalyticsSummary, error)
	GenerateInsights(startTime, endTime time.Time) ([]AnalyticsInsight, error)
	GenerateFunnelReport(funnelID string, startTime, endTime time.Time) (*FunnelAnalysis, error)
	GeneratePerformanceReport(startTime, endTime time.Time) ([]PerformanceEvent, error)
	GenerateBehavioralReport(startTime, endTime time.Time) ([]BehavioralEvent, error)
}

type AnalyticsStorage interface {
	SaveEvent(event AnalyticsEvent) error
	SaveSession(session UserSession) error
	SaveJourney(journey UserJourney) error
	SaveInsight(insight AnalyticsInsight) error
	SaveAlert(alert AnalyticsAlert) error
	SaveReport(report AnalyticsReport) error

	GetEvents(request AnalyticsRequest) ([]AnalyticsEvent, error)
	GetSessions(request AnalyticsRequest) ([]UserSession, error)
	GetJourneys(request AnalyticsRequest) ([]UserJourney, error)
	GetInsights(request AnalyticsRequest) ([]AnalyticsInsight, error)
	GetAlerts(request AnalyticsRequest) ([]AnalyticsAlert, error)
	GetReports(request AnalyticsRequest) ([]AnalyticsReport, error)

	GetSession(sessionID string) (*UserSession, error)
	GetJourney(sessionID string) (*UserJourney, error)
	GetInsight(insightID string) (*AnalyticsInsight, error)
	GetAlert(alertID string) (*AnalyticsAlert, error)
	GetReport(reportID string) (*AnalyticsReport, error)

	UpdateSession(session UserSession) error
	UpdateJourney(journey UserJourney) error
	UpdateInsight(insight AnalyticsInsight) error
	UpdateAlert(alert AnalyticsAlert) error

	DeleteEvent(eventID string) error
	DeleteSession(sessionID string) error
	DeleteJourney(sessionID string) error
	DeleteInsight(insightID string) error
	DeleteAlert(alertID string) error
	DeleteReport(reportID string) error

	CleanupOldData(olderThan time.Time) error
	GetStats() (map[string]interface{}, error)
}

type AnalyticsService interface {
	TrackPageView(event PageViewEvent) error
	TrackBehavioralPattern(event BehavioralEvent) error
	TrackPerformance(event PerformanceEvent) error
	TrackCustomEvent(event AnalyticsEvent) error

	CreateSession(session UserSession) error
	GetSession(sessionID string) (*UserSession, error)
	UpdateSession(session UserSession) error
	EndSession(sessionID string) error

	GetUserJourney(sessionID string) (*UserJourney, error)
	GetFunnelAnalysis(funnelID string, startTime, endTime time.Time) (*FunnelAnalysis, error)
	GetRealTimeMetrics() (*RealTimeMetrics, error)
	GetInsights(sessionID string) ([]AnalyticsInsight, error)
	GetAlerts() ([]AnalyticsAlert, error)
	GetHeatmapData(pagePath string, startTime, endTime time.Time) ([]HeatmapPoint, error)

	GenerateReport(request AnalyticsRequest) (*AnalyticsReport, error)
	GenerateSummary(startTime, endTime time.Time) (*AnalyticsSummary, error)
	GenerateInsights(startTime, endTime time.Time) ([]AnalyticsInsight, error)
	GenerateFunnelReport(funnelID string, startTime, endTime time.Time) (*FunnelAnalysis, error)
	GeneratePerformanceReport(startTime, endTime time.Time) ([]PerformanceEvent, error)
	GenerateBehavioralReport(startTime, endTime time.Time) ([]BehavioralEvent, error)

	SubscribeToRealTimeMetrics(ctx context.Context) (<-chan RealTimeMetrics, error)
	UnsubscribeFromRealTimeMetrics(subscriberID string) error

	CleanupOldData(olderThan time.Time) error
	GetStats() (map[string]interface{}, error)
}

type AnalyticsValidator interface {
	ValidateEvent(event AnalyticsEvent) error
	ValidateSession(session UserSession) error
	ValidateJourney(journey UserJourney) error
	ValidateRequest(request AnalyticsRequest) error
}

type AnalyticsNotifier interface {
	SendAlert(alert AnalyticsAlert) error
	SendInsight(insight AnalyticsInsight) error
	SendReport(report AnalyticsReport) error
}

type AnalyticsAggregator interface {
	AggregatePageViews(startTime, endTime time.Time) ([]PageStats, error)
	AggregateUserSessions(startTime, endTime time.Time) ([]UserSession, error)
	AggregatePerformanceMetrics(startTime, endTime time.Time) ([]PerformanceEvent, error)
	AggregateBehavioralPatterns(startTime, endTime time.Time) ([]BehavioralEvent, error)
	AggregateFunnelData(funnelID string, startTime, endTime time.Time) (*FunnelAnalysis, error)
}

type AnalyticsCalculator interface {
	CalculateBounceRate(sessionID string) (float64, error)
	CalculateConversionRate(funnelID string, startTime, endTime time.Time) (float64, error)
	CalculatePerformanceScore(events []PerformanceEvent) (float64, error)
	CalculateUserEngagement(sessionID string) (float64, error)
	CalculateFunnelDropOff(funnelID string, startTime, endTime time.Time) ([]float64, error)
}

type AnalyticsExporter interface {
	ExportToCSV(data interface{}, filename string) error
	ExportToJSON(data interface{}, filename string) error
	ExportToExcel(data interface{}, filename string) error
	ExportToPDF(report AnalyticsReport, filename string) error
}

type AnalyticsImporter interface {
	ImportFromCSV(filename string) ([]AnalyticsEvent, error)
	ImportFromJSON(filename string) ([]AnalyticsEvent, error)
	ImportFromExcel(filename string) ([]AnalyticsEvent, error)
}

type AnalyticsScheduler interface {
	ScheduleReport(reportID string, cronExpression string) error
	ScheduleCleanup(cronExpression string) error
	ScheduleInsightGeneration(cronExpression string) error
	CancelSchedule(scheduleID string) error
	GetSchedules() ([]string, error)
}
