package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
	"github.com/slack-go/slack"
	"gopkg.in/gomail.v2"
)

type NotificationConfig struct {
	SlackToken    string
	SlackChannel  string
	EmailSMTPHost string
	EmailSMTPPort int
	EmailUsername string
	EmailPassword string
	EmailFrom     string
	EmailTo       []string
}

type NotifierService struct {
	config      *NotificationConfig
	slackClient *slack.Client
	emailDialer *gomail.Dialer
}

func NewNotifierService() *NotifierService {
	config := &NotificationConfig{
		SlackToken:    os.Getenv("SLACK_BOT_TOKEN"),
		SlackChannel:  os.Getenv("SLACK_CHANNEL"),
		EmailSMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
		EmailSMTPPort: 587,
		EmailUsername: os.Getenv("EMAIL_USERNAME"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
		EmailFrom:     os.Getenv("EMAIL_FROM"),
		EmailTo:       []string{os.Getenv("EMAIL_TO")},
	}

	var slackClient *slack.Client
	if config.SlackToken != "" {
		slackClient = slack.New(config.SlackToken)
	}

	var emailDialer *gomail.Dialer
	if config.EmailSMTPHost != "" && config.EmailUsername != "" {
		emailDialer = gomail.NewDialer(config.EmailSMTPHost, config.EmailSMTPPort, config.EmailUsername, config.EmailPassword)
	}

	return &NotifierService{
		config:      config,
		slackClient: slackClient,
		emailDialer: emailDialer,
	}
}

func (ns *NotifierService) SendAlert(alert core.AnalyticsAlert) error {	
	if ns.slackClient != nil && ns.config.SlackChannel != "" {
		if err := ns.sendSlackAlert(alert); err != nil {
			log.Printf("Failed to send Slack alert: %v", err)
		}
	}
	
	if ns.emailDialer != nil && len(ns.config.EmailTo) > 0 {
		if err := ns.sendEmailAlert(alert); err != nil {
			log.Printf("Failed to send email alert: %v", err)
		}
	}
	
	return nil
}

func (ns *NotifierService) SendInsight(insight core.AnalyticsInsight) error {
	if ns.slackClient != nil && ns.config.SlackChannel != "" {
		if err := ns.sendSlackInsight(insight); err != nil {
			log.Printf("Failed to send Slack insight: %v", err)
		}
	}
	
	if ns.emailDialer != nil && len(ns.config.EmailTo) > 0 {
		if err := ns.sendEmailInsight(insight); err != nil {
			log.Printf("Failed to send email insight: %v", err)
		}
	}
	
	return nil
}

func (ns *NotifierService) SendReport(report core.AnalyticsReport) error {
	if ns.slackClient != nil && ns.config.SlackChannel != "" {
		if err := ns.sendSlackReport(report); err != nil {
			log.Printf("Failed to send Slack report: %v", err)
		}
	}
	
	if ns.emailDialer != nil && len(ns.config.EmailTo) > 0 {
		if err := ns.sendEmailReport(report); err != nil {
			log.Printf("Failed to send email report: %v", err)
		}
	}
	
	return nil
}

func (ns *NotifierService) sendSlackAlert(alert core.AnalyticsAlert) error {
	severityColor := map[string]string{
		"low":    "#36a64f",
		"medium": "#ff9500",
		"high":   "#ff0000",
		"critical": "#8b0000",
	}
	
	color := severityColor[alert.Severity]
	if color == "" {
		color = "#36a64f"
	}
	
	attachment := slack.Attachment{
		Color: color,
		Title: fmt.Sprintf("Analytics Alert: %s", alert.Title),
		Text:  alert.Message,
		Fields: []slack.AttachmentField{
			{
				Title: "Severity",
				Value: alert.Severity,
				Short: true,
			},
			{
				Title: "Timestamp",
				Value: alert.CreatedAt.Format(time.RFC3339),
				Short: true,
			},
		},
		Footer: "Analytics System",
		Ts:     json.Number(fmt.Sprintf("%d", alert.CreatedAt.Unix())),
	}
	
	_, _, err := ns.slackClient.PostMessage(ns.config.SlackChannel, slack.MsgOptionAttachments(attachment))
	return err
}

func (ns *NotifierService) sendSlackInsight(insight core.AnalyticsInsight) error {
	attachment := slack.Attachment{
		Color: "#36a64f",
		Title: fmt.Sprintf("💡 Analytics Insight: %s", insight.Title),
		Text:  insight.Description,
		Fields: []slack.AttachmentField{
			{
				Title: "Type",
				Value: insight.Type,
				Short: true,
			},
			{
				Title: "Confidence",
				Value: fmt.Sprintf("%.2f", insight.Confidence),
				Short: true,
			},
		},
		Footer: "Analytics System",
		Ts:     json.Number(fmt.Sprintf("%d", insight.CreatedAt.Unix())),
	}
	
	_, _, err := ns.slackClient.PostMessage(ns.config.SlackChannel, slack.MsgOptionAttachments(attachment))
	return err
}

func (ns *NotifierService) sendSlackReport(report core.AnalyticsReport) error {
	attachment := slack.Attachment{
		Color: "#36a64f",
		Title: fmt.Sprintf("Analytics Report: %s", report.Title),
		Text:  fmt.Sprintf("Report generated with %d events processed", report.TotalEvents),
		Fields: []slack.AttachmentField{
			{
				Title: "Total Events",
				Value: fmt.Sprintf("%d", report.TotalEvents),
				Short: true,
			},
			{
				Title: "Duration",
				Value: report.EndTime.Sub(report.StartTime).String(),
				Short: true,
			},
		},
		Footer: "Analytics System",
		Ts:     json.Number(fmt.Sprintf("%d", report.GeneratedAt.Unix())),
	}
	
	_, _, err := ns.slackClient.PostMessage(ns.config.SlackChannel, slack.MsgOptionAttachments(attachment))
	return err
}

func (ns *NotifierService) sendEmailAlert(alert core.AnalyticsAlert) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ns.config.EmailFrom)
	m.SetHeader("To", ns.config.EmailTo...)
	m.SetHeader("Subject", fmt.Sprintf("Analytics Alert: %s", alert.Title))
	
	body := fmt.Sprintf(`
		<h2>Analytics Alert</h2>
		<p><strong>Title:</strong> %s</p>
		<p><strong>Message:</strong> %s</p>
		<p><strong>Severity:</strong> %s</p>
		<p><strong>Timestamp:</strong> %s</p>
		<hr>
		<p><em>This alert was generated by the Analytics System</em></p>
	`, alert.Title, alert.Message, alert.Severity, alert.CreatedAt.Format(time.RFC3339))
	
	m.SetBody("text/html", body)
	
	return ns.emailDialer.DialAndSend(m)
}

func (ns *NotifierService) sendEmailInsight(insight core.AnalyticsInsight) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ns.config.EmailFrom)
	m.SetHeader("To", ns.config.EmailTo...)
	m.SetHeader("Subject", fmt.Sprintf("Analytics Insight: %s", insight.Title))
	
	body := fmt.Sprintf(`
		<h2>Analytics Insight</h2>
		<p><strong>Title:</strong> %s</p>
		<p><strong>Description:</strong> %s</p>
		<p><strong>Type:</strong> %s</p>
		<p><strong>Confidence:</strong> %.2f</p>
		<p><strong>Timestamp:</strong> %s</p>
		<hr>
		<p><em>This insight was generated by the Analytics System</em></p>
	`, insight.Title, insight.Description, insight.Type, insight.Confidence, insight.CreatedAt.Format(time.RFC3339))
	
	m.SetBody("text/html", body)
	
	return ns.emailDialer.DialAndSend(m)
}

func (ns *NotifierService) sendEmailReport(report core.AnalyticsReport) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ns.config.EmailFrom)
	m.SetHeader("To", ns.config.EmailTo...)
	m.SetHeader("Subject", fmt.Sprintf("Analytics Report: %s", report.Title))
	
	body := fmt.Sprintf(`
		<h2>Analytics Report</h2>
		<p><strong>Title:</strong> %s</p>
		<p><strong>Total Events:</strong> %d</p>
		<p><strong>Start Time:</strong> %s</p>
		<p><strong>End Time:</strong> %s</p>
		<p><strong>Duration:</strong> %s</p>
		<hr>
		<p><em>This report was generated by the Analytics System</em></p>
	`, report.Title, report.TotalEvents, report.StartTime.Format(time.RFC3339), report.EndTime.Format(time.RFC3339), report.EndTime.Sub(report.StartTime).String())
	
	m.SetBody("text/html", body)
	
	return ns.emailDialer.DialAndSend(m)
}