package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type MemoryStorage struct {
	tests      map[string]*core.LoadTestSummary
	alerts     map[string]*core.Alert
	triggers   map[string]*core.AlertTrigger
	templates  map[string]*core.AlertTemplate
	history    map[string]*core.LoadTestHistory
	mu         sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tests:     make(map[string]*core.LoadTestSummary),
		alerts:    make(map[string]*core.Alert),
		triggers:  make(map[string]*core.AlertTrigger),
		templates: make(map[string]*core.AlertTemplate),
		history:   make(map[string]*core.LoadTestHistory),
	}
}

func (ms *MemoryStorage) SaveTest(test *core.LoadTestSummary) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.tests[test.TestID] = test
	return nil
}

func (ms *MemoryStorage) GetTest(testID string) (*core.LoadTestSummary, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	test, exists := ms.tests[testID]
	if !exists {
		return nil, NewStorageError("GetTest", "test", fmt.Errorf("test with ID %s not found", testID))
	}

	return test, nil
}

func (ms *MemoryStorage) GetTestsByUser(userID string) ([]*core.LoadTestSummary, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var tests []*core.LoadTestSummary
	for _, test := range ms.tests {
		tests = append(tests, test)
	}

	return tests, nil
}

func (ms *MemoryStorage) GetAllTests() ([]*core.LoadTestSummary, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var tests []*core.LoadTestSummary
	for _, test := range ms.tests {
		tests = append(tests, test)
	}

	return tests, nil
}

func (ms *MemoryStorage) DeleteTest(testID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.tests[testID]; !exists {
		return NewStorageError("DeleteTest", "test", fmt.Errorf("test with ID %s not found", testID))
	}

	delete(ms.tests, testID)
	return nil
}

func (ms *MemoryStorage) UpdateTest(testID string, test *core.LoadTestSummary) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.tests[testID]; !exists {
		return NewStorageError("UpdateTest", "test", fmt.Errorf("test with ID %s not found", testID))
	}

	ms.tests[testID] = test
	return nil
}

func (ms *MemoryStorage) SaveAlert(alert *core.Alert) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.alerts[alert.ID] = alert
	return nil
}

func (ms *MemoryStorage) GetAlert(alertID string) (*core.Alert, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	alert, exists := ms.alerts[alertID]
	if !exists {
		return nil, NewStorageError("GetAlert", "alert", fmt.Errorf("alert with ID %s not found", alertID))
	}

	return alert, nil
}

func (ms *MemoryStorage) GetAlertsByTest(testID string) ([]*core.Alert, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var alerts []*core.Alert
	for _, alert := range ms.alerts {
		if alert.TestID == testID {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

func (ms *MemoryStorage) GetAllAlerts() ([]*core.Alert, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var alerts []*core.Alert
	for _, alert := range ms.alerts {
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (ms *MemoryStorage) UpdateAlert(alertID string, alert *core.Alert) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.alerts[alertID]; !exists {
		return NewStorageError("UpdateAlert", "alert", fmt.Errorf("alert with ID %s not found", alertID))
	}

	ms.alerts[alertID] = alert
	return nil
}

func (ms *MemoryStorage) DeleteAlert(alertID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.alerts[alertID]; !exists {
		return NewStorageError("DeleteAlert", "alert", fmt.Errorf("alert with ID %s not found", alertID))
	}

	delete(ms.alerts, alertID)
	return nil
}

func (ms *MemoryStorage) SaveAlertTrigger(trigger *core.AlertTrigger) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.triggers[trigger.ID] = trigger
	return nil
}

func (ms *MemoryStorage) GetAlertTriggers(alertID string) ([]*core.AlertTrigger, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var triggers []*core.AlertTrigger
	for _, trigger := range ms.triggers {
		if trigger.AlertID == alertID {
			triggers = append(triggers, trigger)
		}
	}

	return triggers, nil
}

func (ms *MemoryStorage) GetAllAlertTriggers() ([]*core.AlertTrigger, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var triggers []*core.AlertTrigger
	for _, trigger := range ms.triggers {
		triggers = append(triggers, trigger)
	}

	return triggers, nil
}

func (ms *MemoryStorage) SaveAlertTemplate(template *core.AlertTemplate) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.templates[template.ID] = template
	return nil
}

func (ms *MemoryStorage) GetAlertTemplate(templateID string) (*core.AlertTemplate, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	template, exists := ms.templates[templateID]
	if !exists {
		return nil, NewStorageError("GetAlertTemplate", "template", fmt.Errorf("template with ID %s not found", templateID))
	}

	return template, nil
}

func (ms *MemoryStorage) GetAllAlertTemplates() ([]*core.AlertTemplate, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var templates []*core.AlertTemplate
	for _, template := range ms.templates {
		templates = append(templates, template)
	}

	return templates, nil
}

func (ms *MemoryStorage) UpdateAlertTemplate(templateID string, template *core.AlertTemplate) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.templates[templateID]; !exists {
		return NewStorageError("UpdateAlertTemplate", "template", fmt.Errorf("template with ID %s not found", templateID))
	}

	ms.templates[templateID] = template
	return nil
}

func (ms *MemoryStorage) DeleteAlertTemplate(templateID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.templates[templateID]; !exists {
		return NewStorageError("DeleteAlertTemplate", "template", fmt.Errorf("template with ID %s not found", templateID))
	}

	delete(ms.templates, templateID)
	return nil
}

func (ms *MemoryStorage) SaveTestHistory(history *core.LoadTestHistory) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.history[history.TestID] = history
	return nil
}

func (ms *MemoryStorage) GetTestHistory(userID string) ([]*core.LoadTestHistory, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var history []*core.LoadTestHistory
	for _, h := range ms.history {
		history = append(history, h)
	}

	return history, nil
}

func (ms *MemoryStorage) GetAllTestHistory() ([]*core.LoadTestHistory, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var history []*core.LoadTestHistory
	for _, h := range ms.history {
		history = append(history, h)
	}

	return history, nil
}

func (ms *MemoryStorage) GetTestStats() (*TestStats, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	stats := &TestStats{}

	for _, test := range ms.tests {
		stats.TotalTests++
		stats.TotalRequests += test.TotalRequests
		stats.TotalDuration += test.TotalDuration
		stats.AverageResponseTime += test.AverageResponseTime
	}

	if stats.TotalTests > 0 {
		stats.AverageResponseTime = stats.AverageResponseTime / time.Duration(stats.TotalTests)
	}

	return stats, nil
}

func (ms *MemoryStorage) GetAlertStats() (*core.AlertStats, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	stats := &core.AlertStats{}

	for _, alert := range ms.alerts {
		stats.TotalAlerts++
		if alert.IsActive {
			stats.ActiveAlerts++
		}
		if alert.TriggerCount > 0 {
			stats.TriggeredAlerts++
		}

		switch alert.Severity {
		case core.SeverityCritical:
			stats.CriticalAlerts++
		case core.SeverityHigh:
			stats.HighAlerts++
		case core.SeverityMedium:
			stats.MediumAlerts++
		case core.SeverityLow:
			stats.LowAlerts++
		}
	}

	return stats, nil
}

func (ms *MemoryStorage) CleanupOldTests(olderThan time.Time) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for testID, test := range ms.tests {
		if test.EndTime.Before(olderThan) {
			delete(ms.tests, testID)
		}
	}

	return nil
}

func (ms *MemoryStorage) CleanupOldTriggers(olderThan time.Time) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for triggerID, trigger := range ms.triggers {
		if trigger.TriggeredAt.Before(olderThan) {
			delete(ms.triggers, triggerID)
		}
	}

	return nil
}
