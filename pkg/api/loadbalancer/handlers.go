package loadbalancer

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	service LoadBalancerService
	upgrader websocket.Upgrader
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func NewHandler(service LoadBalancerService) *Handler {
	return &Handler{
		service: service,
		upgrader: upgrader,
	}
}

func (h *Handler) StartLoadTest(c *gin.Context) {
	var req core.LoadTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	response, err := h.service.StartLoadTest(req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to start load test")
		return
	}

	h.sendSuccess(c, response, "Load test started successfully")
}

func (h *Handler) GetTestStatus(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	status, err := h.service.GetTestStatus(testID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Test not found")
		return
	}

	h.sendSuccess(c, status)
}

func (h *Handler) GetTestSummary(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	summary, err := h.service.GetTestSummary(testID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Test summary not found")
		return
	}

	h.sendSuccess(c, summary)
}

func (h *Handler) GetTestResults(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	includeResults := c.Query("includeResults") == "true"
	
	response, err := h.service.GetTestResults(testID, includeResults)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Test results not found")
		return
	}

	h.sendSuccess(c, response)
}

func (h *Handler) CancelTest(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	response, err := h.service.CancelTest(testID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Failed to cancel test")
		return
	}

	h.sendSuccess(c, response, "Load test cancelled successfully")
}

func (h *Handler) GetAllTests(c *gin.Context) {
	tests, err := h.service.GetAllTests()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get tests")
		return
	}

	h.sendSuccess(c, tests)
}

func (h *Handler) GetRealTimeMetrics(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	metrics, err := h.service.GetRealTimeMetrics(testID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Test not found or not running")
		return
	}

	h.sendSuccess(c, metrics)
}

func (h *Handler) GenerateReport(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	response, err := h.service.GenerateReport(testID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to generate report")
		return
	}

	h.sendSuccess(c, response, "Report generated successfully")
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get statistics")
		return
	}

	h.sendSuccess(c, stats)
}

func (h *Handler) CleanupOldTests(c *gin.Context) {
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		h.sendError(c, http.StatusBadRequest, err, "Invalid hours parameter")
		return
	}

	olderThan := time.Duration(hours) * time.Hour
	response, err := h.service.CleanupOldTests(olderThan)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to cleanup tests")
		return
	}

	h.sendSuccess(c, response, "Cleanup completed successfully")
}

func (h *Handler) GetTestHistory(c *gin.Context) {
	history, err := h.service.GetTestHistory()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get test history")
		return
	}

	h.sendSuccess(c, history)
}

func (h *Handler) AnalyzeURL(c *gin.Context) {
	var req URLAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	result, err := h.service.AnalyzeURL(req.URL)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to analyze URL")
		return
	}

	h.sendSuccess(c, result)
}

func (h *Handler) StreamMetrics(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to upgrade connection")
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := h.service.GetRealTimeMetrics(testID)
			if err != nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "error",
					"message": "Test not found or completed",
				})
				return
			}

			message := map[string]interface{}{
				"type": "metrics",
				"data": metrics,
				"timestamp": time.Now(),
			}

			if err := conn.WriteJSON(message); err != nil {
				return
			}

		case <-c.Request.Context().Done():
			return
		}
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	health := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now(),
		"version": "1.0.0",
		"uptime": time.Since(time.Now()).String(),
		"services": map[string]string{
			"loadbalancer": "healthy",
			"database": "healthy",
			"websocket": "healthy",
		},
	}

	h.sendSuccess(c, health)
}

func (h *Handler) GetVersion(c *gin.Context) {
	version := map[string]interface{}{
		"version": "1.0.0",
		"build": "2024-01-01",
		"api_version": "v1",
		"features": []string{
			"load_testing",
			"real_time_metrics",
			"websocket_streaming",
			"url_analysis",
			"report_generation",
		},
	}

	h.sendSuccess(c, version)
}

func (h *Handler) sendSuccess(c *gin.Context, data interface{}, message ...string) {
	response := APIResponse{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) sendError(c *gin.Context, statusCode int, err error, message ...string) {
	response := APIResponse{
		Success: false,
		Error:   getErrorMessage(err),
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(statusCode, response)
}

func getErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
