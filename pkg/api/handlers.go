package api

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/intelligence"
	"github.com/cherry-pick/pkg/types"
	"github.com/gin-gonic/gin"
)

var (
	connections = make(map[string]*ConnectionInfo)
	reports     = make(map[string]*types.DatabaseReport)
	services    = make(map[string]*intelligence.Service)
	mutex       sync.RWMutex
)

type ConnectionInfo struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Driver           string     `json:"driver"`
	ConnectionString string     `json:"connectionString"`
	Status           string     `json:"status"`
	LastConnected    *time.Time `json:"lastConnected,omitempty"`
}

type CreateConnectionRequest struct {
	Name             string `json:"name" binding:"required"`
	Driver           string `json:"driver" binding:"required"`
	ConnectionString string `json:"connectionString" binding:"required"`
}

type OptimizeQueryRequest struct {
	Query string `json:"query" binding:"required"`
}

func (s *Server) getConnections(c *gin.Context) {
	mutex.RLock()
	defer mutex.RUnlock()

	connectionList := make([]*ConnectionInfo, 0, len(connections))
	for _, conn := range connections {
		connectionList = append(connectionList, conn)
	}

	s.sendSuccess(c, connectionList)
}

func (s *Server) createConnection(c *gin.Context) {
	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	id := strconv.FormatInt(time.Now().UnixNano(), 36)
	connection := &ConnectionInfo{
		ID:               id,
		Name:             req.Name,
		Driver:           req.Driver,
		ConnectionString: req.ConnectionString,
		Status:           "disconnected",
	}

	connections[id] = connection
	s.sendSuccess(c, connection, "Connection created successfully")
}

func (s *Server) testConnection(c *gin.Context) {
	id := c.Param("id")

	mutex.Lock()
	defer mutex.Unlock()

	connection, exists := connections[id]
	if !exists {
		s.sendError(c, http.StatusNotFound,
			&APIError{Message: "Connection not found"}, "Connection not found")
		return
	}

	service, err := intelligence.CreateSimpleService(connection.Driver, connection.ConnectionString)
	if err != nil {
		connection.Status = "error"
		s.sendError(c, http.StatusBadRequest, err, "Failed to connect to database")
		return
	}

	services[id] = service

	now := time.Now()
	connection.Status = "connected"
	connection.LastConnected = &now

	s.sendSuccess(c, map[string]string{"status": "connected"}, "Connection test successful")
}

func (s *Server) deleteConnection(c *gin.Context) {
	id := c.Param("id")

	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := connections[id]; !exists {
		s.sendError(c, http.StatusNotFound,
			&APIError{Message: "Connection not found"}, "Connection not found")
		return
	}

	if service, exists := services[id]; exists {
		service.Close()
		delete(services, id)
	}

	delete(connections, id)
	delete(reports, id)

	s.sendSuccess(c, nil, "Connection deleted successfully")
}

func (s *Server) getReports(c *gin.Context) {
	mutex.RLock()
	defer mutex.RUnlock()

	reportList := make([]*types.DatabaseReport, 0, len(reports))
	for _, report := range reports {
		reportList = append(reportList, report)
	}

	s.sendSuccess(c, reportList)
}

func (s *Server) getReport(c *gin.Context) {
	id := c.Param("id")

	mutex.RLock()
	defer mutex.RUnlock()

	report, exists := reports[id]
	if !exists {
		s.sendError(c, http.StatusNotFound,
			&APIError{Message: "Report not found"}, "Report not found")
		return
	}

	s.sendSuccess(c, report)
}

func (s *Server) analyzeDatabase(c *gin.Context) {
	id := c.Param("id")

	mutex.RLock()
	service, serviceExists := services[id]
	connExists := false
	if _, exists := connections[id]; exists {
		connExists = true
	}
	mutex.RUnlock()

	if !connExists {
		s.sendError(c, http.StatusNotFound,
			&APIError{Message: "Connection not found"}, "Connection not found")
		return
	}

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	report, err := service.AnalyzeDatabase()
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to analyze database")
		return
	}

	mutex.Lock()
	reports[id] = report
	mutex.Unlock()

	s.sendSuccess(c, report, "Database analysis completed")
}

func (s *Server) getSecurityIssues(c *gin.Context) {
	id := c.Param("id")

	mutex.RLock()
	service, serviceExists := services[id]
	mutex.RUnlock()

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	issues, err := service.AnalyzeSecurity()
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to analyze security")
		return
	}

	s.sendSuccess(c, issues)
}

func (s *Server) analyzeSecurity(c *gin.Context) {
	s.getSecurityIssues(c)
}

func (s *Server) getOptimizationHistory(c *gin.Context) {
	s.sendSuccess(c, []interface{}{})
}

func (s *Server) optimizeQuery(c *gin.Context) {
	id := c.Param("id")
	var req OptimizeQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	mutex.RLock()
	service, serviceExists := services[id]
	mutex.RUnlock()

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	optimization, err := service.OptimizeQuery(req.Query)
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to optimize query")
		return
	}

	s.sendSuccess(c, optimization)
}

func (s *Server) getAlerts(c *gin.Context) {
	connectionID := c.Query("connectionId")

	mutex.RLock()
	service, serviceExists := services[connectionID]
	mutex.RUnlock()

	if connectionID != "" && !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	if serviceExists {
		alerts, err := service.CheckAlerts()
		if err != nil {
			s.sendError(c, http.StatusInternalServerError, err, "Failed to get alerts")
			return
		}
		s.sendSuccess(c, alerts)
	} else {
		s.sendSuccess(c, []interface{}{})
	}
}

func (s *Server) acknowledgeAlert(c *gin.Context) {
	alertID := c.Param("id")
	s.sendSuccess(c, nil, "Alert "+alertID+" acknowledged")
}

func (s *Server) getMetrics(c *gin.Context) {
	metrics := map[string]interface{}{
		"cpu_usage":          75.5,
		"memory_usage":       68.2,
		"connections":        42,
		"queries_per_second": 150,
	}
	s.sendSuccess(c, metrics)
}

func (s *Server) getLineage(c *gin.Context) {
	id := c.Param("id")

	mutex.RLock()
	service, serviceExists := services[id]
	mutex.RUnlock()

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	lineage, err := service.TrackLineage()
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to get lineage")
		return
	}

	s.sendSuccess(c, lineage)
}

func (s *Server) trackLineage(c *gin.Context) {
	s.getLineage(c)
}

type SearchCollectionRequest struct {
	Query string `json:"query" binding:"required"`
}

type CollectionDataResponse struct {
	Documents  []map[string]interface{} `json:"documents"`
	TotalCount int64                    `json:"totalCount"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
}

type FieldStats struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Count        int64       `json:"count"`
	UniqueCount  int64       `json:"uniqueCount"`
	NullCount    int64       `json:"nullCount"`
	MinValue     interface{} `json:"minValue,omitempty"`
	MaxValue     interface{} `json:"maxValue,omitempty"`
	AvgValue     interface{} `json:"avgValue,omitempty"`
	SampleValues []string    `json:"sampleValues"`
}

type CollectionStatsResponse struct {
	CollectionName string       `json:"collectionName"`
	DocumentCount  int64        `json:"documentCount"`
	Fields         []FieldStats `json:"fields"`
	Indexes        []IndexStats `json:"indexes"`
}

type IndexStats struct {
	Name     string   `json:"name"`
	Keys     []string `json:"keys"`
	IsUnique bool     `json:"isUnique"`
	Size     int64    `json:"size"`
}

func (s *Server) getCollectionData(c *gin.Context) {
	connectionID := c.Param("id")
	collectionName := c.Param("collection")

	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	mutex.RLock()
	service, serviceExists := services[connectionID]
	mutex.RUnlock()

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	mongoService := service.GetMongoService()
	if mongoService == nil {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Not a MongoDB connection"}, "Collection data only available for MongoDB")
		return
	}

	documents, totalCount, err := s.getMongoCollectionData(service, collectionName, page, limit)
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to fetch collection data")
		return
	}

	response := CollectionDataResponse{
		Documents:  documents,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}

	s.sendSuccess(c, response)
}

func (s *Server) getCollectionStats(c *gin.Context) {
	connectionID := c.Param("id")
	collectionName := c.Param("collection")

	mutex.RLock()
	service, serviceExists := services[connectionID]
	mutex.RUnlock()

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	mongoService := service.GetMongoService()
	if mongoService == nil {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Not a MongoDB connection"}, "Collection stats only available for MongoDB")
		return
	}

	stats, err := s.getMongoCollectionStats(service, collectionName)
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to fetch collection stats")
		return
	}

	s.sendSuccess(c, stats)
}

func (s *Server) searchCollection(c *gin.Context) {
	connectionID := c.Param("id")
	collectionName := c.Param("collection")

	var req SearchCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.sendError(c, http.StatusBadRequest, err, "Invalid search request")
		return
	}

	mutex.RLock()
	service, serviceExists := services[connectionID]
	mutex.RUnlock()

	if !serviceExists {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Connection not established"}, "Please test the connection first")
		return
	}

	mongoService := service.GetMongoService()
	if mongoService == nil {
		s.sendError(c, http.StatusBadRequest,
			&APIError{Message: "Not a MongoDB connection"}, "Collection search only available for MongoDB")
		return
	}

	documents, err := s.searchMongoCollection(service, collectionName, req.Query)
	if err != nil {
		s.sendError(c, http.StatusInternalServerError, err, "Failed to search collection")
		return
	}

	s.sendSuccess(c, documents)
}

type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
