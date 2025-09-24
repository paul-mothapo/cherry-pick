package alerting

import (
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type AlertManager struct {
	alerts      map[string]*core.Alert
	triggers    map[string]*core.AlertTrigger
	templates   map[string]*core.AlertTemplate
	mu          sync.RWMutex
	notifier    NotificationService
	evaluator   AlertEvaluator
}

type NotificationService interface {
	SendNotification(alert *core.Alert, trigger *core.AlertTrigger) error
}

type AlertEvaluator interface {
	EvaluateAlert(alert *core.Alert, metrics *core.RealTimeMetrics) (bool, float64, error)
}

func NewAlertManager(notifier NotificationService, evaluator AlertEvaluator) *AlertManager {
	return &AlertManager{
		alerts:    make(map[string]*core.Alert),
		triggers:  make(map[string]*core.AlertTrigger),
		templates: make(map[string]*core.AlertTemplate),
		notifier:  notifier,
		evaluator: evaluator,
	}
}

func (am *AlertManager) CreateAlert(testID string, req core.AlertRequest) (*core.Alert, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alertID := generateAlertID()
	now := time.Now()

	alert := &core.Alert{
		ID:             alertID,
		TestID:         testID,
		Name:           req.Name,
		Description:    req.Description,
		Condition:      req.Condition,
		Threshold:      req.Threshold,
		Operator:       req.Operator,
		Metric:         req.Metric,
		IsActive:       req.IsActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		TriggerCount:   0,
		Notifications:  req.Notifications,
		CooldownPeriod: time.Duration(req.CooldownPeriod) * time.Second,
		Severity:       req.Severity,
		Tags:           req.Tags,
	}

	am.alerts[alertID] = alert
	return alert, nil
}

func (am *AlertManager) GetAlert(alertID string) (*core.Alert, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return nil, fmt.Errorf("alert with ID %s not found", alertID)
	}

	return alert, nil
}

func (am *AlertManager) UpdateAlert(alertID string, req core.AlertRequest) (*core.Alert, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.alerts[alertID]
	if !exists {
		return nil, fmt.Errorf("alert with ID %s not found", alertID)
	}

	alert.Name = req.Name
	alert.Description = req.Description
	alert.Condition = req.Condition
	alert.Threshold = req.Threshold
	alert.Operator = req.Operator
	alert.Metric = req.Metric
	alert.IsActive = req.IsActive
	alert.UpdatedAt = time.Now()
	alert.CooldownPeriod = time.Duration(req.CooldownPeriod) * time.Second
	alert.Notifications = req.Notifications
	alert.Severity = req.Severity
	alert.Tags = req.Tags

	return alert, nil
}

func (am *AlertManager) DeleteAlert(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.alerts[alertID]; !exists {
		return fmt.Errorf("alert with ID %s not found", alertID)
	}

	delete(am.alerts, alertID)
	return nil
}

func (am *AlertManager) GetAlertsForTest(testID string) ([]*core.Alert, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var alerts []*core.Alert
	for _, alert := range am.alerts {
		if alert.TestID == testID {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

func (am *AlertManager) GetAllAlerts() ([]*core.Alert, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var alerts []*core.Alert
	for _, alert := range am.alerts {
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (am *AlertManager) EvaluateAlerts(testID string, metrics *core.RealTimeMetrics) error {
	am.mu.RLock()
	alerts := make([]*core.Alert, 0)
	for _, alert := range am.alerts {
		if alert.TestID == testID && alert.IsActive {
			alerts = append(alerts, alert)
		}
	}
	am.mu.RUnlock()

	for _, alert := range alerts {
		if alert.LastTriggered != nil && 
		   time.Since(*alert.LastTriggered) < alert.CooldownPeriod {
			continue
		}

		triggered, value, err := am.evaluator.EvaluateAlert(alert, metrics)
		if err != nil {
			continue
		}

		if triggered {
			err := am.triggerAlert(alert, value, metrics)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (am *AlertManager) triggerAlert(alert *core.Alert, value float64, metrics *core.RealTimeMetrics) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now()
	triggerID := generateTriggerID()

	trigger := &core.AlertTrigger{
		ID:          triggerID,
		AlertID:     alert.ID,
		TestID:      alert.TestID,
		TriggeredAt: now,
		Value:       value,
		Threshold:   alert.Threshold,
		Message:     fmt.Sprintf("Alert '%s' triggered: %s %s %.2f (threshold: %.2f)", 
			alert.Name, alert.Metric, alert.Operator, value, alert.Threshold),
		IsResolved:  false,
	}

	am.triggers[triggerID] = trigger

	alert.LastTriggered = &now
	alert.TriggerCount++

	for _, notification := range alert.Notifications {
		if notification.IsActive {
			err := am.notifier.SendNotification(alert, trigger)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (am *AlertManager) GetAlertTriggers(alertID string) ([]*core.AlertTrigger, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var triggers []*core.AlertTrigger
	for _, trigger := range am.triggers {
		if trigger.AlertID == alertID {
			triggers = append(triggers, trigger)
		}
	}

	return triggers, nil
}

func (am *AlertManager) GetAlertStats() (*core.AlertStats, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	stats := &core.AlertStats{}
	
	for _, alert := range am.alerts {
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

func (am *AlertManager) CreateAlertTemplate(req core.AlertTemplate) (*core.AlertTemplate, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	templateID := generateTemplateID()
	now := time.Now()

	template := &core.AlertTemplate{
		ID:          templateID,
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Threshold:   req.Threshold,
		Operator:    req.Operator,
		Metric:      req.Metric,
		Severity:    req.Severity,
		IsPublic:    req.IsPublic,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   now,
		Tags:        req.Tags,
	}

	am.templates[templateID] = template
	return template, nil
}

func (am *AlertManager) GetAlertTemplates() ([]*core.AlertTemplate, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var templates []*core.AlertTemplate
	for _, template := range am.templates {
		templates = append(templates, template)
	}

	return templates, nil
}

func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}

func generateTriggerID() string {
	return fmt.Sprintf("trigger_%d", time.Now().UnixNano())
}

func generateTemplateID() string {
	return fmt.Sprintf("template_%d", time.Now().UnixNano())
}

func (am *AlertManager) GetSupportedMetrics() []string {
	return []string{
		"error_rate",
		"response_time",
		"average_response_time",
		"requests_per_second",
		"throughput",
		"total_requests",
		"successful_requests",
		"failed_requests",
		"percentile_50",
		"p50",
		"percentile_95",
		"p95",
		"percentile_99",
		"p99",
		"bandwidth",
		"min_response_time",
		"max_response_time",
		"standard_deviation",
		"variance",
	}
}

func (am *AlertManager) GetSupportedOperators() []string {
	return []string{">", "<", ">=", "<=", "==", "!="}
}

func (am *AlertManager) ValidateAlertCondition(condition string) error {
	return nil
}
