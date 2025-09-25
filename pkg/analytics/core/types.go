package core

import "time"

type AnalyticsEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	UserID    string                 `json:"userId,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type PageViewEvent struct {
	AnalyticsEvent
	URL                    string  `json:"url"`
	Path                   string  `json:"path"`
	Title                  string  `json:"title"`
	Referrer               string  `json:"referrer,omitempty"`
	LoadTime               int64   `json:"loadTime"`
	RenderTime             int64   `json:"renderTime"`
	FirstPaint             int64   `json:"firstPaint"`
	FirstContentfulPaint   int64   `json:"firstContentfulPaint"`
	LargestContentfulPaint int64   `json:"largestContentfulPaint"`
	CumulativeLayoutShift  float64 `json:"cumulativeLayoutShift"`
	FirstInputDelay        int64   `json:"firstInputDelay"`
	TimeOnPage             int64   `json:"timeOnPage"`
	ScrollDepth            float64 `json:"scrollDepth"`
	BounceRate             bool    `json:"bounceRate"`
	ExitRate               bool    `json:"exitRate"`
}

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

type BehavioralEvent struct {
	AnalyticsEvent
	PatternType string             `json:"patternType"`
	Element     string             `json:"element,omitempty"`
	Coordinates map[string]float64 `json:"coordinates,omitempty"`
	Duration    int64              `json:"duration"`
	Intensity   float64            `json:"intensity"`
	Frequency   int                `json:"frequency"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

type PerformanceEvent struct {
	AnalyticsEvent
	PageID                 string  `json:"pageId"`
	URL                    string  `json:"url"`
	LoadTime               int64   `json:"loadTime"`
	RenderTime             int64   `json:"renderTime"`
	FirstPaint             int64   `json:"firstPaint"`
	FirstContentfulPaint   int64   `json:"firstContentfulPaint"`
	LargestContentfulPaint int64   `json:"largestContentfulPaint"`
	CumulativeLayoutShift  float64 `json:"cumulativeLayoutShift"`
	FirstInputDelay        int64   `json:"firstInputDelay"`
	TimeToInteractive      int64   `json:"timeToInteractive"`
	TotalBlockingTime      int64   `json:"totalBlockingTime"`
	SpeedIndex             int64   `json:"speedIndex"`
	ResourceCount          int     `json:"resourceCount"`
	ResourceSize           int64   `json:"resourceSize"`
	CacheHitRate           float64 `json:"cacheHitRate"`
	CDNHitRate             float64 `json:"cdnHitRate"`
	CompressionRatio       float64 `json:"compressionRatio"`
	HTTP2Usage             bool    `json:"http2Usage"`
	HTTPSUsage             bool    `json:"httpsUsage"`
}

type UserJourney struct {
	SessionID      string     `json:"sessionId"`
	UserID         string     `json:"userId,omitempty"`
	StartTime      time.Time  `json:"startTime"`
	EndTime        *time.Time `json:"endTime,omitempty"`
	TotalPages     int        `json:"totalPages"`
	TotalTime      int64      `json:"totalTime"`
	BounceRate     bool       `json:"bounceRate"`
	ConversionRate float64    `json:"conversionRate"`
	GoalCompleted  bool       `json:"goalCompleted"`
	FunnelStage    string     `json:"funnelStage,omitempty"`
	JourneyPath    []string   `json:"journeyPath"`
	DropOffPoint   string     `json:"dropOffPoint,omitempty"`
}

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

type FunnelStage struct {
	StageID        string  `json:"stageId"`
	StageName      string  `json:"stageName"`
	PagePath       string  `json:"pagePath"`
	Users          int     `json:"users"`
	ConversionRate float64 `json:"conversionRate"`
	AverageTime    int64   `json:"averageTime"`
	BounceRate     float64 `json:"bounceRate"`
	ExitRate       float64 `json:"exitRate"`
}

type FunnelBottleneck struct {
	StageID         string   `json:"stageId"`
	StageName       string   `json:"stageName"`
	DropOffRate     float64  `json:"dropOffRate"`
	Severity        string   `json:"severity"`
	Impact          float64  `json:"impact"`
	RootCause       string   `json:"rootCause"`
	Recommendations []string `json:"recommendations"`
}

type RealTimeMetrics struct {
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

type PageStats struct {
	Path        string  `json:"path"`
	Title       string  `json:"title"`
	Views       int     `json:"views"`
	UniqueViews int     `json:"uniqueViews"`
	AvgTime     int64   `json:"avgTime"`
	BounceRate  float64 `json:"bounceRate"`
	ExitRate    float64 `json:"exitRate"`
	LoadTime    int64   `json:"loadTime"`
}

type ReferrerStats struct {
	Referrer   string  `json:"referrer"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type CountryStats struct {
	Country    string  `json:"country"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type DeviceStats struct {
	Device     string  `json:"device"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type BrowserStats struct {
	Browser    string  `json:"browser"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type AnalyticsAlert struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Severity  string                 `json:"severity"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Resolved  bool                   `json:"resolved"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type AnalyticsInsight struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Impact          string                 `json:"impact"`
	Confidence      float64                `json:"confidence"`
	Timestamp       time.Time              `json:"timestamp"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Recommendations []string               `json:"recommendations,omitempty"`
	Actionable      bool                   `json:"actionable"`
}

type AnalyticsReport struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Period          string                 `json:"period"`
	StartTime       time.Time              `json:"startTime"`
	EndTime         time.Time              `json:"endTime"`
	GeneratedAt     time.Time              `json:"generatedAt"`
	Summary         AnalyticsSummary       `json:"summary"`
	UserJourneys    []UserJourney          `json:"userJourneys,omitempty"`
	FunnelAnalysis  *FunnelAnalysis        `json:"funnelAnalysis,omitempty"`
	PerformanceData []PerformanceEvent     `json:"performanceData,omitempty"`
	BehavioralData  []BehavioralEvent      `json:"behavioralData,omitempty"`
	Insights        []AnalyticsInsight     `json:"insights"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type AnalyticsSummary struct {
	TotalSessions      int     `json:"totalSessions"`
	TotalPageViews     int     `json:"totalPageViews"`
	UniqueUsers        int     `json:"uniqueUsers"`
	AvgSessionDuration int64   `json:"avgSessionDuration"`
	BounceRate         float64 `json:"bounceRate"`
	ConversionRate     float64 `json:"conversionRate"`
	AvgPageLoadTime    int64   `json:"avgPageLoadTime"`
	PerformanceScore   float64 `json:"performanceScore"`
	TopPage            string  `json:"topPage"`
	TopReferrer        string  `json:"topReferrer"`
	TopCountry         string  `json:"topCountry"`
	TopDevice          string  `json:"topDevice"`
	TopBrowser         string  `json:"topBrowser"`
}

type HeatmapPoint struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Intensity float64 `json:"intensity"`
	Count     int     `json:"count"`
}

type AnalyticsRequest struct {
	SessionID string                 `json:"sessionId,omitempty"`
	UserID    string                 `json:"userId,omitempty"`
	StartTime *time.Time             `json:"startTime,omitempty"`
	EndTime   *time.Time             `json:"endTime,omitempty"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Limit     int                    `json:"limit,omitempty"`
	Offset    int                    `json:"offset,omitempty"`
}

type AnalyticsResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}
