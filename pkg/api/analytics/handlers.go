package analytics

import (
	"net/http"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service AnalyticsService
}

type AnalyticsService interface {
	TrackPageView(event core.PageViewEvent) error
	TrackBehavioralPattern(event core.BehavioralEvent) error
	TrackPerformance(event core.PerformanceEvent) error
	TrackCustomEvent(event core.AnalyticsEvent) error
	CreateSession(session core.UserSession) error
	GetSession(sessionID string) (*core.UserSession, error)
	UpdateSession(session core.UserSession) error
	EndSession(sessionID string) error
	GetUserJourney(sessionID string) (*core.UserJourney, error)
	GetFunnelAnalysis(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error)
	GetRealTimeMetrics() (*core.RealTimeMetrics, error)
	GetInsights(sessionID string) ([]core.AnalyticsInsight, error)
	GetAlerts() ([]core.AnalyticsAlert, error)
	GetHeatmapData(pagePath string, startTime, endTime time.Time) ([]core.HeatmapPoint, error)
	GenerateReport(request core.AnalyticsRequest) (*core.AnalyticsReport, error)
	GenerateSummary(startTime, endTime time.Time) (*core.AnalyticsSummary, error)
	GenerateInsights(startTime, endTime time.Time) ([]core.AnalyticsInsight, error)
	GenerateFunnelReport(funnelID string, startTime, endTime time.Time) (*core.FunnelAnalysis, error)
	GeneratePerformanceReport(startTime, endTime time.Time) ([]core.PerformanceEvent, error)
	GenerateBehavioralReport(startTime, endTime time.Time) ([]core.BehavioralEvent, error)
	SubscribeToRealTimeMetrics(ctx context.Context) (<-chan core.RealTimeMetrics, error)
	UnsubscribeFromRealTimeMetrics(subscriberID string) error
	CleanupOldData(olderThan time.Time) error
	GetStats() (map[string]interface{}, error)
}

func NewHandler(service AnalyticsService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) TrackPageView(c *gin.Context) {
	var event core.PageViewEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	if err := h.service.TrackPageView(event); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to track page view")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Page view tracked successfully",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) TrackBehavioralPattern(c *gin.Context) {
	var event core.BehavioralEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	if err := h.service.TrackBehavioralPattern(event); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to track behavioral pattern")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Behavioral pattern tracked successfully",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) TrackPerformance(c *gin.Context) {
	var event core.PerformanceEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	if err := h.service.TrackPerformance(event); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to track performance")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Performance tracked successfully",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) TrackCustomEvent(c *gin.Context) {
	var event core.AnalyticsEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	if err := h.service.TrackCustomEvent(event); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to track custom event")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Custom event tracked successfully",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) CreateSession(c *gin.Context) {
	var session core.UserSession
	if err := c.ShouldBindJSON(&session); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	if err := h.service.CreateSession(session); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create session")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Session created successfully",
		Data:    session,
	}
	h.sendSuccess(c, response)
}

func (h *Handler) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Session ID is required")
		return
	}

	session, err := h.service.GetSession(sessionID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Session not found")
		return
	}

	h.sendSuccess(c, session)
}

func (h *Handler) UpdateSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Session ID is required")
		return
	}

	var session core.UserSession
	if err := c.ShouldBindJSON(&session); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	session.SessionID = sessionID

	if err := h.service.UpdateSession(session); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to update session")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Session updated successfully",
		Data:    session,
	}
	h.sendSuccess(c, response)
}

func (h *Handler) EndSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Session ID is required")
		return
	}

	if err := h.service.EndSession(sessionID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to end session")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Session ended successfully",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) GetUserJourney(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Session ID is required")
		return
	}

	journey, err := h.service.GetUserJourney(sessionID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "User journey not found")
		return
	}

	h.sendSuccess(c, journey)
}

func (h *Handler) GetFunnelAnalysis(c *gin.Context) {
	funnelID := c.Param("funnelId")
	if funnelID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Funnel ID is required")
		return
	}

	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid start time format")
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid end time format")
			return
		}
	} else {
		endTime = time.Now()
	}

	funnelAnalysis, err := h.service.GetFunnelAnalysis(funnelID, startTime, endTime)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get funnel analysis")
		return
	}

	h.sendSuccess(c, funnelAnalysis)
}

func (h *Handler) GetRealTimeMetrics(c *gin.Context) {
	metrics, err := h.service.GetRealTimeMetrics()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get real-time metrics")
		return
	}

	h.sendSuccess(c, metrics)
}

func (h *Handler) GetInsights(c *gin.Context) {
	sessionID := c.Query("sessionId")
	
	var insights []core.AnalyticsInsight
	var err error
	
	if sessionID != "" {
		insights, err = h.service.GetInsights(sessionID)
	} else {
		insights, err = h.service.GenerateInsights(time.Now().Add(-24*time.Hour), time.Now())
	}
	
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get insights")
		return
	}

	h.sendSuccess(c, insights)
}

func (h *Handler) GetAlerts(c *gin.Context) {
	alerts, err := h.service.GetAlerts()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get alerts")
		return
	}

	h.sendSuccess(c, alerts)
}

func (h *Handler) GetHeatmapData(c *gin.Context) {
	pagePath := c.Param("pagePath")
	if pagePath == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Page path is required")
		return
	}

	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid start time format")
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid end time format")
			return
		}
	} else {
		endTime = time.Now()
	}

	heatmapData, err := h.service.GetHeatmapData(pagePath, startTime, endTime)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get heatmap data")
		return
	}

	h.sendSuccess(c, heatmapData)
}

func (h *Handler) GenerateReport(c *gin.Context) {
	var request core.AnalyticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	report, err := h.service.GenerateReport(request)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to generate report")
		return
	}

	h.sendSuccess(c, report)
}

func (h *Handler) GenerateSummary(c *gin.Context) {
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid start time format")
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid end time format")
			return
		}
	} else {
		endTime = time.Now()
	}

	summary, err := h.service.GenerateSummary(startTime, endTime)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to generate summary")
		return
	}

	h.sendSuccess(c, summary)
}

func (h *Handler) GenerateFunnelReport(c *gin.Context) {
	funnelID := c.Param("funnelId")
	if funnelID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Funnel ID is required")
		return
	}

	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid start time format")
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid end time format")
			return
		}
	} else {
		endTime = time.Now()
	}

	report, err := h.service.GenerateFunnelReport(funnelID, startTime, endTime)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to generate funnel report")
		return
	}

	h.sendSuccess(c, report)
}

func (h *Handler) GeneratePerformanceReport(c *gin.Context) {
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid start time format")
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid end time format")
			return
		}
	} else {
		endTime = time.Now()
	}

	report, err := h.service.GeneratePerformanceReport(startTime, endTime)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to generate performance report")
		return
	}

	h.sendSuccess(c, report)
}

func (h *Handler) GenerateBehavioralReport(c *gin.Context) {
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid start time format")
			return
		}
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, err, "Invalid end time format")
			return
		}
	} else {
		endTime = time.Now()
	}

	report, err := h.service.GenerateBehavioralReport(startTime, endTime)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to generate behavioral report")
		return
	}

	h.sendSuccess(c, report)
}

func (h *Handler) SubscribeToRealTimeMetrics(c *gin.Context) {
	response := core.AnalyticsResponse{
		Success: true,
		Message: "Real-time metrics subscription not implemented yet",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) CleanupOldData(c *gin.Context) {
	olderThanStr := c.Query("olderThan")
	if olderThanStr == "" {
		h.sendError(c, http.StatusBadRequest, nil, "olderThan parameter is required")
		return
	}

	olderThan, err := time.Parse(time.RFC3339, olderThanStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid olderThan format")
		return
	}

	if err := h.service.CleanupOldData(olderThan); err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to cleanup old data")
		return
	}

	response := core.AnalyticsResponse{
		Success: true,
		Message: "Old data cleaned up successfully",
	}
	h.sendSuccess(c, response)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get stats")
		return
	}

	h.sendSuccess(c, stats)
}

func (h *Handler) sendSuccess(c *gin.Context, data interface{}, message ...string) {
	response := core.AnalyticsResponse{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) sendError(c *gin.Context, statusCode int, err error, message ...string) {
	response := core.AnalyticsResponse{
		Success: false,
		Error:   err.Error(),
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(statusCode, response)
}
