package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type Reporter struct {
	outputDir string
}

func NewReporter(outputDir string) *Reporter {
	return &Reporter{
		outputDir: outputDir,
	}
}

func (r *Reporter) GenerateReport(testID string, summary *core.LoadTestSummary, results []core.LoadTestResult) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	jsonReport := LoadTestReport{
		TestID:    testID,
		Generated: time.Now(),
		Summary:   summary,
		Results:   results,
	}

	jsonFile := filepath.Join(r.outputDir, fmt.Sprintf("loadtest_%s.json", testID))
	if err := r.saveJSONReport(jsonFile, jsonReport); err != nil {
		return fmt.Errorf("failed to save JSON report: %w", err)
	}

	htmlFile := filepath.Join(r.outputDir, fmt.Sprintf("loadtest_%s.html", testID))
	if err := r.generateHTMLReport(htmlFile, jsonReport); err != nil {
		return fmt.Errorf("failed to generate HTML report: %w", err)
	}

	csvFile := filepath.Join(r.outputDir, fmt.Sprintf("loadtest_%s.csv", testID))
	if err := r.generateCSVReport(csvFile, results); err != nil {
		return fmt.Errorf("failed to generate CSV report: %w", err)
	}

	return nil
}

type LoadTestReport struct {
	TestID    string                 `json:"testId"`
	Generated time.Time              `json:"generated"`
	Summary   *core.LoadTestSummary  `json:"summary"`
	Results   []core.LoadTestResult  `json:"results"`
}

func (r *Reporter) saveJSONReport(filename string, report LoadTestReport) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func (r *Reporter) generateHTMLReport(filename string, report LoadTestReport) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	html := r.generateHTMLContent(report)
	_, err = file.WriteString(html)
	return err
}

func (r *Reporter) generateHTMLContent(report LoadTestReport) string {
	summary := report.Summary

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Load Test Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #e0e0e0; }
        .metric { display: inline-block; margin: 10px; padding: 15px; background: #f8f9fa; border-radius: 5px; min-width: 150px; text-align: center; }
        .metric-value { font-size: 24px; font-weight: bold; color: #2c3e50; }
        .metric-label { font-size: 14px; color: #7f8c8d; margin-top: 5px; }
        .section { margin: 20px 0; }
        .section h3 { color: #34495e; border-bottom: 1px solid #bdc3c7; padding-bottom: 5px; }
        .status-success { color: #27ae60; }
        .status-error { color: #e74c3c; }
        .table { width: 100%%; border-collapse: collapse; margin-top: 10px; }
        .table th, .table td { padding: 8px 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .table th { background-color: #f8f9fa; font-weight: bold; }
        .chart-container { margin: 20px 0; padding: 20px; background: #f8f9fa; border-radius: 5px; }
        .progress-bar { width: 100%%; height: 20px; background-color: #ecf0f1; border-radius: 10px; overflow: hidden; }
        .progress-fill { height: 100%%; background-color: #3498db; transition: width 0.3s ease; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Load Test Report</h1>
            <p>Test ID: %s | Generated: %s</p>
        </div>

        <div class="section">
            <h3>Test Configuration</h3>
            <table class="table">
                <tr><td><strong>URL</strong></td><td>%s</td></tr>
                <tr><td><strong>Concurrent Users</strong></td><td>%d</td></tr>
                <tr><td><strong>Duration</strong></td><td>%s</td></tr>
                <tr><td><strong>Method</strong></td><td>%s</td></tr>
            </table>
        </div>

        <div class="section">
            <h3>Key Metrics</h3>
            <div class="metric">
                <div class="metric-value">%d</div>
                <div class="metric-label">Total Requests</div>
            </div>
            <div class="metric">
                <div class="metric-value status-success">%d</div>
                <div class="metric-label">Successful</div>
            </div>
            <div class="metric">
                <div class="metric-value status-error">%d</div>
                <div class="metric-label">Failed</div>
            </div>
            <div class="metric">
                <div class="metric-value">%.2f%%</div>
                <div class="metric-label">Error Rate</div>
            </div>
            <div class="metric">
                <div class="metric-value">%s</div>
                <div class="metric-label">Avg Response Time</div>
            </div>
            <div class="metric">
                <div class="metric-value">%.2f</div>
                <div class="metric-label">Requests/sec</div>
            </div>
        </div>

        <div class="section">
            <h3>Response Time Distribution</h3>
            <div class="chart-container">
                %s
            </div>
        </div>

        <div class="section">
            <h3>Status Code Distribution</h3>
            <table class="table">
                <thead>
                    <tr><th>Status Code</th><th>Count</th><th>Percentage</th></tr>
                </thead>
                <tbody>
                    %s
                </tbody>
            </table>
        </div>

        <div class="section">
            <h3>Response Time Statistics</h3>
            <table class="table">
                <tr><td><strong>Minimum</strong></td><td>%s</td></tr>
                <tr><td><strong>Average</strong></td><td>%s</td></tr>
                <tr><td><strong>Maximum</strong></td><td>%s</td></tr>
            </table>
        </div>
    </div>
</body>
</html>`,
		report.TestID,
		report.TestID,
		report.Generated.Format("2006-01-02 15:04:05"),
		summary.Config.URL,
		summary.Config.ConcurrentUsers,
		summary.TotalDuration.String(),
		summary.Config.Method,
		summary.TotalRequests,
		summary.SuccessfulRequests,
		summary.FailedRequests,
		summary.ErrorRate,
		summary.AverageResponseTime.String(),
		summary.RequestsPerSecond,
		r.generateResponseTimeDistributionHTML(summary.ResponseTimeDistribution),
		r.generateStatusCodeDistributionHTML(summary.StatusCodes, summary.TotalRequests),
		summary.MinResponseTime.String(),
		summary.AverageResponseTime.String(),
		summary.MaxResponseTime.String(),
	)
}

// generateResponseTimeDistributionHTML generates HTML for response time distribution
func (r *Reporter) generateResponseTimeDistributionHTML(distribution map[string]int64) string {
	var html string
	total := int64(0)
	for _, count := range distribution {
		total += count
	}

	for timeRange, count := range distribution {
		percentage := float64(count) / float64(total) * 100
		html += fmt.Sprintf(`
			<div style="margin: 10px 0;">
				<div style="display: flex; justify-content: space-between; margin-bottom: 5px;">
					<span>%s</span>
					<span>%d (%.1f%%)</span>
				</div>
				<div class="progress-bar">
					<div class="progress-fill" style="width: %.1f%%;"></div>
				</div>
			</div>`,
			timeRange, count, percentage, percentage)
	}

	return html
}

// generateStatusCodeDistributionHTML generates HTML for status code distribution
func (r *Reporter) generateStatusCodeDistributionHTML(statusCodes map[int]int64, totalRequests int64) string {
	var html string
	for code, count := range statusCodes {
		percentage := float64(count) / float64(totalRequests) * 100
		html += fmt.Sprintf(`
			<tr>
				<td>%d</td>
				<td>%d</td>
				<td>%.1f%%</td>
			</tr>`,
			code, count, percentage)
	}
	return html
}

// generateCSVReport generates a CSV report
func (r *Reporter) generateCSVReport(filename string, results []core.LoadTestResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("RequestID,UserID,StartTime,EndTime,Duration,StatusCode,ResponseSize,Success,Error\n")
	if err != nil {
		return err
	}

	for _, result := range results {
		_, err = file.WriteString(fmt.Sprintf("%s,%d,%s,%s,%s,%d,%d,%t,\"%s\"\n",
			result.RequestID,
			result.UserID,
			result.StartTime.Format(time.RFC3339),
			result.EndTime.Format(time.RFC3339),
			result.Duration.String(),
			result.StatusCode,
			result.ResponseSize,
			result.Success,
			result.Error,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
