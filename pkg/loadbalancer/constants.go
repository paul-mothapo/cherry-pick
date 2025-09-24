package loadbalancer

import "time"

const (
	DefaultTestDuration        = 30 * time.Second
	DefaultRequestDelay        = 100 * time.Millisecond
	DefaultHTTPTimeout         = 30 * time.Second
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 100
	DefaultIdleConnTimeout     = 90 * time.Second
	DefaultMaxConcurrentUsers  = 1000
	DefaultMinConcurrentUsers  = 1
	DefaultURLAnalyzerTimeout  = 10 * time.Second
	DefaultMaxDepth            = 5
	DefaultMaxPages            = 200
)

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
	MethodPATCH  = "PATCH"
)

const (
	BucketUnder100ms = "<100ms"
	Bucket100To500ms = "100-500ms"
	Bucket500msTo1s  = "500ms-1s"
	Bucket1To2s      = "1-2s"
	BucketOver2s     = ">2s"
)

const (
	ErrURLRequired            = "URL is required"
	ErrInvalidConcurrentUsers = "concurrent users must be between 1 and 1000"
	ErrInvalidDuration        = "duration cannot be negative"
	ErrInvalidRequestDelay    = "request delay cannot be negative"
	ErrTestNotFound           = "test not found"
	ErrTestAlreadyExists      = "test already exists"
	ErrTestNotRunning         = "test is not running"
	ErrEngineNotFound         = "engine not found"
	ErrEngineAlreadyExists    = "engine already exists"
)

const (
	HTMLReportExtension = ".html"
	JSONReportExtension = ".json"
	CSVReportExtension  = ".csv"
)

const (
	LoadTestReportPrefix = "loadtest_"
)
