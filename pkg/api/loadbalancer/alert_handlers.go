package loadbalancer

import (
	"net/http"
	"strconv"

	"github.com/cherry-pick/pkg/loadbalancer/core"
	"github.com/gin-gonic/gin"
)


func (h *Handler) CreateAlert(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	var req core.AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	alert, err := h.service.CreateAlert(testID, req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create alert")
		return
	}

	response := core.AlertResponse{
		ID:      alert.ID,
		Status:  "created",
		Message: "Alert created successfully",
		Alert:   alert,
	}

	h.sendSuccess(c, response, "Alert created successfully")
}

func (h *Handler) GetAlert(c *gin.Context) {
	alertID := c.Param("alertId")
	if alertID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Alert ID is required")
		return
	}

	alert, err := h.service.GetAlert(alertID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Alert not found")
		return
	}

	h.sendSuccess(c, alert)
}

func (h *Handler) UpdateAlert(c *gin.Context) {
	alertID := c.Param("alertId")
	if alertID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Alert ID is required")
		return
	}

	var req core.AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	alert, err := h.service.UpdateAlert(alertID, req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to update alert")
		return
	}

	response := core.AlertResponse{
		ID:      alert.ID,
		Status:  "updated",
		Message: "Alert updated successfully",
		Alert:   alert,
	}

	h.sendSuccess(c, response, "Alert updated successfully")
}

func (h *Handler) DeleteAlert(c *gin.Context) {
	alertID := c.Param("alertId")
	if alertID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Alert ID is required")
		return
	}

	err := h.service.DeleteAlert(alertID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to delete alert")
		return
	}

	response := core.AlertResponse{
		ID:      alertID,
		Status:  "deleted",
		Message: "Alert deleted successfully",
	}

	h.sendSuccess(c, response, "Alert deleted successfully")
}

func (h *Handler) GetAlertsForTest(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	alerts, err := h.service.GetAlertsForTest(testID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get alerts")
		return
	}

	h.sendSuccess(c, alerts)
}

func (h *Handler) GetAllAlerts(c *gin.Context) {
	alerts, err := h.service.GetAllAlerts()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get alerts")
		return
	}

	h.sendSuccess(c, alerts)
}

func (h *Handler) GetAlertTriggers(c *gin.Context) {
	alertID := c.Param("alertId")
	if alertID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Alert ID is required")
		return
	}

	triggers, err := h.service.GetAlertTriggers(alertID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get alert triggers")
		return
	}

	h.sendSuccess(c, triggers)
}

func (h *Handler) GetAlertStats(c *gin.Context) {
	stats, err := h.service.GetAlertStats()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get alert statistics")
		return
	}

	h.sendSuccess(c, stats)
}


func (h *Handler) CreateAlertTemplate(c *gin.Context) {
	var req core.AlertTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	template, err := h.service.CreateAlertTemplate(req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to create alert template")
		return
	}

	h.sendSuccess(c, template, "Alert template created successfully")
}

func (h *Handler) GetAlertTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Template ID is required")
		return
	}

	template, err := h.service.GetAlertTemplate(templateID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Alert template not found")
		return
	}

	h.sendSuccess(c, template)
}

func (h *Handler) UpdateAlertTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Template ID is required")
		return
	}

	var req core.AlertTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	template, err := h.service.UpdateAlertTemplate(templateID, req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to update alert template")
		return
	}

	h.sendSuccess(c, template, "Alert template updated successfully")
}

func (h *Handler) DeleteAlertTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Template ID is required")
		return
	}

	err := h.service.DeleteAlertTemplate(templateID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to delete alert template")
		return
	}

	response := map[string]string{
		"templateId": templateID,
		"status":     "deleted",
		"message":    "Alert template deleted successfully",
	}

	h.sendSuccess(c, response, "Alert template deleted successfully")
}

func (h *Handler) GetAllAlertTemplates(c *gin.Context) {
	templates, err := h.service.GetAllAlertTemplates()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get alert templates")
		return
	}

	h.sendSuccess(c, templates)
}


func (h *Handler) EvaluateAlerts(c *gin.Context) {
	testID := c.Param("testId")
	if testID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Test ID is required")
		return
	}

	err := h.service.EvaluateAlerts(testID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to evaluate alerts")
		return
	}

	response := map[string]string{
		"testId":  testID,
		"status":  "evaluated",
		"message": "Alerts evaluated successfully",
	}

	h.sendSuccess(c, response, "Alerts evaluated successfully")
}

func (h *Handler) GetSupportedMetrics(c *gin.Context) {
	metrics := h.service.GetSupportedMetrics()
	h.sendSuccess(c, metrics)
}

func (h *Handler) GetSupportedOperators(c *gin.Context) {
	operators := h.service.GetSupportedOperators()
	h.sendSuccess(c, operators)
}

func (h *Handler) ValidateAlertCondition(c *gin.Context) {
	var req struct {
		Condition string `json:"condition" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	err := h.service.ValidateAlertCondition(req.Condition)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid alert condition")
		return
	}

	response := map[string]string{
		"condition": req.Condition,
		"status":    "valid",
		"message":   "Alert condition is valid",
	}

	h.sendSuccess(c, response, "Alert condition is valid")
}
