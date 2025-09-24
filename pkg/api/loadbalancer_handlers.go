package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer"
	"github.com/gin-gonic/gin"
)

var loadBalancer *loadbalancer.LoadBalancer

// InitializeLoadBalancer initializes the load balancer
func InitializeLoadBalancer() {
	loadBalancer = loadbalancer.NewLoadBalancer("./reports")
}

// URLAnalysisRequest represents the request to analyze a URL
type URLAnalysisRequest struct {
	URL string `json:"url" binding:"required"`
}

// analyzeURL analyzes a URL and discovers its pages
func (s *Server) analyzeURL(c *gin.Context) {
	var req URLAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	// Create URL analyzer
	analyzer := loadbalancer.NewURLAnalyzer()

	// Analyze the URL
	result, err := analyzer.AnalyzeURL(req.URL)
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to analyze URL")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// LoadTestRequest represents the request to start a load test
type LoadTestRequest struct {
	URL             string            `json:"url" binding:"required"`
	ConcurrentUsers int               `json:"concurrentUsers" binding:"required,min=1,max=1000"`
	Duration        int               `json:"duration"`     // in seconds
	RampUpTime      int               `json:"rampUpTime"`   // in seconds
	RequestDelay    int               `json:"requestDelay"` // in milliseconds
	Headers         map[string]string `json:"headers,omitempty"`
	Method          string            `json:"method,omitempty"`
	Body            string            `json:"body,omitempty"`
}

// LoadTestResponse represents the response for load test operations
type LoadTestResponse struct {
	TestID  string `json:"testId"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// startLoadTest starts a new load test
func (s *Server) startLoadTest(c *gin.Context) {
	var req LoadTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	// Generate test ID
	testID := generateTestID()

	// Convert request to config
	config := loadBalancer.ConvertRequestToConfig(loadbalancer.LoadTestRequest{
		URL:             req.URL,
		ConcurrentUsers: req.ConcurrentUsers,
		Duration:        req.Duration,
		RampUpTime:      req.RampUpTime,
		RequestDelay:    req.RequestDelay,
		Headers:         req.Headers,
		Method:          req.Method,
		Body:            req.Body,
	})

	// Validate config
	if err := loadBalancer.ValidateConfig(config); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Invalid configuration")
		return
	}

	// Start the test
	if err := loadBalancer.StartTest(testID, config); err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to start load test")
		return
	}

	response := LoadTestResponse{
		TestID:  testID,
		Status:  "started",
		Message: "Load test started successfully",
	}

	s.sendSuccess(c, response, "Load test started")
}

// getLoadTestStatus returns the status of a load test
func (s *Server) getLoadTestStatus(c *gin.Context) {
	testID := c.Param("testId")

	status, err := loadBalancer.GetTestStatus(testID)
	if err != nil {
		s.sendError(c, http.StatusNotFound, err, "Test not found")
		return
	}

	s.sendSuccess(c, status)
}

// getLoadTestSummary returns the summary of a completed load test
func (s *Server) getLoadTestSummary(c *gin.Context) {
	testID := c.Param("testId")

	summary, err := loadBalancer.GetTestSummary(testID)
	if err != nil {
		s.sendError(c, http.StatusNotFound, err, "Test summary not found")
		return
	}

	s.sendSuccess(c, summary)
}

// getLoadTestResults returns the detailed results of a load test
func (s *Server) getLoadTestResults(c *gin.Context) {
	testID := c.Param("testId")

	// Check if we should include detailed results
	includeResults := c.Query("includeResults") == "true"

	summary, err := loadBalancer.GetTestSummary(testID)
	if err != nil {
		s.sendError(c, http.StatusNotFound, err, "Test summary not found")
		return
	}

	response := map[string]interface{}{
		"summary": summary,
	}

	if includeResults {
		results, err := loadBalancer.GetTestResults(testID)
		if err != nil {
			s.sendError(c, http.StatusInternalServerError, err, "Failed to get test results")
			return
		}
		response["results"] = results
	}

	s.sendSuccess(c, response)
}

// cancelLoadTest cancels a running load test
func (s *Server) cancelLoadTest(c *gin.Context) {
	testID := c.Param("testId")

	if err := loadBalancer.CancelTest(testID); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Failed to cancel test")
		return
	}

	response := LoadTestResponse{
		TestID:  testID,
		Status:  "cancelled",
		Message: "Load test cancelled successfully",
	}

	s.sendSuccess(c, response, "Load test cancelled")
}

// getAllLoadTests returns all load tests
func (s *Server) getAllLoadTests(c *gin.Context) {
	tests := loadBalancer.GetAllTests()
	s.sendSuccess(c, tests)
}

// getRealTimeMetrics returns real-time metrics for a running test
func (s *Server) getRealTimeMetrics(c *gin.Context) {
	testID := c.Param("testId")

	metrics, err := loadBalancer.GetRealTimeMetrics(testID)
	if err != nil {
		s.sendError(c, http.StatusNotFound, err, "Test not found or not running")
		return
	}

	s.sendSuccess(c, metrics)
}

// generateLoadTestReport generates a comprehensive report for a test
func (s *Server) generateLoadTestReport(c *gin.Context) {
	testID := c.Param("testId")

	if err := loadBalancer.GenerateReport(testID); err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to generate report")
		return
	}

	response := map[string]string{
		"testId":  testID,
		"message": "Report generated successfully",
		"files":   "Check ./reports/ directory for generated files",
	}

	s.sendSuccess(c, response, "Report generated")
}

// getLoadBalancerStats returns statistics about the load balancer
func (s *Server) getLoadBalancerStats(c *gin.Context) {
	stats := loadBalancer.GetEngineStats()
	s.sendSuccess(c, stats)
}

// cleanupOldTests removes old completed tests
func (s *Server) cleanupOldTests(c *gin.Context) {
	// Get hours parameter from query string
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		s.sendError(c, http.StatusBadRequest, err, "Invalid hours parameter")
		return
	}

	olderThan := time.Duration(hours) * time.Hour
	loadBalancer.CleanupOldTests(olderThan)

	response := map[string]string{
		"message":   "Cleanup completed",
		"olderThan": olderThan.String(),
	}

	s.sendSuccess(c, response, "Cleanup completed")
}

// getLoadTestHistory returns a list of historical load tests
func (s *Server) getLoadTestHistory(c *gin.Context) {
	tests := loadBalancer.GetAllTests()

	var history []loadbalancer.LoadTestHistory
	for testID, status := range tests {
		history = append(history, loadbalancer.LoadTestHistory{
			TestID:    testID,
			Name:      "Load Test " + testID[:8],
			URL:       "N/A", // Would need to store this separately
			StartTime: status.StartTime,
			EndTime:   status.EndTime,
			Status:    status.Status,
		})
	}

	s.sendSuccess(c, history)
}

// generateTestID generates a unique test ID
func generateTestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
