package storage

import (
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type Storage interface {
	SaveTest(test *core.LoadTestSummary) error
	GetTest(testID string) (*core.LoadTestSummary, error)
	GetTestsByUser(userID string) ([]*core.LoadTestSummary, error)
	GetAllTests() ([]*core.LoadTestSummary, error)
	DeleteTest(testID string) error
	UpdateTest(testID string, test *core.LoadTestSummary) error

	SaveAlert(alert *core.Alert) error
	GetAlert(alertID string) (*core.Alert, error)
	GetAlertsByTest(testID string) ([]*core.Alert, error)
	GetAllAlerts() ([]*core.Alert, error)
	UpdateAlert(alertID string, alert *core.Alert) error
	DeleteAlert(alertID string) error

	SaveAlertTrigger(trigger *core.AlertTrigger) error
	GetAlertTriggers(alertID string) ([]*core.AlertTrigger, error)
	GetAllAlertTriggers() ([]*core.AlertTrigger, error)

	SaveAlertTemplate(template *core.AlertTemplate) error
	GetAlertTemplate(templateID string) (*core.AlertTemplate, error)
	GetAllAlertTemplates() ([]*core.AlertTemplate, error)
	UpdateAlertTemplate(templateID string, template *core.AlertTemplate) error
	DeleteAlertTemplate(templateID string) error

	SaveTestHistory(history *core.LoadTestHistory) error
	GetTestHistory(userID string) ([]*core.LoadTestHistory, error)
	GetAllTestHistory() ([]*core.LoadTestHistory, error)

	GetTestStats() (*TestStats, error)
	GetAlertStats() (*core.AlertStats, error)

	CleanupOldTests(olderThan time.Time) error
	CleanupOldTriggers(olderThan time.Time) error
}

type TestStats struct {
	TotalTests          int64         `json:"totalTests"`
	CompletedTests      int64         `json:"completedTests"`
	RunningTests        int64         `json:"runningTests"`
	FailedTests         int64         `json:"failedTests"`
	CancelledTests      int64         `json:"cancelledTests"`
	TotalRequests       int64         `json:"totalRequests"`
	TotalDuration       time.Duration `json:"totalDuration"`
	AverageResponseTime time.Duration `json:"averageResponseTime"`
}

type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"sslMode"`
	MaxConns int    `json:"maxConns"`
	MinConns int    `json:"minConns"`
}

type StorageConfig struct {
	Database DatabaseConfig `json:"database"`
	Backup   BackupConfig   `json:"backup"`
}

type BackupConfig struct {
	Enabled     bool          `json:"enabled"`
	Interval    time.Duration `json:"interval"`
	Retention   time.Duration `json:"retention"`
	Destination string        `json:"destination"`
}

type StorageError struct {
	Operation string
	Resource  string
	Err       error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("storage error in %s operation on %s: %v", e.Operation, e.Resource, e.Err)
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

func NewStorageError(operation, resource string, err error) *StorageError {
	return &StorageError{
		Operation: operation,
		Resource:  resource,
		Err:       err,
	}
}
