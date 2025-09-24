package analytics

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Tracker handles real-time analytics tracking
type Tracker struct {
	sessions     map[string]*UserSession
	pageViews    []PageView
	journeys     map[string]*UserJourney
	patterns     []BehavioralPattern
	performance  []PerformanceMetrics
	insights     []AnalyticsInsight
	alerts       []AnalyticsAlert
	mu           sync.RWMutex
	subscribers  map[string]chan RealTimeAnalytics
	subscriberMu sync.RWMutex
	terminalUI   *TerminalUI
}

// NewTracker creates a new analytics tracker
func NewTracker() *Tracker {
	tracker := &Tracker{
		sessions:    make(map[string]*UserSession),
		pageViews:   make([]PageView, 0),
		journeys:    make(map[string]*UserJourney),
		patterns:    make([]BehavioralPattern, 0),
		performance: make([]PerformanceMetrics, 0),
		insights:    make([]AnalyticsInsight, 0),
		alerts:      make([]AnalyticsAlert, 0),
		subscribers: make(map[string]chan RealTimeAnalytics),
	}

	// Initialize terminal UI
	tracker.terminalUI = NewTerminalUI(tracker)

	return tracker
}

func (t *Tracker) TrackPageView(ctx *gin.Context, data PageViewData) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	session, err := t.getOrCreateSession(ctx, data.SessionID)
	if err != nil {
		return err
	}

	pageView := PageView{
		ID:                     generateID(),
		SessionID:              session.SessionID,
		URL:                    data.URL,
		Path:                   data.Path,
		Title:                  data.Title,
		Referrer:               data.Referrer,
		Timestamp:              time.Now(),
		LoadTime:               data.LoadTime,
		RenderTime:             data.RenderTime,
		FirstPaint:             data.FirstPaint,
		FirstContentfulPaint:   data.FirstContentfulPaint,
		LargestContentfulPaint: data.LargestContentfulPaint,
		CumulativeLayoutShift:  data.CumulativeLayoutShift,
		FirstInputDelay:        data.FirstInputDelay,
		TimeOnPage:             data.TimeOnPage,
		ScrollDepth:            data.ScrollDepth,
		BounceRate:             data.BounceRate,
		ExitRate:               data.ExitRate,
		Metadata:               data.Metadata,
	}

	// Add to page views
	t.pageViews = append(t.pageViews, pageView)

	// Log to terminal UI
	if t.terminalUI != nil {
		t.terminalUI.LogPageView(pageView)
	}

	// Update journey
	t.updateJourney(session.SessionID, pageView)

	// Track performance metrics
	t.trackPerformanceMetrics(pageView)

	// Generate insights
	go t.generateInsights(session.SessionID, pageView)

	// Broadcast to subscribers
	go t.broadcastRealTimeData()

	return nil
}

// TrackBehavioralPattern tracks user behavior patterns
func (t *Tracker) TrackBehavioralPattern(ctx *gin.Context, data BehavioralPatternData) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	pattern := BehavioralPattern{
		SessionID:   data.SessionID,
		UserID:      data.UserID,
		PatternType: data.PatternType,
		Element:     data.Element,
		Coordinates: data.Coordinates,
		Duration:    data.Duration,
		Intensity:   data.Intensity,
		Frequency:   data.Frequency,
		Timestamp:   time.Now(),
		Context:     data.Context,
		HeatmapData: data.HeatmapData,
	}

	t.patterns = append(t.patterns, pattern)

	// Generate behavioral insights
	go t.generateBehavioralInsights(pattern)

	return nil
}

// GetRealTimeAnalytics returns current real-time analytics
func (t *Tracker) GetRealTimeAnalytics() RealTimeAnalytics {
	t.mu.RLock()
	defer t.mu.RUnlock()

	now := time.Now()
	oneMinuteAgo := now.Add(-time.Minute)

	// Calculate active users and sessions
	activeSessions := 0
	activeUsers := make(map[string]bool)
	pageViewsLastMinute := 0

	for _, session := range t.sessions {
		if session.IsActive && session.StartTime.After(oneMinuteAgo) {
			activeSessions++
			if session.UserID != "" {
				activeUsers[session.UserID] = true
			}
		}
	}

	// Count page views in last minute
	for _, pv := range t.pageViews {
		if pv.Timestamp.After(oneMinuteAgo) {
			pageViewsLastMinute++
		}
	}

	// Calculate top pages
	topPages := t.calculateTopPages(oneMinuteAgo)
	topReferrers := t.calculateTopReferrers(oneMinuteAgo)
	topCountries := t.calculateTopCountries(oneMinuteAgo)
	topDevices := t.calculateTopDevices(oneMinuteAgo)
	topBrowsers := t.calculateTopBrowsers(oneMinuteAgo)

	// Calculate performance score
	performanceScore := t.calculatePerformanceScore(oneMinuteAgo)

	// Calculate bounce rate
	bounceRate := t.calculateBounceRate(oneMinuteAgo)

	// Calculate conversion rate
	conversionRate := t.calculateConversionRate(oneMinuteAgo)

	// Get active alerts
	activeAlerts := t.getActiveAlerts()

	return RealTimeAnalytics{
		Timestamp:          now,
		ActiveUsers:        len(activeUsers),
		ActiveSessions:     activeSessions,
		PageViewsPerMinute: pageViewsLastMinute,
		TopPages:           topPages,
		TopReferrers:       topReferrers,
		TopCountries:       topCountries,
		TopDevices:         topDevices,
		TopBrowsers:        topBrowsers,
		PerformanceScore:   performanceScore,
		ErrorRate:          0, // TODO: Calculate from error logs
		BounceRate:         bounceRate,
		ConversionRate:     conversionRate,
		Alerts:             activeAlerts,
	}
}

// GetUserJourney returns a user's complete journey
func (t *Tracker) GetUserJourney(sessionID string) (*UserJourney, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	journey, exists := t.journeys[sessionID]
	if !exists {
		return nil, fmt.Errorf("journey not found for session %s", sessionID)
	}

	return journey, nil
}

// GetFunnelAnalysis performs funnel analysis for a given funnel
func (t *Tracker) GetFunnelAnalysis(funnelID string, stages []string) (*FunnelAnalysis, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Analyze user journeys through the funnel stages
	funnelStages := make([]FunnelStage, len(stages))
	totalUsers := 0
	dropOffRates := make([]float64, len(stages)-1)

	for i, stage := range stages {
		stageUsers := t.countUsersAtStage(stage)
		if i == 0 {
			totalUsers = stageUsers
		}

		conversionRate := 0.0
		if totalUsers > 0 {
			conversionRate = float64(stageUsers) / float64(totalUsers) * 100
		}

		avgTime := t.calculateAverageTimeAtStage(stage)
		bounceRate := t.calculateBounceRateAtStage(stage)
		exitRate := t.calculateExitRateAtStage(stage)

		funnelStages[i] = FunnelStage{
			StageID:        stage,
			StageName:      stage,
			PagePath:       stage,
			Users:          stageUsers,
			ConversionRate: conversionRate,
			AverageTime:    avgTime,
			BounceRate:     bounceRate,
			ExitRate:       exitRate,
		}

		// Calculate drop-off rate
		if i > 0 {
			prevUsers := funnelStages[i-1].Users
			if prevUsers > 0 {
				dropOffRates[i-1] = float64(prevUsers-stageUsers) / float64(prevUsers) * 100
			}
		}
	}

	// Calculate overall conversion rate
	overallConversionRate := 0.0
	if totalUsers > 0 && len(funnelStages) > 0 {
		lastStageUsers := funnelStages[len(funnelStages)-1].Users
		overallConversionRate = float64(lastStageUsers) / float64(totalUsers) * 100
	}

	// Identify bottlenecks
	bottlenecks := t.identifyFunnelBottlenecks(funnelStages, dropOffRates)

	// Generate insights and recommendations
	insights := t.generateFunnelInsights(funnelStages, dropOffRates)
	recommendations := t.generateFunnelRecommendations(bottlenecks)

	return &FunnelAnalysis{
		FunnelID:        funnelID,
		FunnelName:      fmt.Sprintf("Funnel %s", funnelID),
		Stages:          funnelStages,
		TotalUsers:      totalUsers,
		ConversionRate:  overallConversionRate,
		DropOffRates:    dropOffRates,
		Bottlenecks:     bottlenecks,
		Insights:        insights,
		Recommendations: recommendations,
	}, nil
}

// SubscribeToRealTimeData subscribes to real-time analytics updates
func (t *Tracker) SubscribeToRealTimeData() chan RealTimeAnalytics {
	t.subscriberMu.Lock()
	defer t.subscriberMu.Unlock()

	subscriberID := generateID()
	ch := make(chan RealTimeAnalytics, 100)
	t.subscribers[subscriberID] = ch

	// Send initial data
	go func() {
		ch <- t.GetRealTimeAnalytics()
	}()

	return ch
}

// UnsubscribeFromRealTimeData unsubscribes from real-time analytics updates
func (t *Tracker) UnsubscribeFromRealTimeData(subscriberID string) {
	t.subscriberMu.Lock()
	defer t.subscriberMu.Unlock()

	if ch, exists := t.subscribers[subscriberID]; exists {
		close(ch)
		delete(t.subscribers, subscriberID)
	}
}

// Helper methods

func (t *Tracker) getOrCreateSession(ctx *gin.Context, sessionID string) (*UserSession, error) {
	if session, exists := t.sessions[sessionID]; exists {
		return session, nil
	}

	// Create new session
	session := &UserSession{
		SessionID: sessionID,
		StartTime: time.Now(),
		UserAgent: ctx.GetHeader("User-Agent"),
		IPAddress: ctx.ClientIP(),
		IsActive:  true,
		Metadata:  make(map[string]interface{}),
	}

	t.sessions[sessionID] = session
	return session, nil
}

func (t *Tracker) updateJourney(sessionID string, pageView PageView) {
	journey, exists := t.journeys[sessionID]
	if !exists {
		journey = &UserJourney{
			SessionID:   sessionID,
			StartTime:   pageView.Timestamp,
			PageViews:   make([]PageView, 0),
			JourneyPath: make([]string, 0),
		}
		t.journeys[sessionID] = journey
	}

	journey.PageViews = append(journey.PageViews, pageView)
	journey.JourneyPath = append(journey.JourneyPath, pageView.Path)
	journey.TotalPages = len(journey.PageViews)
	journey.TotalTime = int64(time.Since(journey.StartTime).Milliseconds())

	// Check for bounce (single page view)
	journey.BounceRate = len(journey.PageViews) == 1

	// Update conversion rate (simplified)
	journey.ConversionRate = t.calculateJourneyConversionRate(journey)
}

func (t *Tracker) trackPerformanceMetrics(pageView PageView) {
	metrics := PerformanceMetrics{
		PageID:                 pageView.ID,
		URL:                    pageView.URL,
		Timestamp:              pageView.Timestamp,
		LoadTime:               pageView.LoadTime,
		RenderTime:             pageView.RenderTime,
		FirstPaint:             pageView.FirstPaint,
		FirstContentfulPaint:   pageView.FirstContentfulPaint,
		LargestContentfulPaint: pageView.LargestContentfulPaint,
		CumulativeLayoutShift:  pageView.CumulativeLayoutShift,
		FirstInputDelay:        pageView.FirstInputDelay,
	}

	t.performance = append(t.performance, metrics)
}

func (t *Tracker) generateInsights(sessionID string, pageView PageView) {
	// Generate performance insights
	if pageView.LoadTime > 3000 { // Slow load time
		insight := AnalyticsInsight{
			ID:          generateID(),
			Type:        "performance",
			Title:       "Slow Page Load Detected",
			Description: fmt.Sprintf("Page %s took %dms to load, which is above the recommended 3 seconds", pageView.Path, pageView.LoadTime),
			Impact:      "high",
			Confidence:  0.9,
			Timestamp:   time.Now(),
			Data: map[string]interface{}{
				"pagePath":  pageView.Path,
				"loadTime":  pageView.LoadTime,
				"threshold": 3000,
			},
			Recommendations: []string{
				"Optimize images and assets",
				"Enable compression",
				"Use a CDN",
				"Minify CSS and JavaScript",
			},
			Actionable: true,
		}
		t.insights = append(t.insights, insight)
	}

	// Generate behavioral insights
	if pageView.ScrollDepth < 25 { // Low engagement
		insight := AnalyticsInsight{
			ID:          generateID(),
			Type:        "user_behavior",
			Title:       "Low Engagement Detected",
			Description: fmt.Sprintf("Users are only scrolling %f%% on page %s", pageView.ScrollDepth, pageView.Path),
			Impact:      "medium",
			Confidence:  0.7,
			Timestamp:   time.Now(),
			Data: map[string]interface{}{
				"pagePath":    pageView.Path,
				"scrollDepth": pageView.ScrollDepth,
			},
			Recommendations: []string{
				"Improve content quality and relevance",
				"Optimize page layout and structure",
				"Add engaging visual elements",
				"Improve mobile responsiveness",
			},
			Actionable: true,
		}
		t.insights = append(t.insights, insight)
	}
}

func (t *Tracker) generateBehavioralInsights(pattern BehavioralPattern) {
	// Analyze behavioral patterns and generate insights
	// This would contain complex pattern analysis logic
}

func (t *Tracker) broadcastRealTimeData() {
	t.subscriberMu.RLock()
	defer t.subscriberMu.RUnlock()

	data := t.GetRealTimeAnalytics()
	for _, ch := range t.subscribers {
		select {
		case ch <- data:
		default:
			// Channel is full, skip this update
		}
	}
}

// Additional helper methods for calculations
func (t *Tracker) calculateTopPages(since time.Time) []PageStats {
	pageStats := make(map[string]*PageStats)

	for _, pv := range t.pageViews {
		if pv.Timestamp.After(since) {
			if stats, exists := pageStats[pv.Path]; exists {
				stats.Views++
				stats.AvgTime = (stats.AvgTime + pv.TimeOnPage) / 2
			} else {
				pageStats[pv.Path] = &PageStats{
					Path:        pv.Path,
					Title:       pv.Title,
					Views:       1,
					UniqueViews: 1,
					AvgTime:     pv.TimeOnPage,
					LoadTime:    pv.LoadTime,
				}
			}
		}
	}

	// Convert to slice and sort by views
	result := make([]PageStats, 0, len(pageStats))
	for _, stats := range pageStats {
		result = append(result, *stats)
	}

	// Sort by views (simplified)
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Views < result[j].Views {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

func (t *Tracker) calculateTopReferrers(since time.Time) []ReferrerStats {
	referrerCounts := make(map[string]int)
	total := 0

	for _, pv := range t.pageViews {
		if pv.Timestamp.After(since) && pv.Referrer != "" {
			referrerCounts[pv.Referrer]++
			total++
		}
	}

	result := make([]ReferrerStats, 0, len(referrerCounts))
	for referrer, count := range referrerCounts {
		result = append(result, ReferrerStats{
			Referrer:   referrer,
			Count:      count,
			Percentage: float64(count) / float64(total) * 100,
		})
	}

	return result
}

func (t *Tracker) calculateTopCountries(since time.Time) []CountryStats {
	countryCounts := make(map[string]int)
	total := 0

	for _, session := range t.sessions {
		if session.StartTime.After(since) && session.Country != "" {
			countryCounts[session.Country]++
			total++
		}
	}

	result := make([]CountryStats, 0, len(countryCounts))
	for country, count := range countryCounts {
		result = append(result, CountryStats{
			Country:    country,
			Count:      count,
			Percentage: float64(count) / float64(total) * 100,
		})
	}

	return result
}

func (t *Tracker) calculateTopDevices(since time.Time) []DeviceStats {
	deviceCounts := make(map[string]int)
	total := 0

	for _, session := range t.sessions {
		if session.StartTime.After(since) && session.Device != "" {
			deviceCounts[session.Device]++
			total++
		}
	}

	result := make([]DeviceStats, 0, len(deviceCounts))
	for device, count := range deviceCounts {
		result = append(result, DeviceStats{
			Device:     device,
			Count:      count,
			Percentage: float64(count) / float64(total) * 100,
		})
	}

	return result
}

func (t *Tracker) calculateTopBrowsers(since time.Time) []BrowserStats {
	browserCounts := make(map[string]int)
	total := 0

	for _, session := range t.sessions {
		if session.StartTime.After(since) && session.Browser != "" {
			browserCounts[session.Browser]++
			total++
		}
	}

	result := make([]BrowserStats, 0, len(browserCounts))
	for browser, count := range browserCounts {
		result = append(result, BrowserStats{
			Browser:    browser,
			Count:      count,
			Percentage: float64(count) / float64(total) * 100,
		})
	}

	return result
}

func (t *Tracker) calculatePerformanceScore(since time.Time) float64 {
	if len(t.performance) == 0 {
		return 0.0
	}

	var totalScore float64
	count := 0

	for _, perf := range t.performance {
		if perf.Timestamp.After(since) {
			// Calculate performance score based on Core Web Vitals
			score := 100.0

			// Penalize slow load times
			if perf.LoadTime > 3000 {
				score -= 20
			} else if perf.LoadTime > 2000 {
				score -= 10
			}

			// Penalize poor LCP
			if perf.LargestContentfulPaint > 4000 {
				score -= 20
			} else if perf.LargestContentfulPaint > 2500 {
				score -= 10
			}

			// Penalize poor CLS
			if perf.CumulativeLayoutShift > 0.25 {
				score -= 20
			} else if perf.CumulativeLayoutShift > 0.1 {
				score -= 10
			}

			// Penalize poor FID
			if perf.FirstInputDelay > 300 {
				score -= 20
			} else if perf.FirstInputDelay > 100 {
				score -= 10
			}

			totalScore += score
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return totalScore / float64(count)
}

func (t *Tracker) calculateBounceRate(since time.Time) float64 {
	totalSessions := 0
	bouncedSessions := 0

	for _, journey := range t.journeys {
		if journey.StartTime.After(since) {
			totalSessions++
			if journey.BounceRate {
				bouncedSessions++
			}
		}
	}

	if totalSessions == 0 {
		return 0.0
	}

	return float64(bouncedSessions) / float64(totalSessions) * 100
}

func (t *Tracker) calculateConversionRate(since time.Time) float64 {
	totalSessions := 0
	convertedSessions := 0

	for _, journey := range t.journeys {
		if journey.StartTime.After(since) {
			totalSessions++
			if journey.GoalCompleted {
				convertedSessions++
			}
		}
	}

	if totalSessions == 0 {
		return 0.0
	}

	return float64(convertedSessions) / float64(totalSessions) * 100
}

func (t *Tracker) getActiveAlerts() []AnalyticsAlert {
	activeAlerts := make([]AnalyticsAlert, 0)
	for _, alert := range t.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}
	return activeAlerts
}

// Additional helper methods for funnel analysis
func (t *Tracker) countUsersAtStage(stage string) int {
	count := 0
	for _, journey := range t.journeys {
		for _, path := range journey.JourneyPath {
			if path == stage {
				count++
				break
			}
		}
	}
	return count
}

func (t *Tracker) calculateAverageTimeAtStage(stage string) int64 {
	var totalTime int64
	var count int

	for _, journey := range t.journeys {
		for i, path := range journey.JourneyPath {
			if path == stage {
				if i < len(journey.PageViews) {
					totalTime += journey.PageViews[i].TimeOnPage
					count++
				}
				break
			}
		}
	}

	if count == 0 {
		return 0
	}

	return totalTime / int64(count)
}

func (t *Tracker) calculateBounceRateAtStage(stage string) float64 {
	// Simplified bounce rate calculation
	return 0.0 // TODO: Implement proper bounce rate calculation
}

func (t *Tracker) calculateExitRateAtStage(stage string) float64 {
	// Simplified exit rate calculation
	return 0.0 // TODO: Implement proper exit rate calculation
}

func (t *Tracker) identifyFunnelBottlenecks(stages []FunnelStage, dropOffRates []float64) []FunnelBottleneck {
	bottlenecks := make([]FunnelBottleneck, 0)

	for i, dropOffRate := range dropOffRates {
		if dropOffRate > 50 { // High drop-off rate
			severity := "medium"
			if dropOffRate > 75 {
				severity = "high"
			}

			bottleneck := FunnelBottleneck{
				StageID:     stages[i].StageID,
				StageName:   stages[i].StageName,
				DropOffRate: dropOffRate,
				Severity:    severity,
				Impact:      dropOffRate,
				RootCause:   "High user drop-off detected",
				Recommendations: []string{
					"Improve page content and user experience",
					"Add progress indicators",
					"Simplify the process",
					"Address technical issues",
				},
			}

			bottlenecks = append(bottlenecks, bottleneck)
		}
	}

	return bottlenecks
}

func (t *Tracker) generateFunnelInsights(stages []FunnelStage, dropOffRates []float64) []string {
	insights := make([]string, 0)

	// Overall conversion rate insight
	if len(stages) > 0 {
		overallConversion := stages[len(stages)-1].ConversionRate
		if overallConversion < 10 {
			insights = append(insights, "Overall conversion rate is very low, indicating significant optimization opportunities")
		} else if overallConversion < 25 {
			insights = append(insights, "Conversion rate is below average, consider funnel optimization")
		}
	}

	// Drop-off insights
	for i, dropOffRate := range dropOffRates {
		if dropOffRate > 50 {
			insights = append(insights, fmt.Sprintf("Stage %s has a %f%% drop-off rate, indicating user friction", stages[i].StageName, dropOffRate))
		}
	}

	return insights
}

func (t *Tracker) generateFunnelRecommendations(bottlenecks []FunnelBottleneck) []string {
	recommendations := make([]string, 0)

	for _, bottleneck := range bottlenecks {
		recommendations = append(recommendations, bottleneck.Recommendations...)
	}

	// General recommendations
	recommendations = append(recommendations, "Implement A/B testing to optimize conversion paths")
	recommendations = append(recommendations, "Add user feedback collection at each stage")
	recommendations = append(recommendations, "Monitor and analyze user behavior patterns")

	return recommendations
}

func (t *Tracker) calculateJourneyConversionRate(journey *UserJourney) float64 {
	// Simplified conversion rate calculation
	// In a real implementation, this would check against defined conversion goals
	return 0.0
}

// GetSessions returns all sessions
func (t *Tracker) GetSessions() map[string]*UserSession {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.sessions
}

// GetPageViews returns all page views
func (t *Tracker) GetPageViews() []PageView {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.pageViews
}

// GetJourneys returns all journeys
func (t *Tracker) GetJourneys() map[string]*UserJourney {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.journeys
}

// GetInsights returns all insights
func (t *Tracker) GetInsights() []AnalyticsInsight {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.insights
}

// StartTerminalUI starts the terminal UI
func (t *Tracker) StartTerminalUI() {
	if t.terminalUI != nil {
		go t.terminalUI.Start()
	}
}

// StopTerminalUI stops the terminal UI
func (t *Tracker) StopTerminalUI() {
	if t.terminalUI != nil {
		t.terminalUI.Stop()
	}
}

// Data structures for API requests
type PageViewData struct {
	SessionID              string                 `json:"sessionId"`
	URL                    string                 `json:"url"`
	Path                   string                 `json:"path"`
	Title                  string                 `json:"title"`
	Referrer               string                 `json:"referrer,omitempty"`
	LoadTime               int64                  `json:"loadTime"`
	RenderTime             int64                  `json:"renderTime"`
	FirstPaint             int64                  `json:"firstPaint"`
	FirstContentfulPaint   int64                  `json:"firstContentfulPaint"`
	LargestContentfulPaint int64                  `json:"largestContentfulPaint"`
	CumulativeLayoutShift  float64                `json:"cumulativeLayoutShift"`
	FirstInputDelay        int64                  `json:"firstInputDelay"`
	TimeOnPage             int64                  `json:"timeOnPage"`
	ScrollDepth            float64                `json:"scrollDepth"`
	BounceRate             bool                   `json:"bounceRate"`
	ExitRate               bool                   `json:"exitRate"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}

type BehavioralPatternData struct {
	SessionID   string                 `json:"sessionId"`
	UserID      string                 `json:"userId,omitempty"`
	PatternType string                 `json:"patternType"`
	Element     string                 `json:"element,omitempty"`
	Coordinates map[string]float64     `json:"coordinates,omitempty"`
	Duration    int64                  `json:"duration"`
	Intensity   float64                `json:"intensity"`
	Frequency   int                    `json:"frequency"`
	Context     map[string]interface{} `json:"context,omitempty"`
	HeatmapData []HeatmapPoint         `json:"heatmapData,omitempty"`
}

// Utility functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
