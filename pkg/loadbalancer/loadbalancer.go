package loadbalancer

import (
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
	"github.com/cherry-pick/pkg/loadbalancer/manager"
	"github.com/cherry-pick/pkg/loadbalancer/reporter"
)

type LoadBalancer struct {
	manager   core.TestManager
	reporter  core.TestReporter
	validator core.ConfigValidator
}

func NewLoadBalancer(outputDir string) *LoadBalancer {
	return &LoadBalancer{
		manager:   manager.NewManager(),
		reporter:  reporter.NewReporter(outputDir),
		validator: NewConfigValidator(),
	}
}

func (lb *LoadBalancer) StartTest(testID string, config core.LoadTestConfig) error {
	if config.Method == "" {
		config.Method = "GET"
	}
	if config.Duration == 0 {
		config.Duration = 30 * time.Second
	}
	if config.RequestDelay == 0 {
		config.RequestDelay = 100 * time.Millisecond
	}

	if err := lb.validator.ValidateConfig(config); err != nil {
		return err
	}

	return lb.manager.StartLoadTest(testID, config)
}

func (lb *LoadBalancer) GetTestStatus(testID string) (*core.LoadTestStatus, error) {
	return lb.manager.GetTestStatus(testID)
}

func (lb *LoadBalancer) GetTestSummary(testID string) (*core.LoadTestSummary, error) {
	return lb.manager.GetTestSummary(testID)
}

func (lb *LoadBalancer) GetTestResults(testID string) ([]core.LoadTestResult, error) {
	return lb.manager.GetTestResults(testID)
}

func (lb *LoadBalancer) CancelTest(testID string) error {
	return lb.manager.CancelTest(testID)
}

func (lb *LoadBalancer) GetAllTests() map[string]*core.LoadTestStatus {
	return lb.manager.GetAllTests()
}

func (lb *LoadBalancer) GetRealTimeMetrics(testID string) (*core.RealTimeMetrics, error) {
	return lb.manager.GetRealTimeMetrics(testID)
}

func (lb *LoadBalancer) GenerateReport(testID string) error {
	summary, err := lb.GetTestSummary(testID)
	if err != nil {
		return fmt.Errorf("failed to get test summary: %w", err)
	}

	results, err := lb.GetTestResults(testID)
	if err != nil {
		return fmt.Errorf("failed to get test results: %w", err)
	}

	return lb.reporter.GenerateReport(testID, summary, results)
}

func (lb *LoadBalancer) GetEngineStats() map[string]interface{} {
	return lb.manager.GetEngineStats()
}

func (lb *LoadBalancer) CleanupOldTests(olderThan time.Duration) {
	lb.manager.CleanupOldTests(olderThan)
}

func (lb *LoadBalancer) ValidateConfig(config core.LoadTestConfig) error {
	return lb.validator.ValidateConfig(config)
}

func (lb *LoadBalancer) ConvertRequestToConfig(req core.LoadTestRequest) core.LoadTestConfig {
	config := core.LoadTestConfig{
		URL:             req.URL,
		ConcurrentUsers: req.ConcurrentUsers,
		Headers:         req.Headers,
		Method:          req.Method,
		Body:            req.Body,
	}

	if req.Duration > 0 {
		config.Duration = time.Duration(req.Duration) * time.Second
	}

	if req.RampUpTime > 0 {
		config.RampUpTime = time.Duration(req.RampUpTime) * time.Second
	}

	if req.RequestDelay > 0 {
		config.RequestDelay = time.Duration(req.RequestDelay) * time.Millisecond
	}

	return config
}