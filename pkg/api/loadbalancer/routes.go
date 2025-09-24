package loadbalancer

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	loadbalancer := router.Group("/loadbalancer")
	{
		// Test management
		loadbalancer.POST("/tests", handler.StartLoadTest)
		loadbalancer.GET("/tests", handler.GetAllTests)
		loadbalancer.GET("/tests/:testId", handler.GetTestStatus)
		loadbalancer.GET("/tests/:testId/summary", handler.GetTestSummary)
		loadbalancer.GET("/tests/:testId/results", handler.GetTestResults)
		loadbalancer.DELETE("/tests/:testId", handler.CancelTest)
		
		// Real-time monitoring
		loadbalancer.GET("/tests/:testId/metrics", handler.GetRealTimeMetrics)
		loadbalancer.GET("/tests/:testId/stream", handler.StreamMetrics) // WebSocket
		
		// Reporting
		loadbalancer.POST("/tests/:testId/reports", handler.GenerateReport)
		loadbalancer.GET("/tests/:testId/reports", handler.GetReports)
		
		// Test templates
		loadbalancer.GET("/templates", handler.GetTemplates)
		loadbalancer.POST("/templates", handler.CreateTemplate)
		loadbalancer.GET("/templates/:id", handler.GetTemplate)
		loadbalancer.PUT("/templates/:id", handler.UpdateTemplate)
		loadbalancer.DELETE("/templates/:id", handler.DeleteTemplate)
		
		// Scheduled tests
		loadbalancer.GET("/schedules", handler.GetSchedules)
		loadbalancer.POST("/schedules", handler.CreateSchedule)
		loadbalancer.GET("/schedules/:id", handler.GetSchedule)
		loadbalancer.PUT("/schedules/:id", handler.UpdateSchedule)
		loadbalancer.DELETE("/schedules/:id", handler.DeleteSchedule)
		
		// System management
		loadbalancer.GET("/stats", handler.GetStats)
		loadbalancer.POST("/cleanup", handler.CleanupOldTests)
		loadbalancer.GET("/history", handler.GetTestHistory)
		
		// URL analysis
		loadbalancer.POST("/analyze", handler.AnalyzeURL)
		
		// Health and status
		loadbalancer.GET("/health", handler.HealthCheck)
		loadbalancer.GET("/version", handler.GetVersion)
		
		// Alert management
		loadbalancer.POST("/tests/:testId/alerts", handler.CreateAlert)
		loadbalancer.GET("/tests/:testId/alerts", handler.GetAlertsForTest)
		loadbalancer.GET("/alerts", handler.GetAllAlerts)
		loadbalancer.GET("/alerts/:alertId", handler.GetAlert)
		loadbalancer.PUT("/alerts/:alertId", handler.UpdateAlert)
		loadbalancer.DELETE("/alerts/:alertId", handler.DeleteAlert)
		loadbalancer.GET("/alerts/:alertId/triggers", handler.GetAlertTriggers)
		loadbalancer.GET("/alerts/stats", handler.GetAlertStats)
		loadbalancer.POST("/tests/:testId/alerts/evaluate", handler.EvaluateAlerts)
		
		// Alert templates
		loadbalancer.POST("/alert-templates", handler.CreateAlertTemplate)
		loadbalancer.GET("/alert-templates", handler.GetAllAlertTemplates)
		loadbalancer.GET("/alert-templates/:templateId", handler.GetAlertTemplate)
		loadbalancer.PUT("/alert-templates/:templateId", handler.UpdateAlertTemplate)
		loadbalancer.DELETE("/alert-templates/:templateId", handler.DeleteAlertTemplate)
		
		// Alert utilities
		loadbalancer.GET("/alerts/metrics", handler.GetSupportedMetrics)
		loadbalancer.GET("/alerts/operators", handler.GetSupportedOperators)
		loadbalancer.POST("/alerts/validate", handler.ValidateAlertCondition)
	}
}
