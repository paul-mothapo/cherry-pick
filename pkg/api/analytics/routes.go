package analytics

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	analytics := router.Group("/analytics")
	{
		analytics.POST("/track/pageview", handler.TrackPageView)
		analytics.POST("/track/behavior", handler.TrackBehavioralPattern)
		analytics.POST("/track/performance", handler.TrackPerformance)
		analytics.POST("/track/event", handler.TrackCustomEvent)
		
		analytics.POST("/sessions", handler.CreateSession)
		analytics.GET("/sessions/:sessionId", handler.GetSession)
		analytics.PUT("/sessions/:sessionId", handler.UpdateSession)
		analytics.DELETE("/sessions/:sessionId", handler.EndSession)
		
		analytics.GET("/journey/:sessionId", handler.GetUserJourney)
		analytics.GET("/funnel/:funnelId", handler.GetFunnelAnalysis)
		analytics.GET("/realtime", handler.GetRealTimeMetrics)
		analytics.GET("/insights", handler.GetInsights)
		analytics.GET("/alerts", handler.GetAlerts)
		analytics.GET("/heatmap/:pagePath", handler.GetHeatmapData)
		
		analytics.POST("/reports", handler.GenerateReport)
		analytics.GET("/summary", handler.GenerateSummary)
		analytics.GET("/funnel/:funnelId/report", handler.GenerateFunnelReport)
		analytics.GET("/performance/report", handler.GeneratePerformanceReport)
		analytics.GET("/behavioral/report", handler.GenerateBehavioralReport)
		
		analytics.GET("/stream", handler.SubscribeToRealTimeMetrics)
		
		analytics.POST("/cleanup", handler.CleanupOldData)
		analytics.GET("/stats", handler.GetStats)
	}
}
