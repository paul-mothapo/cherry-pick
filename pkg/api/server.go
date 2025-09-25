package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cherry-pick/pkg/analytics"
	"github.com/cherry-pick/pkg/api/analyzer"
	"github.com/cherry-pick/pkg/api/loadbalancer"
	"github.com/cherry-pick/pkg/analyzer"
	"github.com/cherry-pick/pkg/loadbalancer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router        *gin.Engine
	port          string
	loadBalancer  *loadbalancer.LoadBalancer
	urlAnalyzer   *loadbalancer.URLAnalyzer
}

func NewServer(port string) *Server {
	router := gin.Default()

	allowedOrigins := getCORSOrigins()
	allowedMethods := getCORSMethods()
	allowedHeaders := getCORSHeaders()
	exposeHeaders := getCORSExposeHeaders()
	allowCredentials := getCORSAllowCredentials()
	maxAge := getCORSMaxAge()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     allowedMethods,
		AllowHeaders:     allowedHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	}))

	lb := loadbalancer.NewLoadBalancer("./reports")
	analyzer := loadbalancer.NewURLAnalyzer()

	server := &Server{
		router:       router,
		port:         port,
		loadBalancer: lb,
		urlAnalyzer:  analyzer,
	}

	InitializeAnalytics()

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		// @Connection routes
		connections := api.Group("/connections")
		{
			connections.GET("", s.getConnections)
			connections.POST("", s.createConnection)
			connections.POST("/:id/test", s.testConnection)
			connections.DELETE("/:id", s.deleteConnection)
		}

		// @Analysis routes
		analysis := api.Group("/analysis")
		{
			analysis.GET("/reports", s.getReports)
			analysis.GET("/:id/report", s.getReport)
			analysis.POST("/:id/analyze", s.analyzeDatabase)
		}

		// @Security routes
		security := api.Group("/security")
		{
			security.GET("/:id/issues", s.getSecurityIssues)
			security.POST("/:id/analyze", s.analyzeSecurity)
		}

		// @Optimization routes
		optimization := api.Group("/optimization")
		{
			optimization.GET("/:id/history", s.getOptimizationHistory)
			optimization.POST("/:id/optimize", s.optimizeQuery)
		}

		// @Monitoring routes
		monitoring := api.Group("/monitoring")
		{
			monitoring.GET("/alerts", s.getAlerts)
			monitoring.POST("/alerts/:id/acknowledge", s.acknowledgeAlert)
			monitoring.GET("/:id/metrics", s.getMetrics)
		}

		// @Lineage routes
		lineage := api.Group("/lineage")
		{
			lineage.GET("/:id", s.getLineage)
			lineage.POST("/:id/track", s.trackLineage)
		}

		// @Collection routes
		collections := api.Group("/collections")
		{
			collections.GET("/:id/:collection/data", s.getCollectionData)
			collections.GET("/:id/:collection/stats", s.getCollectionStats)
			collections.POST("/:id/:collection/search", s.searchCollection)
		}

		// @Load Balancer routes
		loadBalancerService := loadbalancer.NewService(s.loadBalancer, s.urlAnalyzer)
		loadBalancerHandler := loadbalancer.NewHandler(loadBalancerService)
		loadbalancer.SetupRoutes(api, loadBalancerHandler)

		// @Analytics routes
		analyticsService := analytics.NewAnalytics()
		analyticsHandler := analytics.NewHandler(analyticsService.GetService())
		analytics.SetupRoutes(api, analyticsHandler)

		// @Analyzer routes
		analyzerService := analyzer.NewAnalyzer()
		analyzerHandler := analyzer.NewHandler(analyzerService.GetService())
		analyzer.SetupRoutes(api, analyzerHandler)
	}

	s.router.Static("/static", "./web/dist/assets")
	s.router.StaticFile("/", "./web/dist/index.html")
	s.router.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.port)
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (s *Server) sendSuccess(c *gin.Context, data interface{}, message ...string) {
	response := APIResponse{
		Success: true,
		Data:    data,
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) sendError(c *gin.Context, statusCode int, err error, message ...string) {
	response := APIResponse{
		Success: false,
		Error:   err.Error(),
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(statusCode, response)
}

func getCORSOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:3000", "http://localhost:8080"}
	}
	return strings.Split(origins, ",")
}

func getCORSMethods() []string {
	methods := os.Getenv("CORS_ALLOWED_METHODS")
	if methods == "" {
		return []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	}
	return strings.Split(methods, ",")
}

func getCORSHeaders() []string {
	headers := os.Getenv("CORS_ALLOWED_HEADERS")
	if headers == "" {
		return []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	}
	return strings.Split(headers, ",")
}

func getCORSExposeHeaders() []string {
	headers := os.Getenv("CORS_EXPOSE_HEADERS")
	if headers == "" {
		return []string{"Content-Length", "Content-Type"}
	}
	return strings.Split(headers, ",")
}

func getCORSAllowCredentials() bool {
	credentials := os.Getenv("CORS_ALLOW_CREDENTIALS")
	return credentials == "true" || credentials == "1"
}

func getCORSMaxAge() time.Duration {
	maxAge := os.Getenv("CORS_MAX_AGE")
	if maxAge == "" {
		return 12 * time.Hour
	}
	
	if duration, err := time.ParseDuration(maxAge); err == nil {
		return duration
	}
	
	return 12 * time.Hour
}
