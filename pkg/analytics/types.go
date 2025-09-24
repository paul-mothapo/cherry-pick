package analytics

import "time"

// UserSession represents a user session with tracking data
type UserSession struct {
	SessionID string                 `json:"sessionId"`
	UserID    string                 `json:"userId,omitempty"`
	StartTime time.Time              `json:"startTime"`
	EndTime   *time.Time             `json:"endTime,omitempty"`
	UserAgent string                 `json:"userAgent"`
	IPAddress string                 `json:"ipAddress"`
	Country   string                 `json:"country,omitempty"`
	Device    string                 `json:"device,omitempty"`
	Browser   string                 `json:"browser,omitempty"`
	OS        string                 `json:"os,omitempty"`
	IsActive  bool                   `json:"isActive"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PageView represents a single page view with detailed metrics
type PageView struct {
	ID                     string                 `json:"id"`
	SessionID              string                 `json:"sessionId"`
	URL                    string                 `json:"url"`
	Path                   string                 `json:"path"`
	Title                  string                 `json:"title"`
	Referrer               string                 `json:"referrer,omitempty"`
	Timestamp              time.Time              `json:"timestamp"`
	LoadTime               int64                  `json:"loadTime"`               // milliseconds
	RenderTime             int64                  `json:"renderTime"`             // milliseconds
	FirstPaint             int64                  `json:"firstPaint"`             // milliseconds
	FirstContentfulPaint   int64                  `json:"firstContentfulPaint"`   // milliseconds
	LargestContentfulPaint int64                  `json:"largestContentfulPaint"` // milliseconds
	CumulativeLayoutShift  float64                `json:"cumulativeLayoutShift"`
	FirstInputDelay        int64                  `json:"firstInputDelay"` // milliseconds
	TimeOnPage             int64                  `json:"timeOnPage"`      // milliseconds
	ScrollDepth            float64                `json:"scrollDepth"`     // percentage 0-100
	BounceRate             bool                   `json:"bounceRate"`
	ExitRate               bool                   `json:"exitRate"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}

// UserJourney represents a complete user journey through the application
type UserJourney struct {
	SessionID      string     `json:"sessionId"`
	UserID         string     `json:"userId,omitempty"`
	StartTime      time.Time  `json:"startTime"`
	EndTime        *time.Time `json:"endTime,omitempty"`
	PageViews      []PageView `json:"pageViews"`
	TotalPages     int        `json:"totalPages"`
	TotalTime      int64      `json:"totalTime"` // milliseconds
	BounceRate     bool       `json:"bounceRate"`
	ConversionRate float64    `json:"conversionRate"`
	GoalCompleted  bool       `json:"goalCompleted"`
	FunnelStage    string     `json:"funnelStage,omitempty"`
	JourneyPath    []string   `json:"journeyPath"` // array of paths visited
	DropOffPoint   string     `json:"dropOffPoint,omitempty"`
}

// PerformanceMetrics represents detailed performance data
type PerformanceMetrics struct {
	PageID                 string    `json:"pageId"`
	URL                    string    `json:"url"`
	Timestamp              time.Time `json:"timestamp"`
	LoadTime               int64     `json:"loadTime"`
	RenderTime             int64     `json:"renderTime"`
	FirstPaint             int64     `json:"firstPaint"`
	FirstContentfulPaint   int64     `json:"firstContentfulPaint"`
	LargestContentfulPaint int64     `json:"largestContentfulPaint"`
	CumulativeLayoutShift  float64   `json:"cumulativeLayoutShift"`
	FirstInputDelay        int64     `json:"firstInputDelay"`
	TimeToInteractive      int64     `json:"timeToInteractive"`
	TotalBlockingTime      int64     `json:"totalBlockingTime"`
	SpeedIndex             int64     `json:"speedIndex"`
	ResourceCount          int       `json:"resourceCount"`
	ResourceSize           int64     `json:"resourceSize"`
	ImageCount             int       `json:"imageCount"`
	ImageSize              int64     `json:"imageSize"`
	ScriptCount            int       `json:"scriptCount"`
	ScriptSize             int64     `json:"scriptSize"`
	StyleCount             int       `json:"styleCount"`
	StyleSize              int64     `json:"styleSize"`
	FontCount              int       `json:"fontCount"`
	FontSize               int64     `json:"fontSize"`
	ThirdPartyCount        int       `json:"thirdPartyCount"`
	ThirdPartySize         int64     `json:"thirdPartySize"`
	CacheHitRate           float64   `json:"cacheHitRate"`
	CDNHitRate             float64   `json:"cdnHitRate"`
	CompressionRatio       float64   `json:"compressionRatio"`
	HTTP2Usage             bool      `json:"http2Usage"`
	HTTPSUsage             bool      `json:"httpsUsage"`
	ServiceWorkerUsage     bool      `json:"serviceWorkerUsage"`
	PWAFeatures            []string  `json:"pwaFeatures,omitempty"`
}

// BehavioralPattern represents user behavior analysis
type BehavioralPattern struct {
	SessionID   string                 `json:"sessionId"`
	UserID      string                 `json:"userId,omitempty"`
	PatternType string                 `json:"patternType"` // "scroll", "click", "hover", "form_fill", "search", "navigation"
	Element     string                 `json:"element,omitempty"`
	Coordinates map[string]float64     `json:"coordinates,omitempty"`
	Duration    int64                  `json:"duration"`  // milliseconds
	Intensity   float64                `json:"intensity"` // 0-1 scale
	Frequency   int                    `json:"frequency"`
	Timestamp   time.Time              `json:"timestamp"`
	Context     map[string]interface{} `json:"context,omitempty"`
	HeatmapData []HeatmapPoint         `json:"heatmapData,omitempty"`
}

// HeatmapPoint represents a point on a heatmap
type HeatmapPoint struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Intensity float64 `json:"intensity"`
	Count     int     `json:"count"`
}

// FunnelAnalysis represents conversion funnel analysis
type FunnelAnalysis struct {
	FunnelID        string             `json:"funnelId"`
	FunnelName      string             `json:"funnelName"`
	Stages          []FunnelStage      `json:"stages"`
	TotalUsers      int                `json:"totalUsers"`
	ConversionRate  float64            `json:"conversionRate"`
	DropOffRates    []float64          `json:"dropOffRates"`
	Bottlenecks     []FunnelBottleneck `json:"bottlenecks"`
	Insights        []string           `json:"insights"`
	Recommendations []string           `json:"recommendations"`
}

// FunnelStage represents a stage in the conversion funnel
type FunnelStage struct {
	StageID        string  `json:"stageId"`
	StageName      string  `json:"stageName"`
	PagePath       string  `json:"pagePath"`
	Users          int     `json:"users"`
	ConversionRate float64 `json:"conversionRate"`
	AverageTime    int64   `json:"averageTime"` // milliseconds
	BounceRate     float64 `json:"bounceRate"`
	ExitRate       float64 `json:"exitRate"`
}

// FunnelBottleneck represents a bottleneck in the funnel
type FunnelBottleneck struct {
	StageID         string   `json:"stageId"`
	StageName       string   `json:"stageName"`
	DropOffRate     float64  `json:"dropOffRate"`
	Severity        string   `json:"severity"` // "low", "medium", "high", "critical"
	Impact          float64  `json:"impact"`   // percentage impact on overall conversion
	RootCause       string   `json:"rootCause"`
	Recommendations []string `json:"recommendations"`
}

// RealTimeAnalytics represents real-time analytics data
type RealTimeAnalytics struct {
	Timestamp          time.Time        `json:"timestamp"`
	ActiveUsers        int              `json:"activeUsers"`
	ActiveSessions     int              `json:"activeSessions"`
	PageViewsPerMinute int              `json:"pageViewsPerMinute"`
	TopPages           []PageStats      `json:"topPages"`
	TopReferrers       []ReferrerStats  `json:"topReferrers"`
	TopCountries       []CountryStats   `json:"topCountries"`
	TopDevices         []DeviceStats    `json:"topDevices"`
	TopBrowsers        []BrowserStats   `json:"topBrowsers"`
	PerformanceScore   float64          `json:"performanceScore"`
	ErrorRate          float64          `json:"errorRate"`
	BounceRate         float64          `json:"bounceRate"`
	ConversionRate     float64          `json:"conversionRate"`
	Alerts             []AnalyticsAlert `json:"alerts,omitempty"`
}

// PageStats represents statistics for a specific page
type PageStats struct {
	Path        string  `json:"path"`
	Title       string  `json:"title"`
	Views       int     `json:"views"`
	UniqueViews int     `json:"uniqueViews"`
	AvgTime     int64   `json:"avgTime"` // milliseconds
	BounceRate  float64 `json:"bounceRate"`
	ExitRate    float64 `json:"exitRate"`
	LoadTime    int64   `json:"loadTime"` // milliseconds
}

// ReferrerStats represents referrer statistics
type ReferrerStats struct {
	Referrer   string  `json:"referrer"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// CountryStats represents country statistics
type CountryStats struct {
	Country    string  `json:"country"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// DeviceStats represents device statistics
type DeviceStats struct {
	Device     string  `json:"device"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// BrowserStats represents browser statistics
type BrowserStats struct {
	Browser    string  `json:"browser"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// AnalyticsAlert represents an analytics alert
type AnalyticsAlert struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`     // "performance", "error", "conversion", "traffic"
	Severity  string                 `json:"severity"` // "low", "medium", "high", "critical"
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Resolved  bool                   `json:"resolved"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AnalyticsInsight represents a generated insight
type AnalyticsInsight struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"` // "performance", "conversion", "user_behavior", "technical"
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Impact          string                 `json:"impact"`     // "low", "medium", "high"
	Confidence      float64                `json:"confidence"` // 0-1 scale
	Timestamp       time.Time              `json:"timestamp"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Recommendations []string               `json:"recommendations,omitempty"`
	Actionable      bool                   `json:"actionable"`
}

// AnalyticsReport represents a comprehensive analytics report
type AnalyticsReport struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Period          string                 `json:"period"` // "hour", "day", "week", "month"
	StartTime       time.Time              `json:"startTime"`
	EndTime         time.Time              `json:"endTime"`
	GeneratedAt     time.Time              `json:"generatedAt"`
	Summary         AnalyticsSummary       `json:"summary"`
	UserJourneys    []UserJourney          `json:"userJourneys,omitempty"`
	FunnelAnalysis  *FunnelAnalysis        `json:"funnelAnalysis,omitempty"`
	PerformanceData []PerformanceMetrics   `json:"performanceData,omitempty"`
	BehavioralData  []BehavioralPattern    `json:"behavioralData,omitempty"`
	Insights        []AnalyticsInsight     `json:"insights"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// AnalyticsSummary represents a summary of analytics data
type AnalyticsSummary struct {
	TotalSessions      int     `json:"totalSessions"`
	TotalPageViews     int     `json:"totalPageViews"`
	UniqueUsers        int     `json:"uniqueUsers"`
	AvgSessionDuration int64   `json:"avgSessionDuration"` // milliseconds
	BounceRate         float64 `json:"bounceRate"`
	ConversionRate     float64 `json:"conversionRate"`
	AvgPageLoadTime    int64   `json:"avgPageLoadTime"` // milliseconds
	PerformanceScore   float64 `json:"performanceScore"`
	TopPage            string  `json:"topPage"`
	TopReferrer        string  `json:"topReferrer"`
	TopCountry         string  `json:"topCountry"`
	TopDevice          string  `json:"topDevice"`
	TopBrowser         string  `json:"topBrowser"`
}
