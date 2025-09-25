package analyzer

import (
	"net/http"
	"strconv"

	"github.com/cherry-pick/pkg/analyzer/core"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service AnalyzerService
}

type AnalyzerService interface {
	AnalyzeDatabase(ctx context.Context, request core.AnalysisRequest) (*core.AnalysisResult, error)
	GetAnalysisHistory(ctx context.Context, limit int) ([]core.AnalysisResult, error)
	GetAnalysisByID(ctx context.Context, analysisID string) (*core.AnalysisResult, error)
	DeleteAnalysis(ctx context.Context, analysisID string) error
	GetSupportedDatabaseTypes() []core.DatabaseType
	GetAnalysisOptions() core.AnalysisOptions
}

func NewHandler(service AnalyzerService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) AnalyzeDatabase(c *gin.Context) {
	var request core.AnalysisRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	result, err := h.service.AnalyzeDatabase(c.Request.Context(), request)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to analyze database")
		return
	}

	h.sendSuccess(c, result, "Database analysis completed successfully")
}

func (h *Handler) GetAnalysisHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, err, "Invalid limit parameter")
		return
	}

	history, err := h.service.GetAnalysisHistory(c.Request.Context(), limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to get analysis history")
		return
	}

	h.sendSuccess(c, history, "Analysis history retrieved successfully")
}

func (h *Handler) GetAnalysisByID(c *gin.Context) {
	analysisID := c.Param("id")
	if analysisID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Analysis ID is required")
		return
	}

	result, err := h.service.GetAnalysisByID(c.Request.Context(), analysisID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, err, "Analysis not found")
		return
	}

	h.sendSuccess(c, result, "Analysis retrieved successfully")
}

func (h *Handler) DeleteAnalysis(c *gin.Context) {
	analysisID := c.Param("id")
	if analysisID == "" {
		h.sendError(c, http.StatusBadRequest, nil, "Analysis ID is required")
		return
	}

	err := h.service.DeleteAnalysis(c.Request.Context(), analysisID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err, "Failed to delete analysis")
		return
	}

	h.sendSuccess(c, nil, "Analysis deleted successfully")
}

func (h *Handler) GetSupportedDatabaseTypes(c *gin.Context) {
	types := h.service.GetSupportedDatabaseTypes()
	h.sendSuccess(c, types, "Supported database types retrieved successfully")
}

func (h *Handler) GetAnalysisOptions(c *gin.Context) {
	options := h.service.GetAnalysisOptions()
	h.sendSuccess(c, options, "Analysis options retrieved successfully")
}

func (h *Handler) sendSuccess(c *gin.Context, data interface{}, message ...string) {
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	if len(message) > 0 {
		response["message"] = message[0]
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) sendError(c *gin.Context, statusCode int, err error, message ...string) {
	response := map[string]interface{}{
		"success": false,
	}
	if err != nil {
		response["error"] = err.Error()
	}
	if len(message) > 0 {
		response["message"] = message[0]
	}
	c.JSON(statusCode, response)
}
