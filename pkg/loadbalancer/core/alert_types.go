package core

import (
	"time"
)

type Alert struct {
	ID              string            `json:"id"`
	TestID          string            `json:"testId"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Condition       string            `json:"condition"`
	Threshold       float64           `json:"threshold"`
	Operator        string            `json:"operator"`
	Metric          string            `json:"metric"`
	IsActive        bool              `json:"isActive"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	LastTriggered   *time.Time        `json:"lastTriggered,omitempty"`
	TriggerCount    int64             `json:"triggerCount"`
	Notifications   []Notification    `json:"notifications"`
	CooldownPeriod  time.Duration     `json:"cooldownPeriod"`
	Severity        AlertSeverity     `json:"severity"`
	Tags            []string          `json:"tags"`
}

type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "low"
	SeverityMedium   AlertSeverity = "medium"
	SeverityHigh     AlertSeverity = "high"
	SeverityCritical AlertSeverity = "critical"
)

type Notification struct {
	ID       string             `json:"id"`
	Type     NotificationType   `json:"type"`
	Config   NotificationConfig `json:"config"`
	IsActive bool               `json:"isActive"`
}

type NotificationType string

const (
	NotificationEmail   NotificationType = "email"
	NotificationSlack   NotificationType = "slack"
	NotificationWebhook NotificationType = "webhook"
	NotificationSMS     NotificationType = "sms"
)

type NotificationConfig struct {
	EmailAddresses []string          `json:"emailAddresses,omitempty"`
	EmailSubject   string            `json:"emailSubject,omitempty"`
	SlackWebhookURL string           `json:"slackWebhookUrl,omitempty"`
	SlackChannel    string           `json:"slackChannel,omitempty"`
	WebhookURL      string           `json:"webhookUrl,omitempty"`
	WebhookHeaders  map[string]string `json:"webhookHeaders,omitempty"`
	PhoneNumbers    []string         `json:"phoneNumbers,omitempty"`
}

type AlertTrigger struct {
	ID          string     `json:"id"`
	AlertID     string     `json:"alertId"`
	TestID      string     `json:"testId"`
	TriggeredAt time.Time  `json:"triggeredAt"`
	Value       float64    `json:"value"`
	Threshold   float64    `json:"threshold"`
	Message     string     `json:"message"`
	ResolvedAt  *time.Time `json:"resolvedAt,omitempty"`
	IsResolved  bool       `json:"isResolved"`
}

type AlertRule struct {
	ID        string        `json:"id"`
	AlertID   string        `json:"alertId"`
	Metric    string        `json:"metric"`
	Operator  string        `json:"operator"`
	Threshold float64       `json:"threshold"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"createdAt"`
}

type AlertTemplate struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Condition   string        `json:"condition"`
	Threshold   float64       `json:"threshold"`
	Operator    string        `json:"operator"`
	Metric      string        `json:"metric"`
	Severity    AlertSeverity `json:"severity"`
	IsPublic    bool          `json:"isPublic"`
	CreatedBy   string        `json:"createdBy"`
	CreatedAt   time.Time     `json:"createdAt"`
	Tags        []string      `json:"tags"`
}

type AlertStats struct {
	TotalAlerts     int64 `json:"totalAlerts"`
	ActiveAlerts    int64 `json:"activeAlerts"`
	TriggeredAlerts int64 `json:"triggeredAlerts"`
	ResolvedAlerts  int64 `json:"resolvedAlerts"`
	CriticalAlerts  int64 `json:"criticalAlerts"`
	HighAlerts      int64 `json:"highAlerts"`
	MediumAlerts    int64 `json:"mediumAlerts"`
	LowAlerts       int64 `json:"lowAlerts"`
}

type AlertRequest struct {
	Name           string         `json:"name" binding:"required"`
	Description    string         `json:"description"`
	Condition      string         `json:"condition" binding:"required"`
	Threshold      float64        `json:"threshold" binding:"required"`
	Operator       string         `json:"operator" binding:"required"`
	Metric         string         `json:"metric" binding:"required"`
	Severity       AlertSeverity  `json:"severity"`
	IsActive       bool           `json:"isActive"`
	CooldownPeriod int            `json:"cooldownPeriod"`
	Notifications  []Notification `json:"notifications"`
	Tags           []string       `json:"tags"`
}

type AlertResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Alert   *Alert `json:"alert,omitempty"`
}
