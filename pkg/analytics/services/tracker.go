package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type TrackerService struct {
	storage   core.AnalyticsStorage
	validator core.AnalyticsValidator
	mu        sync.RWMutex
}

func NewTrackerService(storage core.AnalyticsStorage, validator core.AnalyticsValidator) *TrackerService {
	return &TrackerService{
		storage:   storage,
		validator: validator,
	}
}

func (ts *TrackerService) TrackPageView(event core.PageViewEvent) error {
	if err := ts.validator.ValidateEvent(event.AnalyticsEvent); err != nil {
		return fmt.Errorf("invalid page view event: %w", err)
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.ID == "" {
		event.ID = generateEventID()
	}
	if err := ts.storage.SaveEvent(event.AnalyticsEvent); err != nil {
		return fmt.Errorf("failed to save page view event: %w", err)
	}
	return nil
}

func (ts *TrackerService) TrackBehavioralPattern(event core.BehavioralEvent) error {
	if err := ts.validator.ValidateEvent(event.AnalyticsEvent); err != nil {
		return fmt.Errorf("invalid behavioral event: %w", err)
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.ID == "" {
		event.ID = generateEventID()
	}
	if err := ts.storage.SaveEvent(event.AnalyticsEvent); err != nil {
		return fmt.Errorf("failed to save behavioral event: %w", err)
	}
	return nil
}

func (ts *TrackerService) TrackPerformance(event core.PerformanceEvent) error {
	if err := ts.validator.ValidateEvent(event.AnalyticsEvent); err != nil {
		return fmt.Errorf("invalid performance event: %w", err)
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.ID == "" {
		event.ID = generateEventID()
	}
	if err := ts.storage.SaveEvent(event.AnalyticsEvent); err != nil {
		return fmt.Errorf("failed to save performance event: %w", err)
	}
	return nil
}

func (ts *TrackerService) TrackCustomEvent(event core.AnalyticsEvent) error {
	if err := ts.validator.ValidateEvent(event); err != nil {
		return fmt.Errorf("invalid custom event: %w", err)
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.ID == "" {
		event.ID = generateEventID()
	}
	if err := ts.storage.SaveEvent(event); err != nil {
		return fmt.Errorf("failed to save custom event: %w", err)
	}
	return nil
}

func (ts *TrackerService) GetSession(sessionID string) (*core.UserSession, error) {
	return ts.storage.GetSession(sessionID)
}

func (ts *TrackerService) CreateSession(session core.UserSession) error {
	if err := ts.validator.ValidateSession(session); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}
	if session.StartTime.IsZero() {
		session.StartTime = time.Now()
	}
	session.IsActive = true
	if err := ts.storage.SaveSession(session); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (ts *TrackerService) UpdateSession(session core.UserSession) error {
	if err := ts.validator.ValidateSession(session); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}
	if err := ts.storage.UpdateSession(session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (ts *TrackerService) EndSession(sessionID string) error {
	session, err := ts.storage.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	now := time.Now()
	session.EndTime = &now
	session.IsActive = false
	if err := ts.storage.UpdateSession(*session); err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}
	return nil
}

func generateEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}
