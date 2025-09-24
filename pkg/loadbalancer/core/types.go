package core

import "time"

type LoadTestConfig struct {
	URL             string            `json:"url" binding:"required"`
	ConcurrentUsers int               `json:"concurrentUsers" binding:"required,min=1,max=1000"`
	Duration        time.Duration     `json:"duration"`
	RampUpTime      time.Duration     `json:"rampUpTime"`
	RequestDelay    time.Duration     `json:"requestDelay"`
	Headers         map[string]string `json:"headers"`
	Method          string            `json:"method"`
	Body            string            `json:"body"`
}

type LoadTestResult struct {
	RequestID    string        `json:"requestId"`
	UserID       int           `json:"userId"`
	StartTime    time.Time     `json:"startTime"`
	EndTime      time.Time     `json:"endTime"`
	Duration     time.Duration `json:"duration"`
	StatusCode   int           `json:"statusCode"`
	ResponseSize int64         `json:"responseSize"`
	Error        string        `json:"error,omitempty"`
	Success      bool          `json:"success"`
}

type LoadTestSummary struct {
	TestID                   string           `json:"testId"`
	Config                   LoadTestConfig   `json:"config"`
	StartTime                time.Time        `json:"startTime"`
	EndTime                  time.Time        `json:"endTime"`
	TotalDuration            time.Duration    `json:"totalDuration"`
	TotalRequests            int64            `json:"totalRequests"`
	SuccessfulRequests       int64            `json:"successfulRequests"`
	FailedRequests           int64            `json:"failedRequests"`
	AverageResponseTime      time.Duration    `json:"averageResponseTime"`
	MinResponseTime          time.Duration    `json:"minResponseTime"`
	MaxResponseTime          time.Duration    `json:"maxResponseTime"`
	RequestsPerSecond        float64          `json:"requestsPerSecond"`
	ErrorRate                float64          `json:"errorRate"`
	StatusCodes              map[int]int64    `json:"statusCodes"`
	ResponseTimeDistribution map[string]int64 `json:"responseTimeDistribution"`
	Results                  []LoadTestResult `json:"results,omitempty"`
}

type LoadTestStatus struct {
	TestID    string    `json:"testId"`
	Status    string    `json:"status"`
	Progress  float64   `json:"progress"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type RealTimeMetrics struct {
	TestID              string        `json:"testId"`
	Timestamp           time.Time     `json:"timestamp"`
	ActiveUsers         int           `json:"activeUsers"`
	RequestsPerSecond   float64       `json:"requestsPerSecond"`
	AverageResponseTime time.Duration `json:"averageResponseTime"`
	ErrorRate           float64       `json:"errorRate"`
	TotalRequests       int64         `json:"totalRequests"`
	SuccessfulRequests  int64         `json:"successfulRequests"`
	FailedRequests      int64         `json:"failedRequests"`
	Percentile50        time.Duration `json:"percentile50"`
	Percentile95        time.Duration `json:"percentile95"`
	Percentile99        time.Duration `json:"percentile99"`
	Throughput          float64       `json:"throughput"`
	Bandwidth           float64       `json:"bandwidth"`
	MinResponseTime     time.Duration `json:"minResponseTime"`
	MaxResponseTime     time.Duration `json:"maxResponseTime"`
	StandardDeviation   time.Duration `json:"standardDeviation"`
	Variance            float64       `json:"variance"`
}

type LoadTestRequest struct {
	URL             string            `json:"url" binding:"required"`
	ConcurrentUsers int               `json:"concurrentUsers" binding:"required,min=1,max=1000"`
	Duration        int               `json:"duration"`
	RampUpTime      int               `json:"rampUpTime"`
	RequestDelay    int               `json:"requestDelay"`
	Headers         map[string]string `json:"headers,omitempty"`
	Method          string            `json:"method,omitempty"`
	Body            string            `json:"body,omitempty"`
}

type LoadTestResponse struct {
	TestID  string `json:"testId"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type LoadTestHistory struct {
	TestID    string          `json:"testId"`
	Name      string          `json:"name"`
	URL       string          `json:"url"`
	StartTime time.Time       `json:"startTime"`
	EndTime   time.Time       `json:"endTime"`
	Status    string          `json:"status"`
	Summary   LoadTestSummary `json:"summary,omitempty"`
}
