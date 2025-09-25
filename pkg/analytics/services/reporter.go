package services

import (
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type ReporterService struct {
	storage    core.AnalyticsStorage
	processor  core.AnalyticsProcessor
	aggregator core.AnalyticsAggregator
}

func NewReporterService(
	storage core.AnalyticsStorage,
	processor core.AnalyticsProcessor,
	aggregator core.AnalyticsAggregator,
) *ReporterService {
	return &ReporterService{
		storage:    storage,
		processor:  processor,
		aggregator: aggregator,
	}
}

func (rs *ReporterService) GenerateReport(request core.AnalyticsRequest) (*core.AnalyticsReport, error) {
	if request.StartTime == nil {
		startTime := time.Now().Add(-24 * time.Hour)
		request.StartTime = &startTime
	}
	if request.EndTime == nil {
		endTime := time.Now()
		request.EndTime = &endTime
	}

	summary, err := rs.GenerateSummary(*request.StartTime, *request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	journeys, err := rs.storage.GetJourneys(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get user journeys: %w", err)
	}

	var funnelAnalysis *core.FunnelAnalysis
	if request.Filters != nil {
		if funnelID, ok := request.Filters["funnel_id"].(string); ok {
			funnelAnalysis, err = rs.ProcessFunnelAnalysis(funnelID, *request.StartTime, *request.EndTime)
			if err != nil {
				return nil, fmt.Errorf("failed to process funnel analysis: %w", err)
			}
		}
	}

	performanceData, err := rs.GeneratePerformanceReport(*request.StartTime, *request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate performance report: %w", err)
	}

	behavioralData, err := rs.GenerateBehavioralReport(*request.StartTime, *request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate behavioral report: %w", err)
	}

	insights, err := rs.GenerateInsights(*request.StartTime, *request.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate insights: %w", err)
	}

	report := &core.AnalyticsReport{
		ID:              generateReportID(),
		Title:           "Analytics Report",
		Period:          "custom",
		StartTime:       *request.StartTime,
		EndTime:         *request.EndTime,
		GeneratedAt:     time.Now(),
		Summary:         *summary,
		UserJourneys:    journeys,
		FunnelAnalysis:  funnelAnalysis,
		PerformanceData: performanceData,
		BehavioralData:  behavioralData,
		Insights:        insights,
		Recommendations: rs.generateRecommendations(summary, insights),
		Metadata: map[string]interface{}{
			"generated_by": "analytics_service",
			"version":      "1.0.0",
		},
	}

	if err := rs.storage.SaveReport(*report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return report, nil
}

func (rs *ReporterService) GenerateSummary(startTime, endTime time.Time) (*core.AnalyticsSummary, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	sessions, err := rs.storage.GetSessions(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	events, err := rs.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	summary := &core.AnalyticsSummary{
		TotalSessions: len(sessions),
		UniqueUsers:   rs.countUniqueUsers(sessions),
	}

	pageViewCount := 0
	for _, event := range events {
		if event.Type == "page_view" {
			pageViewCount++
		}
	}
	summary.TotalPageViews = pageViewCount

	var totalDuration int64
	for _, session := range sessions {
		if session.EndTime != nil {
			totalDuration += session.EndTime.Sub(session.StartTime).Milliseconds()
		}
	}
	if len(sessions) > 0 {
		summary.AvgSessionDuration = totalDuration / int64(len(sessions))
	}

	summary.BounceRate = rs.calculateBounceRate(sessions)

	var totalLoadTime int64
	var loadTimeCount int
	for _, event := range events {
		if event.Type == "page_view" {
			if loadTime, ok := event.Metadata["loadTime"].(int64); ok {
				totalLoadTime += loadTime
				loadTimeCount++
			}
		}
	}
	if loadTimeCount > 0 {
		summary.AvgPageLoadTime = totalLoadTime / int64(loadTimeCount)
	}

	summary.TopPage = rs.getTopPage(events)
	summary.TopReferrer = rs.getTopReferrer(events)
	summary.TopCountry = rs.getTopCountry(sessions)
	summary.TopDevice = rs.getTopDevice(sessions)
	summary.TopBrowser = rs.getTopBrowser(sessions)

	return summary, nil
}

func (rs *ReporterService) GenerateInsights(startTime, endTime time.Time) ([]core.AnalyticsInsight, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	insights, err := rs.storage.GetInsights(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get insights: %w", err)
	}

	additionalInsights := rs.generateAdditionalInsights(startTime, endTime)
	insights = append(insights, additionalInsights...)

	return insights, nil
}

func (rs *ReporterService) GenerateFunnelReport(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error) {
	return rs.ProcessFunnelAnalysis(funnelID, startTime, endTime)
}

func (rs *ReporterService) GeneratePerformanceReport(startTime, endTime time.Time) ([]core.PerformanceEvent, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"type": "performance",
		},
	}
	events, err := rs.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance events: %w", err)
	}

	var performanceEvents []core.PerformanceEvent
	for _, event := range events {
		perfEvent := core.PerformanceEvent{
			AnalyticsEvent: event,
		}
		performanceEvents = append(performanceEvents, perfEvent)
	}

	return performanceEvents, nil
}

func (rs *ReporterService) GenerateBehavioralReport(startTime, endTime time.Time) ([]core.BehavioralEvent, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"type": "behavioral",
		},
	}
	events, err := rs.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get behavioral events: %w", err)
	}

	var behavioralEvents []core.BehavioralEvent
	for _, event := range events {
		behaviorEvent := core.BehavioralEvent{
			AnalyticsEvent: event,
		}
		behavioralEvents = append(behavioralEvents, behaviorEvent)
	}

	return behavioralEvents, nil
}

func (rs *ReporterService) countUniqueUsers(sessions []core.UserSession) int {
	userMap := make(map[string]bool)
	for _, session := range sessions {
		if session.UserID != "" {
			userMap[session.UserID] = true
		}
	}
	return len(userMap)
}

func (rs *ReporterService) calculateBounceRate(sessions []core.UserSession) float64 {
	if len(sessions) == 0 {
		return 0
	}

	bounceCount := 0
	for _, session := range sessions {
		if session.EndTime != nil && session.StartTime.Add(time.Minute).After(*session.EndTime) {
			bounceCount++
		}
	}

	return float64(bounceCount) / float64(len(sessions))
}

func (rs *ReporterService) getTopPage(events []core.AnalyticsEvent) string {
	pageCounts := make(map[string]int)
	for _, event := range events {
		if event.Type == "page_view" {
			if path, ok := event.Metadata["path"].(string); ok {
				pageCounts[path]++
			}
		}
	}

	var topPage string
	var maxCount int
	for page, count := range pageCounts {
		if count > maxCount {
			maxCount = count
			topPage = page
		}
	}

	return topPage
}

func (rs *ReporterService) getTopReferrer(events []core.AnalyticsEvent) string {
	referrerCounts := make(map[string]int)
	for _, event := range events {
		if event.Type == "page_view" {
			if referrer, ok := event.Metadata["referrer"].(string); ok && referrer != "" {
				referrerCounts[referrer]++
			}
		}
	}

	var topReferrer string
	var maxCount int
	for referrer, count := range referrerCounts {
		if count > maxCount {
			maxCount = count
			topReferrer = referrer
		}
	}

	return topReferrer
}

func (rs *ReporterService) getTopCountry(sessions []core.UserSession) string {
	countryCounts := make(map[string]int)
	for _, session := range sessions {
		if session.Country != "" {
			countryCounts[session.Country]++
		}
	}

	var topCountry string
	var maxCount int
	for country, count := range countryCounts {
		if count > maxCount {
			maxCount = count
			topCountry = country
		}
	}

	return topCountry
}

func (rs *ReporterService) getTopDevice(sessions []core.UserSession) string {
	deviceCounts := make(map[string]int)
	for _, session := range sessions {
		if session.Device != "" {
			deviceCounts[session.Device]++
		}
	}

	var topDevice string
	var maxCount int
	for device, count := range deviceCounts {
		if count > maxCount {
			maxCount = count
			topDevice = device
		}
	}

	return topDevice
}

func (rs *ReporterService) getTopBrowser(sessions []core.UserSession) string {
	browserCounts := make(map[string]int)
	for _, session := range sessions {
		if session.Browser != "" {
			browserCounts[session.Browser]++
		}
	}

	var topBrowser string
	var maxCount int
	for browser, count := range browserCounts {
		if count > maxCount {
			maxCount = count
			topBrowser = browser
		}
	}

	return topBrowser
}

func (rs *ReporterService) generateRecommendations(summary *core.AnalyticsSummary, insights []core.AnalyticsInsight) []string {
	var recommendations []string

	if summary.AvgPageLoadTime > 3000 {
		recommendations = append(recommendations, "Optimize page load times - current average is above 3 seconds")
	}

	if summary.BounceRate > 0.7 {
		recommendations = append(recommendations, "High bounce rate detected - consider improving page content and user experience")
	}

	if summary.AvgSessionDuration < 60000 {
		recommendations = append(recommendations, "Short session durations - consider improving engagement")
	}

	for _, insight := range insights {
		if insight.Actionable && len(insight.Recommendations) > 0 {
			recommendations = append(recommendations, insight.Recommendations...)
		}
	}

	return recommendations
}

func (rs *ReporterService) generateAdditionalInsights(startTime, endTime time.Time) []core.AnalyticsInsight {
	var insights []core.AnalyticsInsight

	duration := endTime.Sub(startTime)
	if duration.Hours() > 24 {
		insight := core.AnalyticsInsight{
			ID:          generateInsightID(),
			Type:        "temporal",
			Title:       "Extended Analysis Period",
			Description: fmt.Sprintf("Analysis covers %.1f hours of data", duration.Hours()),
			Impact:      "low",
			Confidence:  1.0,
			Timestamp:   time.Now(),
			Actionable:  false,
		}
		insights = append(insights, insight)
	}

	return insights
}

func generateReportID() string {
	return fmt.Sprintf("report_%d", time.Now().UnixNano())
}
