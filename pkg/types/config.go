package types

type Config struct {
	DatabaseConnections map[string]string `json:"database_connections"`
	AnalysisSettings    AnalysisSettings  `json:"analysis_settings"`
	AlertSettings       AlertSettings     `json:"alert_settings"`
	SecuritySettings    SecuritySettings  `json:"security_settings"`
}

type AnalysisSettings struct {
	SampleSize           int     `json:"sample_size"`
	LargeTableThreshold  int64   `json:"large_table_threshold"`
	QualityScoreMinimum  float64 `json:"quality_score_minimum"`
	AutoAnalysisInterval string  `json:"auto_analysis_interval"`
}

type AlertSettings struct {
	EnableAlerts    bool     `json:"enable_alerts"`
	EmailRecipients []string `json:"email_recipients"`
	SlackWebhook    string   `json:"slack_webhook"`
}

type SecuritySettings struct {
	EnablePIIDetection bool     `json:"enable_pii_detection"`
	PIIPatterns        []string `json:"pii_patterns"`
	RequireEncryption  bool     `json:"require_encryption"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}
