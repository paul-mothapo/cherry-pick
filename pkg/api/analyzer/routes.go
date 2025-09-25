package analyzer

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	analyzer := router.Group("/analyzer")
	{
		analyzer.POST("/analyze", handler.AnalyzeDatabase)
		analyzer.GET("/history", handler.GetAnalysisHistory)
		analyzer.GET("/:id", handler.GetAnalysisByID)
		analyzer.DELETE("/:id", handler.DeleteAnalysis)
		analyzer.GET("/types", handler.GetSupportedDatabaseTypes)
		analyzer.GET("/options", handler.GetAnalysisOptions)
	}
}
