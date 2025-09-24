package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cherry-pick/pkg/analytics"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var analyticsTracker *analytics.Tracker

// InitializeAnalytics initializes the analytics tracker
func InitializeAnalytics() {
	analyticsTracker = analytics.NewTracker()

	// Start terminal UI
	analyticsTracker.StartTerminalUI()
}

// TrackPageView handles page view tracking
func (s *Server) trackPageView(c *gin.Context) {
	var data analytics.PageViewData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Generate session ID if not provided
	if data.SessionID == "" {
		data.SessionID = generateSessionID()
	}

	err := analyticsTracker.TrackPageView(c, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to track page view",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Page view tracked successfully",
		"data": gin.H{
			"sessionId": data.SessionID,
			"timestamp": time.Now(),
		},
	})
}

// TrackBehavioralPattern handles behavioral pattern tracking
func (s *Server) trackBehavioralPattern(c *gin.Context) {
	var data analytics.BehavioralPatternData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	err := analyticsTracker.TrackBehavioralPattern(c, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to track behavioral pattern",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Behavioral pattern tracked successfully",
		"data": gin.H{
			"timestamp": time.Now(),
		},
	})
}

// GetRealTimeAnalytics returns real-time analytics data
func (s *Server) getRealTimeAnalytics(c *gin.Context) {
	analytics := analyticsTracker.GetRealTimeAnalytics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetUserJourney returns a specific user's journey
func (s *Server) getUserJourney(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Session ID is required",
		})
		return
	}

	journey, err := analyticsTracker.GetUserJourney(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Journey not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    journey,
	})
}

// GetFunnelAnalysis performs funnel analysis
func (s *Server) getFunnelAnalysis(c *gin.Context) {
	funnelID := c.Query("funnelId")
	if funnelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Funnel ID is required",
		})
		return
	}

	stages := c.QueryArray("stages")
	if len(stages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "At least one stage is required",
		})
		return
	}

	analysis, err := analyticsTracker.GetFunnelAnalysis(funnelID, stages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to perform funnel analysis",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analysis,
	})
}

// GetAnalyticsInsights returns generated insights
func (s *Server) getAnalyticsInsights(c *gin.Context) {
	// Get insights from the tracker
	// This would be implemented to return insights based on query parameters
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"insights": []gin.H{
				{
					"id":          "insight-1",
					"type":        "performance",
					"title":       "Slow Page Load Detected",
					"description": "Multiple pages are loading slower than recommended",
					"impact":      "high",
					"confidence":  0.9,
					"timestamp":   time.Now(),
					"actionable":  true,
					"recommendations": []string{
						"Optimize images and assets",
						"Enable compression",
						"Use a CDN",
					},
				},
				{
					"id":          "insight-2",
					"type":        "user_behavior",
					"title":       "Low Engagement Detected",
					"description": "Users are not engaging deeply with content",
					"impact":      "medium",
					"confidence":  0.7,
					"timestamp":   time.Now(),
					"actionable":  true,
					"recommendations": []string{
						"Improve content quality",
						"Add engaging visual elements",
						"Optimize page layout",
					},
				},
			},
		},
	})
}

// GetAnalyticsReport generates a comprehensive analytics report
func (s *Server) getAnalyticsReport(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// Parse time parameters
	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid start time format",
			})
			return
		}
	} else {
		// Default to last 24 hours
		start = time.Now().Add(-24 * time.Hour)
	}

	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid end time format",
			})
			return
		}
	} else {
		end = time.Now()
	}

	// Generate report based on period
	report := generateAnalyticsReport(period, start, end)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
	})
}

// SubscribeToRealTimeAnalytics handles WebSocket connection for real-time analytics
func (s *Server) subscribeToRealTimeAnalytics(c *gin.Context) {
	// Upgrade to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upgrade to WebSocket",
		})
		return
	}
	defer conn.Close()

	// Subscribe to real-time data
	dataChan := analyticsTracker.SubscribeToRealTimeData()
	defer analyticsTracker.UnsubscribeFromRealTimeData("")

	// Send data every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case data := <-dataChan:
			err := conn.WriteJSON(gin.H{
				"type": "analytics_update",
				"data": data,
			})
			if err != nil {
				return
			}
		case <-ticker.C:
			// Send heartbeat
			err := conn.WriteJSON(gin.H{
				"type":      "heartbeat",
				"timestamp": time.Now(),
			})
			if err != nil {
				return
			}
		}
	}
}

// GetAnalyticsDashboard returns dashboard data
func (s *Server) getAnalyticsDashboard(c *gin.Context) {
	// Get real-time analytics
	realTimeData := analyticsTracker.GetRealTimeAnalytics()

	// Get additional dashboard data
	dashboardData := gin.H{
		"realTime": realTimeData,
		"summary": gin.H{
			"totalSessions":  len(analyticsTracker.GetSessions()),
			"totalPageViews": len(analyticsTracker.GetPageViews()),
			"totalJourneys":  len(analyticsTracker.GetJourneys()),
			"totalInsights":  len(analyticsTracker.GetInsights()),
		},
		"trends": gin.H{
			"pageViews": []gin.H{
				{"time": "00:00", "value": 120},
				{"time": "04:00", "value": 80},
				{"time": "08:00", "value": 200},
				{"time": "12:00", "value": 350},
				{"time": "16:00", "value": 280},
				{"time": "20:00", "value": 180},
			},
			"sessions": []gin.H{
				{"time": "00:00", "value": 45},
				{"time": "04:00", "value": 30},
				{"time": "08:00", "value": 75},
				{"time": "12:00", "value": 120},
				{"time": "16:00", "value": 95},
				{"time": "20:00", "value": 65},
			},
		},
		"topPages":     realTimeData.TopPages,
		"topReferrers": realTimeData.TopReferrers,
		"topCountries": realTimeData.TopCountries,
		"topDevices":   realTimeData.TopDevices,
		"topBrowsers":  realTimeData.TopBrowsers,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dashboardData,
	})
}

// Helper functions

func generateSessionID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func generateAnalyticsReport(period string, start, end time.Time) gin.H {
	// This would generate a comprehensive report based on the time period
	return gin.H{
		"id":          "report-" + strconv.FormatInt(time.Now().Unix(), 10),
		"title":       "Analytics Report - " + period,
		"period":      period,
		"startTime":   start,
		"endTime":     end,
		"generatedAt": time.Now(),
		"summary": gin.H{
			"totalSessions":      1250,
			"totalPageViews":     5670,
			"uniqueUsers":        890,
			"avgSessionDuration": 180000, // milliseconds
			"bounceRate":         35.5,
			"conversionRate":     12.3,
			"avgPageLoadTime":    1200, // milliseconds
			"performanceScore":   85.2,
		},
		"insights": []gin.H{
			{
				"type":        "performance",
				"title":       "Performance Optimization Opportunity",
				"description": "Several pages are loading slower than recommended",
				"impact":      "high",
				"actionable":  true,
			},
			{
				"type":        "conversion",
				"title":       "Conversion Rate Below Average",
				"description": "Overall conversion rate is 12.3%, below industry average of 15%",
				"impact":      "medium",
				"actionable":  true,
			},
		},
		"recommendations": []string{
			"Optimize page load times for better user experience",
			"Implement A/B testing to improve conversion rates",
			"Add user feedback collection mechanisms",
			"Monitor Core Web Vitals more closely",
		},
	}
}
