package loadbalancer

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateTestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		testID := c.Param("testId")
		if testID == "" {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Test ID is required",
			})
			c.Abort()
			return
		}
		
		if len(testID) < 3 {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Test ID must be at least 3 characters",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

func ValidateURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req URLAnalysisRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Invalid request data",
			})
			c.Abort()
			return
		}
		
		if req.URL == "" {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "URL is required",
			})
			c.Abort()
			return
		}
		
		if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "URL must start with http:// or https://",
			})
			c.Abort()
			return
		}
		
		c.Set("url", req.URL)
		c.Next()
	}
}

func ValidateLoadTestRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoadTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Invalid request data",
			})
			c.Abort()
			return
		}
		
		if req.URL == "" {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "URL is required",
			})
			c.Abort()
			return
		}
		
		if req.ConcurrentUsers < 1 || req.ConcurrentUsers > 1000 {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Concurrent users must be between 1 and 1000",
			})
			c.Abort()
			return
		}
		
		if req.Duration < 0 {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Duration cannot be negative",
			})
			c.Abort()
			return
		}
		
		if req.RequestDelay < 0 {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "Request delay cannot be negative",
			})
			c.Abort()
			return
		}
		
		c.Set("loadTestRequest", req)
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}
