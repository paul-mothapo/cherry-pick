package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type MemoryStorage struct {
	events     map[string]core.AnalyticsEvent
	sessions   map[string]core.UserSession
	journeys   map[string]core.UserJourney
	insights   map[string]core.AnalyticsInsight
	alerts     map[string]core.AnalyticsAlert
	reports    map[string]core.AnalyticsReport
	mu         sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		events:   make(map[string]core.AnalyticsEvent),
		sessions: make(map[string]core.UserSession),
		journeys: make(map[string]core.UserJourney),
		insights: make(map[string]core.AnalyticsInsight),
		alerts:   make(map[string]core.AnalyticsAlert),
		reports:  make(map[string]core.AnalyticsReport),
	}
}

func (ms *MemoryStorage) SaveEvent(event core.AnalyticsEvent) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.events[event.ID] = event
	return nil
}

func (ms *MemoryStorage) GetEvents(request core.AnalyticsRequest) ([]core.AnalyticsEvent, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var events []core.AnalyticsEvent
	for _, event := range ms.events {
		if ms.matchesEventFilters(event, request) {
			events = append(events, event)
		}
	}
	if request.Limit > 0 {
		start := request.Offset
		end := start + request.Limit
		if end > len(events) {
			end = len(events)
		}
		if start < len(events) {
			events = events[start:end]
		} else {
			events = []core.AnalyticsEvent{}
		}
	}
	return events, nil
}

func (ms *MemoryStorage) DeleteEvent(eventID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.events[eventID]; !exists {
		return fmt.Errorf("event with ID %s not found", eventID)
	}
	delete(ms.events, eventID)
	return nil
}

func (ms *MemoryStorage) SaveSession(session core.UserSession) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.sessions[session.SessionID] = session
	return nil
}

func (ms *MemoryStorage) GetSession(sessionID string) (*core.UserSession, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	session, exists := ms.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}
	return &session, nil
}

func (ms *MemoryStorage) GetSessions(request core.AnalyticsRequest) ([]core.UserSession, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var sessions []core.UserSession
	for _, session := range ms.sessions {
		if ms.matchesSessionFilters(session, request) {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (ms *MemoryStorage) UpdateSession(session core.UserSession) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.sessions[session.SessionID]; !exists {
		return fmt.Errorf("session with ID %s not found", session.SessionID)
	}
	ms.sessions[session.SessionID] = session
	return nil
}

func (ms *MemoryStorage) DeleteSession(sessionID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.sessions[sessionID]; !exists {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}
	delete(ms.sessions, sessionID)
	return nil
}

func (ms *MemoryStorage) SaveJourney(journey core.UserJourney) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.journeys[journey.SessionID] = journey
	return nil
}

func (ms *MemoryStorage) GetJourney(sessionID string) (*core.UserJourney, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	journey, exists := ms.journeys[sessionID]
	if !exists {
		return nil, fmt.Errorf("journey with session ID %s not found", sessionID)
	}
	return &journey, nil
}

func (ms *MemoryStorage) GetJourneys(request core.AnalyticsRequest) ([]core.UserJourney, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var journeys []core.UserJourney
	for _, journey := range ms.journeys {
		if ms.matchesJourneyFilters(journey, request) {
			journeys = append(journeys, journey)
		}
	}
	return journeys, nil
}

func (ms *MemoryStorage) UpdateJourney(journey core.UserJourney) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.journeys[journey.SessionID]; !exists {
		return fmt.Errorf("journey with session ID %s not found", journey.SessionID)
	}
	ms.journeys[journey.SessionID] = journey
	return nil
}

func (ms *MemoryStorage) DeleteJourney(sessionID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.journeys[sessionID]; !exists {
		return fmt.Errorf("journey with session ID %s not found", sessionID)
	}
	delete(ms.journeys, sessionID)
	return nil
}

func (ms *MemoryStorage) SaveInsight(insight core.AnalyticsInsight) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.insights[insight.ID] = insight
	return nil
}

func (ms *MemoryStorage) GetInsight(insightID string) (*core.AnalyticsInsight, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	insight, exists := ms.insights[insightID]
	if !exists {
		return nil, fmt.Errorf("insight with ID %s not found", insightID)
	}
	return &insight, nil
}

func (ms *MemoryStorage) GetInsights(request core.AnalyticsRequest) ([]core.AnalyticsInsight, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var insights []core.AnalyticsInsight
	for _, insight := range ms.insights {
		if ms.matchesInsightFilters(insight, request) {
			insights = append(insights, insight)
		}
	}
	return insights, nil
}

func (ms *MemoryStorage) UpdateInsight(insight core.AnalyticsInsight) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.insights[insight.ID]; !exists {
		return fmt.Errorf("insight with ID %s not found", insight.ID)
	}
	ms.insights[insight.ID] = insight
	return nil
}

func (ms *MemoryStorage) DeleteInsight(insightID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.insights[insightID]; !exists {
		return fmt.Errorf("insight with ID %s not found", insightID)
	}
	delete(ms.insights, insightID)
	return nil
}

func (ms *MemoryStorage) SaveAlert(alert core.AnalyticsAlert) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.alerts[alert.ID] = alert
	return nil
}

func (ms *MemoryStorage) GetAlert(alertID string) (*core.AnalyticsAlert, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	alert, exists := ms.alerts[alertID]
	if !exists {
		return nil, fmt.Errorf("alert with ID %s not found", alertID)
	}
	return &alert, nil
}

func (ms *MemoryStorage) GetAlerts(request core.AnalyticsRequest) ([]core.AnalyticsAlert, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var alerts []core.AnalyticsAlert
	for _, alert := range ms.alerts {
		if ms.matchesAlertFilters(alert, request) {
			alerts = append(alerts, alert)
		}
	}
	return alerts, nil
}

func (ms *MemoryStorage) UpdateAlert(alert core.AnalyticsAlert) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.alerts[alert.ID]; !exists {
		return fmt.Errorf("alert with ID %s not found", alert.ID)
	}
	ms.alerts[alert.ID] = alert
	return nil
}

func (ms *MemoryStorage) DeleteAlert(alertID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.alerts[alertID]; !exists {
		return fmt.Errorf("alert with ID %s not found", alertID)
	}
	delete(ms.alerts, alertID)
	return nil
}

func (ms *MemoryStorage) SaveReport(report core.AnalyticsReport) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.reports[report.ID] = report
	return nil
}

func (ms *MemoryStorage) GetReport(reportID string) (*core.AnalyticsReport, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	report, exists := ms.reports[reportID]
	if !exists {
		return nil, fmt.Errorf("report with ID %s not found", reportID)
	}
	return &report, nil
}

func (ms *MemoryStorage) GetReports(request core.AnalyticsRequest) ([]core.AnalyticsReport, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var reports []core.AnalyticsReport
	for _, report := range ms.reports {
		if ms.matchesReportFilters(report, request) {
			reports = append(reports, report)
		}
	}
	return reports, nil
}

func (ms *MemoryStorage) DeleteReport(reportID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.reports[reportID]; !exists {
		return fmt.Errorf("report with ID %s not found", reportID)
	}
	delete(ms.reports, reportID)
	return nil
}

func (ms *MemoryStorage) CleanupOldData(olderThan time.Time) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for id, event := range ms.events {
		if event.Timestamp.Before(olderThan) {
			delete(ms.events, id)
		}
	}
	for id, session := range ms.sessions {
		if session.StartTime.Before(olderThan) {
			delete(ms.sessions, id)
		}
	}
	for id, insight := range ms.insights {
		if insight.Timestamp.Before(olderThan) {
			delete(ms.insights, id)
		}
	}
	for id, alert := range ms.alerts {
		if alert.Timestamp.Before(olderThan) {
			delete(ms.alerts, id)
		}
	}
	return nil
}

func (ms *MemoryStorage) GetStats() (map[string]interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	stats := map[string]interface{}{
		"total_events":   len(ms.events),
		"total_sessions": len(ms.sessions),
		"total_journeys": len(ms.journeys),
		"total_insights": len(ms.insights),
		"total_alerts":   len(ms.alerts),
		"total_reports":  len(ms.reports),
	}
	return stats, nil
}

func (ms *MemoryStorage) matchesEventFilters(event core.AnalyticsEvent, request core.AnalyticsRequest) bool {
	if request.StartTime != nil && event.Timestamp.Before(*request.StartTime) {
		return false
	}
	if request.EndTime != nil && event.Timestamp.After(*request.EndTime) {
		return false
	}
	if request.SessionID != "" && event.SessionID != request.SessionID {
		return false
	}
	if request.UserID != "" && event.UserID != request.UserID {
		return false
	}
	if request.Filters != nil {
		for key, value := range request.Filters {
			if eventValue, exists := event.Metadata[key]; !exists || eventValue != value {
				return false
			}
		}
	}
	return true
}

func (ms *MemoryStorage) matchesSessionFilters(session core.UserSession, request core.AnalyticsRequest) bool {
	if request.StartTime != nil && session.StartTime.Before(*request.StartTime) {
		return false
	}
	if request.EndTime != nil && session.StartTime.After(*request.EndTime) {
		return false
	}
	if request.UserID != "" && session.UserID != request.UserID {
		return false
	}
	return true
}

func (ms *MemoryStorage) matchesJourneyFilters(journey core.UserJourney, request core.AnalyticsRequest) bool {
	if request.StartTime != nil && journey.StartTime.Before(*request.StartTime) {
		return false
	}
	if request.EndTime != nil && journey.StartTime.After(*request.EndTime) {
		return false
	}
	if request.UserID != "" && journey.UserID != request.UserID {
		return false
	}
	return true
}

func (ms *MemoryStorage) matchesInsightFilters(insight core.AnalyticsInsight, request core.AnalyticsRequest) bool {
	if request.StartTime != nil && insight.Timestamp.Before(*request.StartTime) {
		return false
	}
	if request.EndTime != nil && insight.Timestamp.After(*request.EndTime) {
		return false
	}
	return true
}

func (ms *MemoryStorage) matchesAlertFilters(alert core.AnalyticsAlert, request core.AnalyticsRequest) bool {
	if request.StartTime != nil && alert.Timestamp.Before(*request.StartTime) {
		return false
	}
	if request.EndTime != nil && alert.Timestamp.After(*request.EndTime) {
		return false
	}
	return true
}

func (ms *MemoryStorage) matchesReportFilters(report core.AnalyticsReport, request core.AnalyticsRequest) bool {
	if request.StartTime != nil && report.GeneratedAt.Before(*request.StartTime) {
		return false
	}
	if request.EndTime != nil && report.GeneratedAt.After(*request.EndTime) {
		return false
	}
	return true
}
