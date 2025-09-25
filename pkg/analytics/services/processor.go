package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type ProcessorService struct {
	storage     core.AnalyticsStorage
	calculator  core.AnalyticsCalculator
	aggregator  core.AnalyticsAggregator
}

func NewProcessorService(
	storage core.AnalyticsStorage,
	calculator core.AnalyticsCalculator,
	aggregator core.AnalyticsAggregator,
) *ProcessorService {
	return &ProcessorService{
		storage:    storage,
		calculator: calculator,
		aggregator: aggregator,
	}
}

func (ps *ProcessorService) ProcessUserJourney(sessionID string) (*core.UserJourney, error) {
	session, err := ps.storage.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	request := core.AnalyticsRequest{
		SessionID: sessionID,
	}
	events, err := ps.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	journey := &core.UserJourney{
		SessionID: sessionID,
		UserID:    session.UserID,
		StartTime: session.StartTime,
		EndTime:   session.EndTime,
	}

	pageViews := 0
	var totalTime int64
	var journeyPath []string

	for _, event := range events {
		if event.Type == "page_view" {
			pageViews++
			journeyPath = append(journeyPath, event.Metadata["path"].(string))
		}
	}

	journey.TotalPages = pageViews
	journey.JourneyPath = journeyPath

	if session.EndTime != nil {
		journey.TotalTime = session.EndTime.Sub(session.StartTime).Milliseconds()
	} else {
		journey.TotalTime = time.Since(session.StartTime).Milliseconds()
	}

	bounceRate, err := ps.calculator.CalculateBounceRate(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate bounce rate: %w", err)
	}
	journey.BounceRate = bounceRate > 0.5

	if err := ps.storage.SaveJourney(*journey); err != nil {
		return nil, fmt.Errorf("failed to save journey: %w", err)
	}

	return journey, nil
}

func (ps *ProcessorService) ProcessFunnelAnalysis(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error) {
	funnelAnalysis, err := ps.aggregator.AggregateFunnelData(funnelID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate funnel data: %w", err)
	}

	conversionRate, err := ps.calculator.CalculateConversionRate(funnelID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate conversion rate: %w", err)
	}
	funnelAnalysis.ConversionRate = conversionRate

	dropOffRates, err := ps.calculator.CalculateFunnelDropOff(funnelID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate drop-off rates: %w", err)
	}
	funnelAnalysis.DropOffRates = dropOffRates

	funnelAnalysis.Bottlenecks = ps.identifyBottlenecks(funnelAnalysis.Stages, dropOffRates)

	funnelAnalysis.Insights = ps.generateFunnelInsights(funnelAnalysis)

	funnelAnalysis.Recommendations = ps.generateFunnelRecommendations(funnelAnalysis)

	return funnelAnalysis, nil
}

func (ps *ProcessorService) ProcessRealTimeMetrics() (*core.RealTimeMetrics, error) {
	now := time.Now()
	startTime := now.Add(-5 * time.Minute)

	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &now,
	}
	sessions, err := ps.storage.GetSessions(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	events, err := ps.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	metrics := &core.RealTimeMetrics{
		Timestamp:      now,
		ActiveUsers:    len(sessions),
		ActiveSessions: len(sessions),
	}

	pageViewCount := 0
	for _, event := range events {
		if event.Type == "page_view" {
			pageViewCount++
		}
	}
	metrics.PageViewsPerMinute = pageViewCount / 5

	topPages, err := ps.aggregateTopPages(events)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate top pages: %w", err)
	}
	metrics.TopPages = topPages

	topReferrers, err := ps.aggregateTopReferrers(events)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate top referrers: %w", err)
	}
	metrics.TopReferrers = topReferrers

	topCountries, err := ps.aggregateTopCountries(sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate top countries: %w", err)
	}
	metrics.TopCountries = topCountries

	topDevices, err := ps.aggregateTopDevices(sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate top devices: %w", err)
	}
	metrics.TopDevices = topDevices

	topBrowsers, err := ps.aggregateTopBrowsers(sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate top browsers: %w", err)
	}
	metrics.TopBrowsers = topBrowsers

	performanceEvents := ps.filterPerformanceEvents(events)
	performanceScore, err := ps.calculator.CalculatePerformanceScore(performanceEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate performance score: %w", err)
	}
	metrics.PerformanceScore = performanceScore

	metrics.BounceRate = ps.calculateBounceRate(sessions)

	alerts, err := ps.storage.GetAlerts(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}
	metrics.Alerts = alerts

	return metrics, nil
}

func (ps *ProcessorService) ProcessInsights(sessionID string) ([]core.AnalyticsInsight, error) {
	session, err := ps.storage.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	request := core.AnalyticsRequest{
		SessionID: sessionID,
	}
	events, err := ps.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	insights := ps.generateInsights(session, events)

	for _, insight := range insights {
		if err := ps.storage.SaveInsight(insight); err != nil {
			return nil, fmt.Errorf("failed to save insight: %w", err)
		}
	}

	return insights, nil
}

func (ps *ProcessorService) ProcessAlerts() ([]core.AnalyticsAlert, error) {
	metrics, err := ps.ProcessRealTimeMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to process real-time metrics: %w", err)
	}

	alerts := ps.generateAlerts(metrics)

	for _, alert := range alerts {
		if err := ps.storage.SaveAlert(alert); err != nil {
			return nil, fmt.Errorf("failed to save alert: %w", err)
		}
	}

	return alerts, nil
}

func (ps *ProcessorService) ProcessHeatmapData(pagePath string, startTime, endTime time.Time) ([]core.HeatmapPoint, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"page_path": pagePath,
			"type":      "behavioral",
		},
	}
	events, err := ps.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	heatmapPoints := ps.processHeatmapPoints(events)

	return heatmapPoints, nil
}

func (ps *ProcessorService) identifyBottlenecks(stages []core.FunnelStage, dropOffRates []float64) []core.FunnelBottleneck {
	var bottlenecks []core.FunnelBottleneck

	for i, stage := range stages {
		if i < len(dropOffRates) && dropOffRates[i] > 0.3 {
			severity := "low"
			if dropOffRates[i] > 0.6 {
				severity = "high"
			} else if dropOffRates[i] > 0.4 {
				severity = "medium"
			}

			bottleneck := core.FunnelBottleneck{
				StageID:     stage.StageID,
				StageName:   stage.StageName,
				DropOffRate: dropOffRates[i],
				Severity:    severity,
				Impact:      dropOffRates[i] * 100,
				RootCause:   ps.identifyRootCause(stage, dropOffRates[i]),
				Recommendations: ps.generateBottleneckRecommendations(stage, dropOffRates[i]),
			}
			bottlenecks = append(bottlenecks, bottleneck)
		}
	}

	return bottlenecks
}

func (ps *ProcessorService) generateFunnelInsights(funnel *core.FunnelAnalysis) []string {
	var insights []string

	if funnel.ConversionRate < 0.1 {
		insights = append(insights, "Low conversion rate detected. Consider optimizing the funnel flow.")
	}

	if len(funnel.Bottlenecks) > 0 {
		insights = append(insights, fmt.Sprintf("Found %d bottlenecks in the funnel that need attention.", len(funnel.Bottlenecks)))
	}

	return insights
}

func (ps *ProcessorService) generateFunnelRecommendations(funnel *core.FunnelAnalysis) []string {
	var recommendations []string

	if funnel.ConversionRate < 0.1 {
		recommendations = append(recommendations, "Optimize the user experience in early funnel stages")
		recommendations = append(recommendations, "A/B test different funnel flows")
	}

	for _, bottleneck := range funnel.Bottlenecks {
		recommendations = append(recommendations, bottleneck.Recommendations...)
	}

	return recommendations
}

func (ps *ProcessorService) aggregateTopPages(events []core.AnalyticsEvent) ([]core.PageStats, error) {
	pageCounts := make(map[string]int)
	pageTimes := make(map[string][]int64)

	for _, event := range events {
		if event.Type == "page_view" {
			path := event.Metadata["path"].(string)
			pageCounts[path]++
			
			if timeOnPage, ok := event.Metadata["timeOnPage"].(int64); ok {
				pageTimes[path] = append(pageTimes[path], timeOnPage)
			}
		}
	}

	var topPages []core.PageStats
	for path, count := range pageCounts {
		var avgTime int64
		if times, exists := pageTimes[path]; exists && len(times) > 0 {
			var total int64
			for _, t := range times {
				total += t
			}
			avgTime = total / int64(len(times))
		}

		topPages = append(topPages, core.PageStats{
			Path:    path,
			Views:   count,
			AvgTime: avgTime,
		})
	}

	sort.Slice(topPages, func(i, j int) bool {
		return topPages[i].Views > topPages[j].Views
	})

	if len(topPages) > 10 {
		topPages = topPages[:10]
	}

	return topPages, nil
}

func (ps *ProcessorService) aggregateTopReferrers(events []core.AnalyticsEvent) ([]core.ReferrerStats, error) {
	referrerCounts := make(map[string]int)
	totalEvents := len(events)

	for _, event := range events {
		if referrer, ok := event.Metadata["referrer"].(string); ok && referrer != "" {
			referrerCounts[referrer]++
		}
	}

	var topReferrers []core.ReferrerStats
	for referrer, count := range referrerCounts {
		percentage := float64(count) / float64(totalEvents) * 100
		topReferrers = append(topReferrers, core.ReferrerStats{
			Referrer:   referrer,
			Count:      count,
			Percentage: percentage,
		})
	}

	sort.Slice(topReferrers, func(i, j int) bool {
		return topReferrers[i].Count > topReferrers[j].Count
	})

	if len(topReferrers) > 10 {
		topReferrers = topReferrers[:10]
	}

	return topReferrers, nil
}

func (ps *ProcessorService) aggregateTopCountries(sessions []core.UserSession) ([]core.CountryStats, error) {
	countryCounts := make(map[string]int)
	totalSessions := len(sessions)

	for _, session := range sessions {
		if session.Country != "" {
			countryCounts[session.Country]++
		}
	}

	var topCountries []core.CountryStats
	for country, count := range countryCounts {
		percentage := float64(count) / float64(totalSessions) * 100
		topCountries = append(topCountries, core.CountryStats{
			Country:    country,
			Count:      count,
			Percentage: percentage,
		})
	}

	sort.Slice(topCountries, func(i, j int) bool {
		return topCountries[i].Count > topCountries[j].Count
	})

	if len(topCountries) > 10 {
		topCountries = topCountries[:10]
	}

	return topCountries, nil
}

func (ps *ProcessorService) aggregateTopDevices(sessions []core.UserSession) ([]core.DeviceStats, error) {
	deviceCounts := make(map[string]int)
	totalSessions := len(sessions)

	for _, session := range sessions {
		if session.Device != "" {
			deviceCounts[session.Device]++
		}
	}

	var topDevices []core.DeviceStats
	for device, count := range deviceCounts {
		percentage := float64(count) / float64(totalSessions) * 100
		topDevices = append(topDevices, core.DeviceStats{
			Device:     device,
			Count:      count,
			Percentage: percentage,
		})
	}

	sort.Slice(topDevices, func(i, j int) bool {
		return topDevices[i].Count > topDevices[j].Count
	})

	if len(topDevices) > 10 {
		topDevices = topDevices[:10]
	}

	return topDevices, nil
}

func (ps *ProcessorService) aggregateTopBrowsers(sessions []core.UserSession) ([]core.BrowserStats, error) {
	browserCounts := make(map[string]int)
	totalSessions := len(sessions)

	for _, session := range sessions {
		if session.Browser != "" {
			browserCounts[session.Browser]++
		}
	}

	var topBrowsers []core.BrowserStats
	for browser, count := range browserCounts {
		percentage := float64(count) / float64(totalSessions) * 100
		topBrowsers = append(topBrowsers, core.BrowserStats{
			Browser:    browser,
			Count:      count,
			Percentage: percentage,
		})
	}

	sort.Slice(topBrowsers, func(i, j int) bool {
		return topBrowsers[i].Count > topBrowsers[j].Count
	})

	if len(topBrowsers) > 10 {
		topBrowsers = topBrowsers[:10]
	}

	return topBrowsers, nil
}

func (ps *ProcessorService) filterPerformanceEvents(events []core.AnalyticsEvent) []core.PerformanceEvent {
	var performanceEvents []core.PerformanceEvent

	for _, event := range events {
		if event.Type == "performance" {
			perfEvent := core.PerformanceEvent{
				AnalyticsEvent: event,
			}
			performanceEvents = append(performanceEvents, perfEvent)
		}
	}

	return performanceEvents
}

func (ps *ProcessorService) calculateBounceRate(sessions []core.UserSession) float64 {
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

func (ps *ProcessorService) generateInsights(session *core.UserSession, events []core.AnalyticsEvent) []core.AnalyticsInsight {
	var insights []core.AnalyticsInsight

	if len(events) > 10 {
		insight := core.AnalyticsInsight{
			ID:          generateInsightID(),
			Type:        "user_behavior",
			Title:       "High Engagement Session",
			Description: "User showed high engagement with multiple page views",
			Impact:      "medium",
			Confidence:  0.8,
			Timestamp:   time.Now(),
			Actionable:  true,
		}
		insights = append(insights, insight)
	}

	return insights
}

func (ps *ProcessorService) generateAlerts(metrics *core.RealTimeMetrics) []core.AnalyticsAlert {
	var alerts []core.AnalyticsAlert

	if metrics.PerformanceScore < 0.5 {
		alert := core.AnalyticsAlert{
			ID:        generateAlertID(),
			Type:      "performance",
			Severity:  "high",
			Title:     "Low Performance Score",
			Message:   fmt.Sprintf("Performance score is %.2f, below threshold", metrics.PerformanceScore),
			Timestamp: time.Now(),
			Resolved:  false,
		}
		alerts = append(alerts, alert)
	}

	if metrics.BounceRate > 0.7 {
		alert := core.AnalyticsAlert{
			ID:        generateAlertID(),
			Type:      "conversion",
			Severity:  "medium",
			Title:     "High Bounce Rate",
			Message:   fmt.Sprintf("Bounce rate is %.2f%%, above threshold", metrics.BounceRate*100),
			Timestamp: time.Now(),
			Resolved:  false,
		}
		alerts = append(alerts, alert)
	}

	return alerts
}

func (ps *ProcessorService) processHeatmapPoints(events []core.AnalyticsEvent) []core.HeatmapPoint {
	pointMap := make(map[string]core.HeatmapPoint)

	for _, event := range events {
		if coordinates, ok := event.Metadata["coordinates"].(map[string]float64); ok {
			x := coordinates["x"]
			y := coordinates["y"]
			key := fmt.Sprintf("%.2f,%.2f", x, y)

			if point, exists := pointMap[key]; exists {
				point.Count++
				point.Intensity = float64(point.Count) / 10.0
				pointMap[key] = point
			} else {
				pointMap[key] = core.HeatmapPoint{
					X:         x,
					Y:         y,
					Intensity: 1.0,
					Count:     1,
				}
			}
		}
	}

	var points []core.HeatmapPoint
	for _, point := range pointMap {
		points = append(points, point)
	}

	return points
}

func (ps *ProcessorService) identifyRootCause(stage core.FunnelStage, dropOffRate float64) string {
	if dropOffRate > 0.6 {
		return "Critical drop-off point - likely UX issue"
	} else if dropOffRate > 0.4 {
		return "Significant drop-off - consider A/B testing"
	}
	return "Minor drop-off - monitor for trends"
}

func (ps *ProcessorService) generateBottleneckRecommendations(stage core.FunnelStage, dropOffRate float64) []string {
	var recommendations []string

	if dropOffRate > 0.6 {
		recommendations = append(recommendations, "Immediate UX review required")
		recommendations = append(recommendations, "Consider simplifying the flow")
	} else if dropOffRate > 0.4 {
		recommendations = append(recommendations, "A/B test alternative approaches")
		recommendations = append(recommendations, "Gather user feedback")
	}

	return recommendations
}

func generateInsightID() string {
	return fmt.Sprintf("insight_%d", time.Now().UnixNano())
}

func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}
