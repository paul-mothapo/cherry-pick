package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cherry-pick/pkg/analytics/core"
)

type ValidatorService struct{}

func NewValidatorService() *ValidatorService {
	return &ValidatorService{}
}

func (vs *ValidatorService) ValidateEvent(event core.AnalyticsEvent) error {
	if event.ID == "" {
		return fmt.Errorf("event ID is required")
	}
	if event.Type == "" {
		return fmt.Errorf("event type is required")
	}
	if event.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if event.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	validTypes := []string{"page_view", "behavioral", "performance", "custom"}
	if !vs.isValidEventType(event.Type, validTypes) {
		return fmt.Errorf("invalid event type: %s", event.Type)
	}
	if event.Timestamp.After(time.Now()) {
		return fmt.Errorf("timestamp cannot be in the future")
	}
	if event.Timestamp.Before(time.Now().Add(-365 * 24 * time.Hour)) {
		return fmt.Errorf("timestamp is too old")
	}
	return nil
}

func (vs *ValidatorService) ValidateSession(session core.UserSession) error {
	if session.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if session.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if !vs.isValidSessionID(session.SessionID) {
		return fmt.Errorf("invalid session ID format")
	}
	if session.StartTime.After(time.Now()) {
		return fmt.Errorf("start time cannot be in the future")
	}
	if session.EndTime != nil {
		if session.EndTime.Before(session.StartTime) {
			return fmt.Errorf("end time cannot be before start time")
		}
	}
	if session.UserAgent == "" {
		return fmt.Errorf("user agent is required")
	}
	if session.IPAddress != "" && !vs.isValidIPAddress(session.IPAddress) {
		return fmt.Errorf("invalid IP address format")
	}
	return nil
}

func (vs *ValidatorService) ValidateJourney(journey core.UserJourney) error {
	if journey.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if journey.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if !vs.isValidSessionID(journey.SessionID) {
		return fmt.Errorf("invalid session ID format")
	}
	if journey.StartTime.After(time.Now()) {
		return fmt.Errorf("start time cannot be in the future")
	}
	if journey.EndTime != nil {
		if journey.EndTime.Before(journey.StartTime) {
			return fmt.Errorf("end time cannot be before start time")
		}
	}
	if journey.TotalPages < 0 {
		return fmt.Errorf("total pages cannot be negative")
	}
	if journey.TotalTime < 0 {
		return fmt.Errorf("total time cannot be negative")
	}
	if journey.ConversionRate < 0 || journey.ConversionRate > 1 {
		return fmt.Errorf("conversion rate must be between 0 and 1")
	}
	return nil
}

func (vs *ValidatorService) ValidateRequest(request core.AnalyticsRequest) error {
	if request.StartTime != nil && request.EndTime != nil {
		if request.StartTime.After(*request.EndTime) {
			return fmt.Errorf("start time cannot be after end time")
		}
		if request.EndTime.Sub(*request.StartTime) > 365*24*time.Hour {
			return fmt.Errorf("time range cannot exceed 1 year")
		}
	}
	if request.SessionID != "" && !vs.isValidSessionID(request.SessionID) {
		return fmt.Errorf("invalid session ID format")
	}
	if request.UserID != "" && !vs.isValidUserID(request.UserID) {
		return fmt.Errorf("invalid user ID format")
	}
	if request.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if request.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if request.Limit > 10000 {
		return fmt.Errorf("limit cannot exceed 10000")
	}
	return nil
}

func (vs *ValidatorService) isValidEventType(eventType string, validTypes []string) bool {
	for _, validType := range validTypes {
		if eventType == validType {
			return true
		}
	}
	return false
}

func (vs *ValidatorService) isValidSessionID(sessionID string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{10,50}$`, sessionID)
	return matched
}

func (vs *ValidatorService) isValidUserID(userID string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{5,50}$`, userID)
	return matched
}

func (vs *ValidatorService) isValidIPAddress(ip string) bool {
	ipv4Regex := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	ipv6Regex := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	return ipv4Regex.MatchString(ip) || ipv6Regex.MatchString(ip)
}

func (vs *ValidatorService) isValidURL(url string) bool {
	matched, _ := regexp.MatchString(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`, url)
	return matched
}

func (vs *ValidatorService) isValidEmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

func (vs *ValidatorService) isValidPhoneNumber(phone string) bool {
	matched, _ := regexp.MatchString(`^[\+]?[0-9\s\-\(\)]{10,20}$`, phone)
	return matched
}

func (vs *ValidatorService) isValidCountryCode(country string) bool {
	matched, _ := regexp.MatchString(`^[A-Z]{2}$`, strings.ToUpper(country))
	return matched
}

func (vs *ValidatorService) isValidDeviceType(device string) bool {
	validDevices := []string{"desktop", "mobile", "tablet", "tv", "watch", "other"}
	for _, validDevice := range validDevices {
		if strings.ToLower(device) == validDevice {
			return true
		}
	}
	return false
}

func (vs *ValidatorService) isValidBrowserType(browser string) bool {
	validBrowsers := []string{"chrome", "firefox", "safari", "edge", "opera", "ie", "other"}
	for _, validBrowser := range validBrowsers {
		if strings.ToLower(browser) == validBrowser {
			return true
		}
	}
	return false
}

func (vs *ValidatorService) isValidOSType(os string) bool {
	validOS := []string{"windows", "macos", "linux", "android", "ios", "other"}
	for _, validOS := range validOS {
		if strings.ToLower(os) == validOS {
			return true
		}
	}
	return false
}
