package core

import (
	"context"
	"time"
)

type LoadTestEngine interface {
	StartLoadTest(testID string, config LoadTestConfig) error
	GetTestStatus(testID string) (*LoadTestStatus, error)
	GetTestSummary(testID string) (*LoadTestSummary, error)
	GetTestResults(testID string) ([]LoadTestResult, error)
	CancelTest(testID string) error
	GetAllTests() map[string]*LoadTestStatus
	GetRealTimeMetrics(testID string) (*RealTimeMetrics, error)
	CleanupOldTests(olderThan time.Duration)
}

type TestReporter interface {
	GenerateReport(testID string, summary *LoadTestSummary, results []LoadTestResult) error
}

type TestManager interface {
	CreateEngine(engineID string) (LoadTestEngine, error)
	GetEngine(engineID string) (LoadTestEngine, error)
	DeleteEngine(engineID string) error
	ListEngines() []string
	GetDefaultEngine() LoadTestEngine
	StartLoadTest(testID string, config LoadTestConfig) error
	GetTestStatus(testID string) (*LoadTestStatus, error)
	GetTestSummary(testID string) (*LoadTestSummary, error)
	GetTestResults(testID string) ([]LoadTestResult, error)
	CancelTest(testID string) error
	GetAllTests() map[string]*LoadTestStatus
	GetRealTimeMetrics(testID string) (*RealTimeMetrics, error)
	CleanupOldTests(olderThan time.Duration)
	GetEngineStats() map[string]interface{}
}

type ConfigValidator interface {
	ValidateConfig(config LoadTestConfig) error
	ValidateTestID(testID string) error
	ValidateEngineID(engineID string) error
}

type TestExecutor interface {
	ExecuteTest(ctx context.Context, testID string, config LoadTestConfig) error
	CancelTest(testID string) error
	GetTestStatus(testID string) (*LoadTestStatus, error)
}

type ResultAggregator interface {
	AggregateResults(testID string, results []LoadTestResult) (*LoadTestSummary, error)
	GetTestResults(testID string) ([]LoadTestResult, error)
}

type MetricsCollector interface {
	CollectMetrics(testID string) (*RealTimeMetrics, error)
	GetEngineStats() map[string]interface{}
}

type ReportGenerator interface {
	GenerateHTMLReport(testID string, summary *LoadTestSummary, results []LoadTestResult) error
	GenerateJSONReport(testID string, summary *LoadTestSummary, results []LoadTestResult) error
	GenerateCSVReport(testID string, results []LoadTestResult) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoadBalancer interface {
	StartTest(testID string, config LoadTestConfig) error
	GetTestStatus(testID string) (*LoadTestStatus, error)
	GetTestSummary(testID string) (*LoadTestSummary, error)
	GetTestResults(testID string) ([]LoadTestResult, error)
	CancelTest(testID string) error
	GetAllTests() map[string]*LoadTestStatus
	GetRealTimeMetrics(testID string) (*RealTimeMetrics, error)
	GenerateReport(testID string) error
	GetEngineStats() map[string]interface{}
	CleanupOldTests(olderThan time.Duration)
	ValidateConfig(config LoadTestConfig) error
	ConvertRequestToConfig(req LoadTestRequest) LoadTestConfig
}
