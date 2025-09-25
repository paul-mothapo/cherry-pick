package services

import (
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type AggregatorService struct {
	storage core.AnalyticsStorage
}

func NewAggregatorService(storage core.AnalyticsStorage) *AggregatorService {
	return &AggregatorService{
		storage: storage,
	}
}

func (as *AggregatorService) AggregatePageViews(startTime, endTime time.Time) ([]core.PageStats, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"type": "page_view",
		},
	}
	events, err := as.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	pageStats := make(map[string]*core.PageStats)

	for _, event := range events {
		if path, ok := event.Metadata["path"].(string); ok {
			if stats, exists := pageStats[path]; exists {
				stats.Views++
				if timeOnPage, ok := event.Metadata["timeOnPage"].(int64); ok {
					stats.AvgTime = (stats.AvgTime + timeOnPage) / 2
				}
			} else {
				stats := &core.PageStats{
					Path:  path,
					Views: 1,
				}
				if timeOnPage, ok := event.Metadata["timeOnPage"].(int64); ok {
					stats.AvgTime = timeOnPage
				}
				pageStats[path] = stats
			}
		}
	}

	var result []core.PageStats
	for _, stats := range pageStats {
		result = append(result, *stats)
	}

	return result, nil
}

func (as *AggregatorService) AggregateUserSessions(startTime, endTime time.Time) ([]core.UserSession, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	return as.storage.GetSessions(request)
}

func (as *AggregatorService) AggregatePerformanceMetrics(startTime, endTime time.Time) ([]core.PerformanceEvent, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"type": "performance",
		},
	}
	events, err := as.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
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

func (as *AggregatorService) AggregateBehavioralPatterns(startTime, endTime time.Time) ([]core.BehavioralEvent, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"type": "behavioral",
		},
	}
	events, err := as.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
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

func (as *AggregatorService) AggregateFunnelData(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error) {
	request := core.AnalyticsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
		Filters: map[string]interface{}{
			"funnel_id": funnelID,
		},
	}
	events, err := as.storage.GetEvents(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	funnelAnalysis := &core.FunnelAnalysis{
		FunnelID:   funnelID,
		FunnelName: fmt.Sprintf("Funnel %s", funnelID),
		Stages:     []core.FunnelStage{},
	}

	sessionEvents := make(map[string][]core.AnalyticsEvent)
	for _, event := range events {
		sessionEvents[event.SessionID] = append(sessionEvents[event.SessionID], event)
	}

	funnelAnalysis.TotalUsers = len(sessionEvents)

	stages := []string{"landing", "product", "cart", "checkout", "purchase"}

	for i, stage := range stages {
		stageUsers := 0
		for _, sessionEventList := range sessionEvents {
			if as.hasReachedStage(sessionEventList, stage) {
				stageUsers++
			}
		}

		funnelStage := core.FunnelStage{
			StageID:        fmt.Sprintf("stage_%d", i),
			StageName:      stage,
			PagePath:       fmt.Sprintf("/%s", stage),
			Users:          stageUsers,
			ConversionRate: float64(stageUsers) / float64(funnelAnalysis.TotalUsers),
			AverageTime:    0,
			BounceRate:     0,
			ExitRate:       0,
		}

		funnelAnalysis.Stages = append(funnelAnalysis.Stages, funnelStage)
	}

	return funnelAnalysis, nil
}

func (as *AggregatorService) hasReachedStage(events []core.AnalyticsEvent, stage string) bool {
	for _, event := range events {
		if event.Type == "page_view" {
			if path, ok := event.Metadata["path"].(string); ok {
				if as.isStagePage(path, stage) {
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

func (as *AggregatorService) isStagePage(path, stage string) bool {
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
