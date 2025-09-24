package loadbalancer

import "time"

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type URLAnalysisRequest struct {
	URL string `json:"url" binding:"required"`
}

type LoadTestRequest struct {
	URL             string            `json:"url" binding:"required"`
	ConcurrentUsers int               `json:"concurrentUsers" binding:"required,min=1,max=1000"`
	Duration        int               `json:"duration"`     // in seconds
	RampUpTime      int               `json:"rampUpTime"`   // in seconds
	RequestDelay    int               `json:"requestDelay"` // in milliseconds
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
	TestID    string    `json:"testId"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Status    string    `json:"status"`
}

type EngineStats struct {
	TotalTests      int `json:"totalTests"`
	RunningTests    int `json:"runningTests"`
	CompletedTests  int `json:"completedTests"`
	FailedTests     int `json:"failedTests"`
	CancelledTests  int `json:"cancelledTests"`
}

type CleanupResponse struct {
	Message   string `json:"message"`
	OlderThan string `json:"olderThan"`
}

type ReportResponse struct {
	TestID  string `json:"testId"`
	Message string `json:"message"`
	Files   string `json:"files"`
}
