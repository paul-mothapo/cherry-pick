package loadbalancer

import (
	"fmt"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer"
	"github.com/cherry-pick/pkg/loadbalancer/alerting"
	"github.com/cherry-pick/pkg/loadbalancer/core"
	"github.com/cherry-pick/pkg/loadbalancer/storage"
	"github.com/cherry-pick/pkg/loadbalancer/utils"
)

type LoadBalancerService interface {
	StartLoadTest(req core.LoadTestRequest) (*core.LoadTestResponse, error)
	GetTestStatus(testID string) (*core.LoadTestStatus, error)
	GetTestSummary(testID string) (*core.LoadTestSummary, error)
	GetTestResults(testID string, includeResults bool) (map[string]interface{}, error)
	CancelTest(testID string) (*core.LoadTestResponse, error)
	GetAllTests() (map[string]*core.LoadTestStatus, error)
	GetRealTimeMetrics(testID string) (*core.RealTimeMetrics, error)
	GenerateReport(testID string) (map[string]string, error)
	GetStats() (map[string]interface{}, error)
	CleanupOldTests(olderThan time.Duration) (map[string]string, error)
	GetTestHistory() ([]core.LoadTestHistory, error)
	AnalyzeURL(url string) (*core.URLAnalysisResult, error)
	
	CreateAlert(testID string, req core.AlertRequest) (*core.Alert, error)
	GetAlert(alertID string) (*core.Alert, error)
	UpdateAlert(alertID string, req core.AlertRequest) (*core.Alert, error)
	DeleteAlert(alertID string) error
	GetAlertsForTest(testID string) ([]*core.Alert, error)
	GetAllAlerts() ([]*core.Alert, error)
	GetAlertTriggers(alertID string) ([]*core.AlertTrigger, error)
	GetAlertStats() (*core.AlertStats, error)
	EvaluateAlerts(testID string) error
	
	CreateAlertTemplate(req core.AlertTemplate) (*core.AlertTemplate, error)
	GetAlertTemplate(templateID string) (*core.AlertTemplate, error)
	UpdateAlertTemplate(templateID string, req core.AlertTemplate) (*core.AlertTemplate, error)
	DeleteAlertTemplate(templateID string) error
	GetAllAlertTemplates() ([]*core.AlertTemplate, error)
	
	GetSupportedMetrics() []string
	GetSupportedOperators() []string
	ValidateAlertCondition(condition string) error
}

type service struct {
	loadBalancer     loadbalancer.LoadBalancer
	analyzer         loadbalancer.URLAnalyzer
	metricsCalculator *utils.MetricsCalculator
	alertManager     *alerting.AlertManager
	storage          storage.Storage
}

func NewService(loadBalancer loadbalancer.LoadBalancer, analyzer loadbalancer.URLAnalyzer) LoadBalancerService {
	storage := storage.NewMemoryStorage()
	
	notifier := alerting.NewNotificationService()
	
	evaluator := alerting.NewAlertEvaluator()
	
	alertManager := alerting.NewAlertManager(notifier, evaluator)
	
	return &service{
		loadBalancer:      loadBalancer,
		analyzer:          analyzer,
		metricsCalculator: utils.NewMetricsCalculator(),
		alertManager:      alertManager,
		storage:           storage,
	}
}

func (s *service) StartLoadTest(req core.LoadTestRequest) (*core.LoadTestResponse, error) {
	testID := generateTestID()
	
	config := s.loadBalancer.ConvertRequestToConfig(req)
	
	if err := s.loadBalancer.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	if err := s.loadBalancer.StartTest(testID, config); err != nil {
		return nil, fmt.Errorf("failed to start load test: %w", err)
	}
	
	return &core.LoadTestResponse{
		TestID:  testID,
		Status:  "started",
		Message: "Load test started successfully",
	}, nil
}

func (s *service) GetTestStatus(testID string) (*core.LoadTestStatus, error) {
	return s.loadBalancer.GetTestStatus(testID)
}

func (s *service) GetTestSummary(testID string) (*core.LoadTestSummary, error) {
	return s.loadBalancer.GetTestSummary(testID)
}

func (s *service) GetTestResults(testID string, includeResults bool) (map[string]interface{}, error) {
	summary, err := s.loadBalancer.GetTestSummary(testID)
	if err != nil {
		return nil, err
	}
	
	response := map[string]interface{}{
		"summary": summary,
	}
	
	if includeResults {
		results, err := s.loadBalancer.GetTestResults(testID)
		if err != nil {
			return nil, err
		}
		response["results"] = results
	}
	
	return response, nil
}

func (s *service) CancelTest(testID string) (*core.LoadTestResponse, error) {
	if err := s.loadBalancer.CancelTest(testID); err != nil {
		return nil, err
	}
	
	return &core.LoadTestResponse{
		TestID:  testID,
		Status:  "cancelled",
		Message: "Load test cancelled successfully",
	}, nil
}

func (s *service) GetAllTests() (map[string]*core.LoadTestStatus, error) {
	return s.loadBalancer.GetAllTests(), nil
}

func (s *service) GetRealTimeMetrics(testID string) (*core.RealTimeMetrics, error) {
	basicMetrics, err := s.loadBalancer.GetRealTimeMetrics(testID)
	if err != nil {
		return nil, err
	}

	results, err := s.loadBalancer.GetTestResults(testID)
	if err != nil {
		return basicMetrics, nil
	}

	advancedMetrics := s.metricsCalculator.CalculateAdvancedMetrics(results)

	mergedMetrics := &core.RealTimeMetrics{
		TestID:              basicMetrics.TestID,
		Timestamp:           basicMetrics.Timestamp,
		ActiveUsers:         basicMetrics.ActiveUsers,
		RequestsPerSecond:   basicMetrics.RequestsPerSecond,
		AverageResponseTime: basicMetrics.AverageResponseTime,
		ErrorRate:           basicMetrics.ErrorRate,
		TotalRequests:       basicMetrics.TotalRequests,
		SuccessfulRequests:  basicMetrics.SuccessfulRequests,
		FailedRequests:     basicMetrics.FailedRequests,
		Percentile50:        advancedMetrics.Percentile50,
		Percentile95:        advancedMetrics.Percentile95,
		Percentile99:        advancedMetrics.Percentile99,
		Throughput:          advancedMetrics.Throughput,
		Bandwidth:           advancedMetrics.Bandwidth,
		MinResponseTime:     advancedMetrics.MinResponseTime,
		MaxResponseTime:     advancedMetrics.MaxResponseTime,
		StandardDeviation:   advancedMetrics.StandardDeviation,
		Variance:            advancedMetrics.Variance,
	}

	return mergedMetrics, nil
}

func (s *service) GenerateReport(testID string) (map[string]string, error) {
	if err := s.loadBalancer.GenerateReport(testID); err != nil {
		return nil, err
	}
	
	return map[string]string{
		"testId":  testID,
		"message": "Report generated successfully",
		"files":   "Check ./reports/ directory for generated files",
	}, nil
}

func (s *service) GetStats() (map[string]interface{}, error) {
	return s.loadBalancer.GetEngineStats(), nil
}

func (s *service) CleanupOldTests(olderThan time.Duration) (map[string]string, error) {
	s.loadBalancer.CleanupOldTests(olderThan)
	
	return map[string]string{
		"message":   "Cleanup completed",
		"olderThan": olderThan.String(),
	}, nil
}

func (s *service) GetTestHistory() ([]core.LoadTestHistory, error) {
	tests := s.loadBalancer.GetAllTests()
	
	var history []core.LoadTestHistory
	for testID, status := range tests {
		history = append(history, core.LoadTestHistory{
			TestID:    testID,
			Name:      "Load Test " + testID[:8],
			URL:       "N/A",
			StartTime: status.StartTime,
			EndTime:   status.EndTime,
			Status:    status.Status,
		})
	}
	
	return history, nil
}

func (s *service) AnalyzeURL(url string) (*core.URLAnalysisResult, error) {
	return s.analyzer.AnalyzeURL(url)
}


func (s *service) CreateAlert(testID string, req core.AlertRequest) (*core.Alert, error) {
	alert, err := s.alertManager.CreateAlert(testID, req)
	if err != nil {
		return nil, err
	}
	
	err = s.storage.SaveAlert(alert)
	if err != nil {
		return nil, err
	}
	
	return alert, nil
}

func (s *service) GetAlert(alertID string) (*core.Alert, error) {
	return s.storage.GetAlert(alertID)
}

func (s *service) UpdateAlert(alertID string, req core.AlertRequest) (*core.Alert, error) {
	alert, err := s.alertManager.UpdateAlert(alertID, req)
	if err != nil {
		return nil, err
	}
	
	err = s.storage.UpdateAlert(alertID, alert)
	if err != nil {
		return nil, err
	}
	
	return alert, nil
}

func (s *service) DeleteAlert(alertID string) error {
	err := s.alertManager.DeleteAlert(alertID)
	if err != nil {
		return err
	}
	
	return s.storage.DeleteAlert(alertID)
}

func (s *service) GetAlertsForTest(testID string) ([]*core.Alert, error) {
	return s.storage.GetAlertsByTest(testID)
}

func (s *service) GetAllAlerts() ([]*core.Alert, error) {
	return s.storage.GetAllAlerts()
}

func (s *service) GetAlertTriggers(alertID string) ([]*core.AlertTrigger, error) {
	return s.storage.GetAlertTriggers(alertID)
}

func (s *service) GetAlertStats() (*core.AlertStats, error) {
	return s.storage.GetAlertStats()
}

func (s *service) EvaluateAlerts(testID string) error {
	metrics, err := s.GetRealTimeMetrics(testID)
	if err != nil {
		return err
	}
	
	return s.alertManager.EvaluateAlerts(testID, metrics)
}


func (s *service) CreateAlertTemplate(req core.AlertTemplate) (*core.AlertTemplate, error) {
	template, err := s.alertManager.CreateAlertTemplate(req)
	if err != nil {
		return nil, err
	}
	
	err = s.storage.SaveAlertTemplate(template)
	if err != nil {
		return nil, err
	}
	
	return template, nil
}

func (s *service) GetAlertTemplate(templateID string) (*core.AlertTemplate, error) {
	return s.storage.GetAlertTemplate(templateID)
}

func (s *service) UpdateAlertTemplate(templateID string, req core.AlertTemplate) (*core.AlertTemplate, error) {
	err := s.storage.UpdateAlertTemplate(templateID, &req)
	if err != nil {
		return nil, err
	}
	
	return &req, nil
}

func (s *service) DeleteAlertTemplate(templateID string) error {
	return s.storage.DeleteAlertTemplate(templateID)
}

func (s *service) GetAllAlertTemplates() ([]*core.AlertTemplate, error) {
	return s.storage.GetAllAlertTemplates()
}

func (s *service) GetSupportedMetrics() []string {
	return s.alertManager.GetSupportedMetrics()
}

func (s *service) GetSupportedOperators() []string {
	return s.alertManager.GetSupportedOperators()
}

func (s *service) ValidateAlertCondition(condition string) error {
	return s.alertManager.ValidateAlertCondition(condition)
}

func generateTestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
