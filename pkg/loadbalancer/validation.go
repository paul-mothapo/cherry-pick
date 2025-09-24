package loadbalancer

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type ConfigValidator struct{}

func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

func (v *ConfigValidator) ValidateConfig(config core.LoadTestConfig) error {
	if err := v.validateURL(config.URL); err != nil {
		return err
	}
	if err := v.validateConcurrentUsers(config.ConcurrentUsers); err != nil {
		return err
	}
	if err := v.validateDuration(config.Duration); err != nil {
		return err
	}
	if err := v.validateRequestDelay(config.RequestDelay); err != nil {
		return err
	}
	if err := v.validateMethod(config.Method); err != nil {
		return err
	}
	if err := v.validateHeaders(config.Headers); err != nil {
		return err
	}
	return nil
}

func (v *ConfigValidator) validateURL(urlStr string) error {
	if urlStr == "" {
		return NewValidationError("URL", urlStr, "required", "URL is required")
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return NewValidationError("URL", urlStr, "valid_url", "invalid URL format")
	}
	if parsedURL.Scheme == "" {
		return NewValidationError("URL", urlStr, "scheme", "URL must include scheme (http/https)")
	}
	if parsedURL.Host == "" {
		return NewValidationError("URL", urlStr, "host", "URL must include host")
	}
	return nil
}

func (v *ConfigValidator) validateConcurrentUsers(users int) error {
	if users < 1 {
		return NewValidationError("ConcurrentUsers", users, "min", 
			"concurrent users must be at least 1")
	}
	if users > 1000 {
		return NewValidationError("ConcurrentUsers", users, "max", 
			"concurrent users cannot exceed 1000")
	}
	return nil
}

func (v *ConfigValidator) validateDuration(duration time.Duration) error {
	if duration < 0 {
		return NewValidationError("Duration", duration, "positive", "duration cannot be negative")
	}
	return nil
}

func (v *ConfigValidator) validateRequestDelay(delay time.Duration) error {
	if delay < 0 {
		return NewValidationError("RequestDelay", delay, "positive", "request delay cannot be negative")
	}
	return nil
}

func (v *ConfigValidator) validateMethod(method string) error {
	if method == "" {
		return nil
	}
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	for _, validMethod := range validMethods {
		if strings.ToUpper(method) == validMethod {
			return nil
		}
	}
	return NewValidationError("Method", method, "valid_method", 
		"method must be one of: GET, POST, PUT, DELETE, PATCH")
}

func (v *ConfigValidator) validateHeaders(headers map[string]string) error {
	if headers == nil {
		return nil
	}
	for key, value := range headers {
		if key == "" {
			return NewValidationError("Headers", key, "non_empty", "header key cannot be empty")
		}
		if strings.Contains(key, "\n") || strings.Contains(key, "\r") {
			return NewValidationError("Headers", key, "no_newlines", "header key cannot contain newlines")
		}
		if strings.Contains(value, "\n") || strings.Contains(value, "\r") {
			return NewValidationError("Headers", value, "no_newlines", "header value cannot contain newlines")
		}
	}
	return nil
}

func (v *ConfigValidator) ValidateTestID(testID string) error {
	if testID == "" {
		return NewValidationError("TestID", testID, "required", "test ID is required")
	}
	if len(testID) < 3 {
		return NewValidationError("TestID", testID, "min_length", "test ID must be at least 3 characters")
	}
	if len(testID) > 100 {
		return NewValidationError("TestID", testID, "max_length", "test ID cannot exceed 100 characters")
	}
	for _, char := range testID {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return NewValidationError("TestID", testID, "valid_chars", 
				"test ID can only contain alphanumeric characters, hyphens, and underscores")
		}
	}
	return nil
}

func (v *ConfigValidator) ValidateEngineID(engineID string) error {
	if engineID == "" {
		return NewValidationError("EngineID", engineID, "required", "engine ID is required")
	}
	if len(engineID) < 1 {
		return NewValidationError("EngineID", engineID, "min_length", "engine ID must be at least 1 character")
	}
	if len(engineID) > 50 {
		return NewValidationError("EngineID", engineID, "max_length", "engine ID cannot exceed 50 characters")
	}
	return nil
}