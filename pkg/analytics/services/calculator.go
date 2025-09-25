package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type CalculatorService struct {
	storage core.AnalyticsStorage
}

func NewCalculatorService(storage core.AnalyticsStorage) *CalculatorService {
	return &CalculatorService{
		storage: storage,
	}
}

func (cs *CalculatorService) CalculateBounceRate(sessionID string) (float64, error) {
	session, err := cs.storage.GetSession(sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	request := core.AnalyticsRequest{
		SessionID: sessionID,
	}
	events, err := cs.storage.GetEvents(request)
	if err != nil {
		return 0, fmt.Errorf("failed to get events: %w", err)
	}

	pageViewCount := 0
	for _, event := range events {
		if event.Type == "page_view" {
			pageViewCount++
		}
	}

	var sessionDuration time.Duration
	if session.EndTime != nil {
		sessionDuration = session.EndTime.Sub(session.StartTime)
	} else {
		sessionDuration = time.Since(session.StartTime)
	}

	if pageViewCount <= 1 {
		if sessionDuration < 30*time.Second {
			hasInteraction := false
			for _, event := range events {
				if event.Type == "behavioral" {
					hasInteraction = true
					break
				}
			}

			if !hasInteraction {
				return 1.0, nil
			}
		}
	}

	return 0.0, nil
}

func (cs *CalculatorService) CalculateConversionRate(funnelID string, startTime, endTime time.Time) (float64, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"funnel_id": funnelID,
		},
	}
	events, err := cs.storage.GetEvents(request)
	if err != nil {
		return 0, fmt.Errorf("failed to get events: %w", err)
	}

	sessionEvents := make(map[string][]core.AnalyticsEvent)
	for _, event := range events {
		sessionEvents[event.SessionID] = append(sessionEvents[event.SessionID], event)
	}

	completedSessions := 0
	totalSessions := len(sessionEvents)

	for _, sessionEventList := range sessionEvents {
		if cs.hasCompletedFunnel(sessionEventList, funnelID) {
			completedSessions++
		}
	}

	if totalSessions == 0 {
		return 0, nil
	}

	return float64(completedSessions) / float64(totalSessions), nil
}

func (cs *CalculatorService) CalculatePerformanceScore(events []core.PerformanceEvent) (float64, error) {
	if len(events) == 0 {
		return 0, nil
	}

	var totalScore float64
	var validEvents int

	for _, event := range events {
		score := cs.calculateEventPerformanceScore(event)
		if score > 0 {
			totalScore += score
			validEvents++
		}
	}

	if validEvents == 0 {
		return 0, nil
	}

	return totalScore / float64(validEvents), nil
}

func (cs *CalculatorService) CalculateUserEngagement(sessionID string) (float64, error) {
	session, err := cs.storage.GetSession(sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	request := core.AnalyticsRequest{
		SessionID: sessionID,
	}
	events, err := cs.storage.GetEvents(request)
	if err != nil {
		return 0, fmt.Errorf("failed to get events: %w", err)
	}

	var engagementScore float64

	sessionDuration := time.Since(session.StartTime)
	if session.EndTime != nil {
		sessionDuration = session.EndTime.Sub(session.StartTime)
	}
	durationScore := math.Min(float64(sessionDuration.Minutes())/30.0, 1.0)
	engagementScore += durationScore * 0.3

	pageViewCount := 0
	for _, event := range events {
		if event.Type == "page_view" {
			pageViewCount++
		}
	}
	pageViewScore := math.Min(float64(pageViewCount)/10.0, 1.0)
	engagementScore += pageViewScore * 0.3

	interactionCount := 0
	for _, event := range events {
		if event.Type == "behavioral" {
			interactionCount++
		}
	}
	interactionScore := math.Min(float64(interactionCount)/20.0, 1.0)
	engagementScore += interactionScore * 0.4

	return math.Min(engagementScore, 1.0), nil
}

func (cs *CalculatorService) CalculateFunnelDropOff(funnelID string, startTime, endTime time.Time) ([]float64, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"funnel_id": funnelID,
		},
	}
	events, err := cs.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	sessionEvents := make(map[string][]core.AnalyticsEvent)
	for _, event := range events {
		sessionEvents[event.SessionID] = append(sessionEvents[event.SessionID], event)
	}

	stages := []string{"landing", "product", "cart", "checkout", "purchase"}
	stageCounts := make([]int, len(stages))

	for _, sessionEventList := range sessionEvents {
		for i, stage := range stages {
			if cs.hasReachedStage(sessionEventList, stage) {
				stageCounts[i]++
			}
		}
	}

	dropOffRates := make([]float64, len(stages)-1)
	for i := 0; i < len(stages)-1; i++ {
		if stageCounts[i] == 0 {
			dropOffRates[i] = 0
		} else {
			dropOffRates[i] = float64(stageCounts[i]-stageCounts[i+1]) / float64(stageCounts[i])
		}
	}

	return dropOffRates, nil
}

func (cs *CalculatorService) hasCompletedFunnel(events []core.AnalyticsEvent, funnelID string) bool {
	for _, event := range events {
		if event.Type == "custom" {
			if event.Metadata["funnel_completed"] == true {
				return true
			}
		}
	}

	completionPages := []string{"/thank-you", "/success", "/confirmation"}
	for _, event := range events {
		if event.Type == "page_view" {
			if path, ok := event.Metadata["path"].(string); ok {
				for _, completionPage := range completionPages {
					if path == completionPage {
						return true
					}
				}
			}
		}
	}

	return false
}

func (cs *CalculatorService) hasReachedStage(events []core.AnalyticsEvent, stage string) bool {
	for _, event := range events {
		if event.Type == "page_view" {
			if path, ok := event.Metadata["path"].(string); ok {
				if cs.isStagePage(path, stage) {
					return true
				}
			}
		}

		if event.Type == "custom" {
			if event.Metadata["stage"] == stage {
				return true
			}
		}
	}

	return false
}

func (cs *CalculatorService) isStagePage(path, stage string) bool {
	stagePages := map[string][]string{
		"landing":  {"/", "/home", "/landing"},
		"product":  {"/product", "/products", "/item"},
		"cart":     {"/cart", "/basket"},
		"checkout": {"/checkout", "/payment"},
		"purchase": {"/thank-you", "/success", "/confirmation"},
	}

	if pages, exists := stagePages[stage]; exists {
		for _, page := range pages {
			if path == page {
				return true
			}
		}
	}

	return false
}

func (cs *CalculatorService) calculateEventPerformanceScore(event core.PerformanceEvent) float64 {
	var score float64

	if event.LoadTime > 0 {
		loadTimeScore := math.Max(0, 1.0-float64(event.LoadTime)/5000.0)
		score += loadTimeScore * 0.3
	}

	if event.FirstContentfulPaint > 0 {
		fcpScore := math.Max(0, 1.0-float64(event.FirstContentfulPaint)/3000.0)
		score += fcpScore * 0.3
	}

	if event.LargestContentfulPaint > 0 {
		lcpScore := math.Max(0, 1.0-float64(event.LargestContentfulPaint)/4000.0)
		score += lcpScore * 0.2
	}

	if event.CumulativeLayoutShift >= 0 {
		clsScore := math.Max(0, 1.0-event.CumulativeLayoutShift/0.25)
		score += clsScore * 0.1
	}

	if event.FirstInputDelay > 0 {
		fidScore := math.Max(0, 1.0-float64(event.FirstInputDelay)/300.0)
		score += fidScore * 0.1
	}

	return math.Min(score, 1.0)
}

func (cs *CalculatorService) CalculateAverageSessionDuration(sessions []core.UserSession) time.Duration {
	if len(sessions) == 0 {
		return 0
	}

	var totalDuration time.Duration
	var validSessions int

	for _, session := range sessions {
		if session.EndTime != nil {
			duration := session.EndTime.Sub(session.StartTime)
			totalDuration += duration
			validSessions++
		}
	}

	if validSessions == 0 {
		return 0
	}

	return totalDuration / time.Duration(validSessions)
}

func (cs *CalculatorService) CalculatePageViewDistribution(events []core.AnalyticsEvent) map[string]int {
	distribution := make(map[string]int)

	for _, event := range events {
		if event.Type == "page_view" {
			if path, ok := event.Metadata["path"].(string); ok {
				distribution[path]++
			}
		}
	}

	return distribution
}

func (cs *CalculatorService) CalculateDeviceDistribution(sessions []core.UserSession) map[string]int {
	distribution := make(map[string]int)

	for _, session := range sessions {
		device := session.Device
		if device == "" {
			device = "unknown"
		}
		distribution[device]++
	}

	return distribution
}

func (cs *CalculatorService) CalculateBrowserDistribution(sessions []core.UserSession) map[string]int {
	distribution := make(map[string]int)

	for _, session := range sessions {
		browser := session.Browser
		if browser == "" {
			browser = "unknown"
		}
		distribution[browser]++
	}

	return distribution
}

func (cs *CalculatorService) CalculateCountryDistribution(sessions []core.UserSession) map[string]int {
	distribution := make(map[string]int)

	for _, session := range sessions {
		country := session.Country
		if country == "" {
			country = "unknown"
		}
		distribution[country]++
	}

	return distribution
}

func (cs *CalculatorService) CalculateTopPages(events []core.AnalyticsEvent, limit int) []core.PageStats {
	pageCounts := make(map[string]int)
	pageTimes := make(map[string][]int64)

	for _, event := range events {
		if event.Type == "page_view" {
			if path, ok := event.Metadata["path"].(string); ok {
				pageCounts[path]++

				if timeOnPage, ok := event.Metadata["timeOnPage"].(int64); ok {
					pageTimes[path] = append(pageTimes[path], timeOnPage)
				}
			}
		}
	}

	var pages []core.PageStats
	for path, count := range pageCounts {
		var avgTime int64
		if times, exists := pageTimes[path]; exists && len(times) > 0 {
			var total int64
			for _, t := range times {
				total += t
			}
			avgTime = total / int64(len(times))
		}

		pages = append(pages, core.PageStats{
			Path:    path,
			Views:   count,
			AvgTime: avgTime,
		})
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Views > pages[j].Views
	})

	if limit > 0 && len(pages) > limit {
		pages = pages[:limit]
	}

	return pages
}
