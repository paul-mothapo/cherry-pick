package loadbalancer

import "errors"

var (
	ErrURLRequired           = errors.New("URL is required")
	ErrInvalidConcurrentUsers = errors.New("concurrent users must be between 1 and 1000")
	ErrInvalidDuration       = errors.New("duration cannot be negative")
	ErrInvalidRequestDelay   = errors.New("request delay cannot be negative")
	ErrTestNotFound          = errors.New("test not found")
	ErrTestAlreadyExists     = errors.New("test already exists")
	ErrTestNotRunning        = errors.New("test is not running")
	ErrTestNotCompleted      = errors.New("test is not completed")
	ErrEngineNotFound        = errors.New("engine not found")
	ErrEngineAlreadyExists   = errors.New("engine already exists")
	ErrReportGenerationFailed = errors.New("failed to generate report")
	ErrInvalidReportFormat    = errors.New("invalid report format")
	ErrInvalidURL            = errors.New("invalid URL")
	ErrURLAnalysisFailed     = errors.New("URL analysis failed")
	ErrMaxDepthExceeded      = errors.New("maximum depth exceeded")
	ErrMaxPagesExceeded      = errors.New("maximum pages exceeded")
	ErrHTTPRequestFailed     = errors.New("HTTP request failed")
	ErrHTTPTimeout           = errors.New("HTTP request timeout")
	ErrInvalidHTTPResponse   = errors.New("invalid HTTP response")
)

type LoadTestError struct {
	TestID    string
	Operation string
	Err       error
}

func (e *LoadTestError) Error() string {
	return e.Operation + " failed for test " + e.TestID + ": " + e.Err.Error()
}

func (e *LoadTestError) Unwrap() error {
	return e.Err
}

func NewLoadTestError(testID, operation string, err error) *LoadTestError {
	return &LoadTestError{
		TestID:    testID,
		Operation: operation,
		Err:       err,
	}
}

type ConfigError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ConfigError) Error() string {
	return "configuration error for field '" + e.Field + "' with value '" + 
		string(rune(e.Value.(int))) + "': " + e.Message
}

func NewConfigError(field string, value interface{}, message string) *ConfigError {
	return &ConfigError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

type ValidationError struct {
	Field   string
	Value   interface{}
	Rule    string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error for field '" + e.Field + "': " + e.Message
}

func NewValidationError(field string, value interface{}, rule, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Rule:    rule,
		Message: message,
	}
}
