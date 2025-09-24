package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type NotificationServiceImpl struct {
	httpClient *http.Client
}

func NewNotificationService() *NotificationServiceImpl {
	return &NotificationServiceImpl{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ns *NotificationServiceImpl) SendNotification(alert *core.Alert, trigger *core.AlertTrigger) error {
	for _, notification := range alert.Notifications {
		if !notification.IsActive {
			continue
		}

		switch notification.Type {
		case core.NotificationEmail:
			return ns.sendEmailNotification(alert, trigger, notification)
		case core.NotificationSlack:
			return ns.sendSlackNotification(alert, trigger, notification)
		case core.NotificationWebhook:
			return ns.sendWebhookNotification(alert, trigger, notification)
		case core.NotificationSMS:
			return ns.sendSMSNotification(alert, trigger, notification)
		default:
			return fmt.Errorf("unsupported notification type: %s", notification.Type)
		}
	}

	return nil
}

func (ns *NotificationServiceImpl) sendEmailNotification(alert *core.Alert, trigger *core.AlertTrigger, notification core.Notification) error {
	emailContent := ns.buildEmailContent(alert, trigger)
	
	fmt.Printf("EMAIL NOTIFICATION:\n")
	fmt.Printf("To: %v\n", notification.Config.EmailAddresses)
	fmt.Printf("Subject: %s\n", notification.Config.EmailSubject)
	fmt.Printf("Content: %s\n", emailContent)
	
	return nil
}

func (ns *NotificationServiceImpl) sendSlackNotification(alert *core.Alert, trigger *core.AlertTrigger, notification core.Notification) error {
	slackMessage := ns.buildSlackMessage(alert, trigger)
	
	payload := map[string]interface{}{
		"text":    slackMessage,
		"channel": notification.Config.SlackChannel,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", notification.Config.SlackWebhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create Slack request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ns.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Slack API returned error status: %d", resp.StatusCode)
	}
	
	return nil
}

func (ns *NotificationServiceImpl) sendWebhookNotification(alert *core.Alert, trigger *core.AlertTrigger, notification core.Notification) error {
	webhookPayload := ns.buildWebhookPayload(alert, trigger)
	
	jsonPayload, err := json.Marshal(webhookPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", notification.Config.WebhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	for key, value := range notification.Config.WebhookHeaders {
		req.Header.Set(key, value)
	}
	
	resp, err := ns.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook notification: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
	}
	
	return nil
}

func (ns *NotificationServiceImpl) sendSMSNotification(alert *core.Alert, trigger *core.AlertTrigger, notification core.Notification) error {
	smsContent := ns.buildSMSContent(alert, trigger)
	
	fmt.Printf("SMS NOTIFICATION:\n")
	fmt.Printf("To: %v\n", notification.Config.PhoneNumbers)
	fmt.Printf("Content: %s\n", smsContent)
	
	return nil
}

func (ns *NotificationServiceImpl) buildEmailContent(alert *core.Alert, trigger *core.AlertTrigger) string {
	subject := fmt.Sprintf("Alert: %s", alert.Name)
	if notification := alert.Notifications[0]; notification.Config.EmailSubject != "" {
		subject = notification.Config.EmailSubject
	}
	
	content := fmt.Sprintf(`
Alert Triggered: %s

Test ID: %s
Alert ID: %s
Severity: %s

Condition: %s
Current Value: %.2f
Threshold: %.2f

Message: %s

Triggered At: %s
Trigger Count: %d

Please check your load test dashboard for more details.
`, alert.Name, alert.TestID, alert.ID, alert.Severity, alert.Condition, 
   trigger.Value, trigger.Threshold, trigger.Message, 
   trigger.TriggeredAt.Format(time.RFC3339), alert.TriggerCount)
	
	return content
}

func (ns *NotificationServiceImpl) buildSlackMessage(alert *core.Alert, trigger *core.AlertTrigger) string {
	severityEmoji := ns.getSeverityEmoji(alert.Severity)
	
	message := fmt.Sprintf(`%s *Alert Triggered: %s*

*Test ID:* %s
*Alert ID:* %s
*Severity:* %s

*Condition:* %s
*Current Value:* %.2f
*Threshold:* %.2f

*Message:* %s

*Triggered At:* %s
*Trigger Count:* %d

Please check your load test dashboard for more details.`, 
		severityEmoji, alert.Name, alert.TestID, alert.ID, alert.Severity,
		alert.Condition, trigger.Value, trigger.Threshold, trigger.Message,
		trigger.TriggeredAt.Format(time.RFC3339), alert.TriggerCount)
	
	return message
}

func (ns *NotificationServiceImpl) buildWebhookPayload(alert *core.Alert, trigger *core.AlertTrigger) map[string]interface{} {
	return map[string]interface{}{
		"event": "alert_triggered",
		"timestamp": time.Now().Format(time.RFC3339),
		"alert": map[string]interface{}{
			"id": alert.ID,
			"name": alert.Name,
			"description": alert.Description,
			"severity": alert.Severity,
			"condition": alert.Condition,
			"threshold": alert.Threshold,
			"operator": alert.Operator,
			"metric": alert.Metric,
		},
		"trigger": map[string]interface{}{
			"id": trigger.ID,
			"test_id": trigger.TestID,
			"triggered_at": trigger.TriggeredAt.Format(time.RFC3339),
			"value": trigger.Value,
			"threshold": trigger.Threshold,
			"message": trigger.Message,
		},
		"test": map[string]interface{}{
			"id": alert.TestID,
		},
	}
}

func (ns *NotificationServiceImpl) buildSMSContent(alert *core.Alert, trigger *core.AlertTrigger) string {
	return fmt.Sprintf("ALERT: %s - %s (Value: %.2f, Threshold: %.2f) - Test: %s", 
		alert.Name, alert.Severity, trigger.Value, trigger.Threshold, alert.TestID)
}

func (ns *NotificationServiceImpl) getSeverityEmoji(severity core.AlertSeverity) string {
	switch severity {
	case core.SeverityCritical:
		return "🚨"
	case core.SeverityHigh:
		return "⚠️"
	case core.SeverityMedium:
		return "🔶"
	case core.SeverityLow:
		return "ℹ️"
	default:
		return "📢"
	}
}
