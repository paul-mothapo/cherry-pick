package analytics

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
)

// TerminalUI provides a real-time terminal interface for analytics
type TerminalUI struct {
	tracker *Tracker
	running bool
}

// NewTerminalUI creates a new terminal UI
func NewTerminalUI(tracker *Tracker) *TerminalUI {
	return &TerminalUI{
		tracker: tracker,
		running: false,
	}
}

// Start starts the terminal UI
func (t *TerminalUI) Start() {
	t.running = true

	// Clear screen
	fmt.Print("\033[2J\033[H")

	// Set up signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		t.Stop()
		os.Exit(0)
	}()

	// Start the display loop
	t.displayLoop()
}

// Stop stops the terminal UI
func (t *TerminalUI) Stop() {
	t.running = false
	fmt.Print("\033[2J\033[H")
	fmt.Println("Analytics terminal stopped.")
}

// displayLoop runs the main display loop
func (t *TerminalUI) displayLoop() {
	for t.running {
		// Move cursor to top-left
		fmt.Print("\033[H")

		// Get current analytics data
		analytics := t.tracker.GetRealTimeAnalytics()

		// Display header
		t.displayHeader()

		// Display real-time metrics
		t.displayMetrics(analytics)

		// Display recent page views
		t.displayRecentPageViews()

		// Display top pages
		t.displayTopPages(analytics)

		// Display user sessions
		t.displayActiveSessions()

		// Display performance alerts
		t.displayAlerts(analytics)

		// Display footer
		t.displayFooter()

		// Wait before next update
		time.Sleep(2 * time.Second)
	}
}

// displayHeader displays the header
func (t *TerminalUI) displayHeader() {
	header := color.New(color.FgCyan, color.Bold)
	header.Println("┌─────────────────────────────────────────────────────────────────────────────┐")
	header.Printf("│                    🚀 CHERRY PICK ANALYTICS DASHBOARD 🚀                    │\n")
	header.Printf("│                           Real-Time User Activity                          │\n")
	header.Println("└─────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
}

// displayMetrics displays key metrics
func (t *TerminalUI) displayMetrics(analytics RealTimeAnalytics) {
	// Active users
	activeColor := color.New(color.FgGreen, color.Bold)
	activeColor.Printf("👥 Active Users: %d", analytics.ActiveUsers)

	// Sessions
	sessionColor := color.New(color.FgBlue, color.Bold)
	sessionColor.Printf("    📊 Sessions: %d", analytics.ActiveSessions)

	// Page views per minute
	pageColor := color.New(color.FgYellow, color.Bold)
	pageColor.Printf("    📄 Page Views/min: %d", analytics.PageViewsPerMinute)

	// Performance score
	perfColor := color.New(color.FgMagenta, color.Bold)
	perfColor.Printf("    ⚡ Performance: %.1f/100", analytics.PerformanceScore)

	fmt.Println()
	fmt.Println()
}

// displayRecentPageViews displays recent page views
func (t *TerminalUI) displayRecentPageViews() {
	title := color.New(color.FgWhite, color.Bold)
	title.Println("📋 RECENT PAGE VIEWS")

	// Get recent page views (last 10)
	recentViews := t.getRecentPageViews(10)

	if len(recentViews) == 0 {
		noData := color.New(color.FgRed)
		noData.Println("   No recent page views")
		fmt.Println()
		return
	}

	for i, pv := range recentViews {
		timestamp := pv.Timestamp.Format("15:04:05")
		path := pv.Path
		if len(path) > 40 {
			path = path[:37] + "..."
		}

		// Color code based on load time
		var loadColor *color.Color
		if pv.LoadTime < 1000 {
			loadColor = color.New(color.FgGreen)
		} else if pv.LoadTime < 3000 {
			loadColor = color.New(color.FgYellow)
		} else {
			loadColor = color.New(color.FgRed)
		}

		fmt.Printf("   %d. [%s] %-40s %s%4dms%s\n",
			i+1,
			timestamp,
			path,
			loadColor.Sprintf(""),
			pv.LoadTime,
			loadColor.Sprintf(""))
	}
	fmt.Println()
}

// displayTopPages displays top pages
func (t *TerminalUI) displayTopPages(analytics RealTimeAnalytics) {
	title := color.New(color.FgWhite, color.Bold)
	title.Println("🏆 TOP PAGES (Last Minute)")

	if len(analytics.TopPages) == 0 {
		noData := color.New(color.FgRed)
		noData.Println("   No page data available")
		fmt.Println()
		return
	}

	for i, page := range analytics.TopPages {
		if i >= 5 { // Show only top 5
			break
		}

		path := page.Path
		if len(path) > 35 {
			path = path[:32] + "..."
		}

		viewsColor := color.New(color.FgCyan)
		timeColor := color.New(color.FgGreen)

		fmt.Printf("   %d. %-35s %s%3d views%s %s%4dms%s\n",
			i+1,
			path,
			viewsColor.Sprintf(""),
			page.Views,
			viewsColor.Sprintf(""),
			timeColor.Sprintf(""),
			page.AvgTime,
			timeColor.Sprintf(""))
	}
	fmt.Println()
}

// displayActiveSessions displays active user sessions
func (t *TerminalUI) displayActiveSessions() {
	title := color.New(color.FgWhite, color.Bold)
	title.Println("👤 ACTIVE USER SESSIONS")

	sessions := t.tracker.GetSessions()
	activeCount := 0

	for _, session := range sessions {
		if session.IsActive {
			activeCount++
			if activeCount <= 5 { // Show only first 5 active sessions
				duration := time.Since(session.StartTime)
				device := session.Device
				if device == "" {
					device = "Unknown"
				}

				country := session.Country
				if country == "" {
					country = "Unknown"
				}

				sessionColor := color.New(color.FgGreen)
				fmt.Printf("   • %s | %s | %s | %s\n",
					session.SessionID[:8],
					device,
					country,
					sessionColor.Sprintf("%.0fs", duration.Seconds()))
			}
		}
	}

	if activeCount == 0 {
		noData := color.New(color.FgRed)
		noData.Println("   No active sessions")
	} else if activeCount > 5 {
		moreColor := color.New(color.FgYellow)
		moreColor.Printf("   ... and %d more active sessions\n", activeCount-5)
	}

	fmt.Println()
}

// displayAlerts displays performance alerts
func (t *TerminalUI) displayAlerts(analytics RealTimeAnalytics) {
	if len(analytics.Alerts) == 0 {
		return
	}

	title := color.New(color.FgWhite, color.Bold)
	title.Println("🚨 ACTIVE ALERTS")

	for _, alert := range analytics.Alerts {
		var alertColor *color.Color
		switch alert.Severity {
		case "critical":
			alertColor = color.New(color.FgRed, color.Bold)
		case "high":
			alertColor = color.New(color.FgRed)
		case "medium":
			alertColor = color.New(color.FgYellow)
		default:
			alertColor = color.New(color.FgBlue)
		}

		fmt.Printf("   %s [%s] %s\n",
			alertColor.Sprintf("⚠️"),
			strings.ToUpper(alert.Severity),
			alert.Title)
	}
	fmt.Println()
}

// displayFooter displays the footer
func (t *TerminalUI) displayFooter() {
	footer := color.New(color.FgCyan)
	footer.Println("┌─────────────────────────────────────────────────────────────────────────────┐")
	footer.Printf("│ Press Ctrl+C to stop the analytics dashboard                              │\n")
	footer.Printf("│ Updates every 2 seconds • %s                                              │\n", time.Now().Format("15:04:05"))
	footer.Println("└─────────────────────────────────────────────────────────────────────────────┘")
}

// getRecentPageViews gets recent page views
func (t *TerminalUI) getRecentPageViews(limit int) []PageView {
	pageViews := t.tracker.GetPageViews()

	// Sort by timestamp (most recent first)
	recent := make([]PageView, 0, limit)
	count := 0

	// Iterate in reverse order to get most recent
	for i := len(pageViews) - 1; i >= 0 && count < limit; i-- {
		recent = append(recent, pageViews[i])
		count++
	}

	return recent
}

// LogPageView logs a page view to the terminal
func (t *TerminalUI) LogPageView(pv PageView) {
	timestamp := pv.Timestamp.Format("15:04:05")
	path := pv.Path
	if len(path) > 50 {
		path = path[:47] + "..."
	}

	// Color code based on load time
	var loadColor *color.Color
	if pv.LoadTime < 1000 {
		loadColor = color.New(color.FgGreen)
	} else if pv.LoadTime < 3000 {
		loadColor = color.New(color.FgYellow)
	} else {
		loadColor = color.New(color.FgRed)
	}

	// Log the page view
	fmt.Printf("🌐 [%s] %-50s %s%4dms%s\n",
		timestamp,
		path,
		loadColor.Sprintf(""),
		pv.LoadTime,
		loadColor.Sprintf(""))
}

// LogUserAction logs a user action to the terminal
func (t *TerminalUI) LogUserAction(action, details string) {
	timestamp := time.Now().Format("15:04:05")
	actionColor := color.New(color.FgCyan)

	fmt.Printf("👆 [%s] %s: %s\n",
		timestamp,
		actionColor.Sprintf(action),
		details)
}

// LogAlert logs an alert to the terminal
func (t *TerminalUI) LogAlert(severity, message string) {
	timestamp := time.Now().Format("15:04:05")

	var alertColor *color.Color
	switch severity {
	case "critical":
		alertColor = color.New(color.FgRed, color.Bold)
	case "high":
		alertColor = color.New(color.FgRed)
	case "medium":
		alertColor = color.New(color.FgYellow)
	default:
		alertColor = color.New(color.FgBlue)
	}

	fmt.Printf("🚨 [%s] %s: %s\n",
		timestamp,
		alertColor.Sprintf(strings.ToUpper(severity)),
		message)
}
