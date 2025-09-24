package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type Engine struct {
	client    HTTPClient
	results   map[string][]core.LoadTestResult
	summaries map[string]*core.LoadTestSummary
	statuses  map[string]*core.LoadTestStatus
	mu        sync.RWMutex
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewEngine() *Engine {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &Engine{
		client:    client,
		results:   make(map[string][]core.LoadTestResult),
		summaries: make(map[string]*core.LoadTestSummary),
		statuses:  make(map[string]*core.LoadTestStatus),
	}
}

func (e *Engine) StartLoadTest(testID string, config core.LoadTestConfig) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.statuses[testID]; exists {
		return fmt.Errorf("test with ID %s already exists", testID)
	}

	e.statuses[testID] = &core.LoadTestStatus{
		TestID:   testID,
		Status:   "pending",
		Progress: 0.0,
	}

	go e.runLoadTest(testID, config)
	return nil
}

func (e *Engine) runLoadTest(testID string, config core.LoadTestConfig) {
	e.mu.Lock()
	e.statuses[testID].Status = "running"
	e.statuses[testID].StartTime = time.Now()
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.statuses[testID].Status = "completed"
		e.statuses[testID].EndTime = time.Now()
		e.statuses[testID].Progress = 1.0
		e.mu.Unlock()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	resultsChan := make(chan core.LoadTestResult, config.ConcurrentUsers*10)
	done := make(chan bool)

	go e.collectResults(testID, resultsChan, done)

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			e.runUser(ctx, userID, config, resultsChan)
		}(i)
	}

	wg.Wait()
	close(resultsChan)
	<-done

	e.generateSummary(testID, config, startTime)
}

func (e *Engine) runUser(ctx context.Context, userID int, config core.LoadTestConfig, resultsChan chan<- core.LoadTestResult) {
	ticker := time.NewTicker(config.RequestDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result := e.makeRequest(userID, config)
			select {
			case resultsChan <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (e *Engine) makeRequest(userID int, config core.LoadTestConfig) core.LoadTestResult {
	startTime := time.Now()
	result := core.LoadTestResult{
		RequestID: fmt.Sprintf("%d-%d", userID, startTime.UnixNano()),
		UserID:    userID,
		StartTime: startTime,
	}

	req, err := http.NewRequest(config.Method, config.URL, nil)
	if err != nil {
		result.Error = err.Error()
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Success = false
		return result
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		result.Error = err.Error()
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Success = false
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Success = false
		return result
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.StatusCode = resp.StatusCode
	result.ResponseSize = int64(len(body))
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	return result
}

func (e *Engine) collectResults(testID string, resultsChan <-chan core.LoadTestResult, done chan<- bool) {
	var results []core.LoadTestResult

	for result := range resultsChan {
		results = append(results, result)

		e.mu.Lock()
		if status, exists := e.statuses[testID]; exists {
			if status.StartTime.IsZero() {
				status.Progress = 0.0
			} else {
				elapsed := time.Since(status.StartTime)
				status.Progress = float64(elapsed) / float64(30*time.Second)
				if status.Progress > 1.0 {
					status.Progress = 1.0
				}
			}
		}
		e.mu.Unlock()
	}

	e.mu.Lock()
	e.results[testID] = results
	e.mu.Unlock()

	done <- true
}

func (e *Engine) generateSummary(testID string, config core.LoadTestConfig, startTime time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()

	results := e.results[testID]
	if len(results) == 0 {
		return
	}

	summary := &core.LoadTestSummary{
		TestID:                   testID,
		Config:                   config,
		StartTime:                startTime,
		EndTime:                  time.Now(),
		TotalDuration:            time.Since(startTime),
		TotalRequests:            int64(len(results)),
		StatusCodes:              make(map[int]int64),
		ResponseTimeDistribution: make(map[string]int64),
	}

	var totalResponseTime time.Duration
	var minResponseTime, maxResponseTime time.Duration
	var successfulRequests int64

	for i, result := range results {
		if result.Success {
			successfulRequests++
		}

		summary.StatusCodes[result.StatusCode]++

		totalResponseTime += result.Duration
		if i == 0 || result.Duration < minResponseTime {
			minResponseTime = result.Duration
		}
		if result.Duration > maxResponseTime {
			maxResponseTime = result.Duration
		}

		switch {
		case result.Duration < 100*time.Millisecond:
			summary.ResponseTimeDistribution["<100ms"]++
		case result.Duration < 500*time.Millisecond:
			summary.ResponseTimeDistribution["100-500ms"]++
		case result.Duration < 1000*time.Millisecond:
			summary.ResponseTimeDistribution["500ms-1s"]++
		case result.Duration < 2000*time.Millisecond:
			summary.ResponseTimeDistribution["1-2s"]++
		default:
			summary.ResponseTimeDistribution[">2s"]++
		}
	}

	summary.SuccessfulRequests = successfulRequests
	summary.FailedRequests = summary.TotalRequests - successfulRequests
	summary.AverageResponseTime = totalResponseTime / time.Duration(len(results))
	summary.MinResponseTime = minResponseTime
	summary.MaxResponseTime = maxResponseTime
	summary.RequestsPerSecond = float64(summary.TotalRequests) / summary.TotalDuration.Seconds()
	summary.ErrorRate = float64(summary.FailedRequests) / float64(summary.TotalRequests) * 100

	e.summaries[testID] = summary
}

func (e *Engine) GetTestStatus(testID string) (*core.LoadTestStatus, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	status, exists := e.statuses[testID]
	if !exists {
		return nil, fmt.Errorf("test with ID %s not found", testID)
	}

	return status, nil
}

func (e *Engine) GetTestSummary(testID string) (*core.LoadTestSummary, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	summary, exists := e.summaries[testID]
	if !exists {
		return nil, fmt.Errorf("test summary with ID %s not found", testID)
	}

	return summary, nil
}

func (e *Engine) GetTestResults(testID string) ([]core.LoadTestResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	results, exists := e.results[testID]
	if !exists {
		return nil, fmt.Errorf("test results with ID %s not found", testID)
	}

	return results, nil
}

func (e *Engine) CancelTest(testID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	status, exists := e.statuses[testID]
	if !exists {
		return fmt.Errorf("test with ID %s not found", testID)
	}

	if status.Status != "running" {
		return fmt.Errorf("test with ID %s is not running", testID)
	}

	status.Status = "cancelled"
	status.EndTime = time.Now()
	status.Progress = 1.0

	return nil
}

func (e *Engine) GetAllTests() map[string]*core.LoadTestStatus {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tests := make(map[string]*core.LoadTestStatus)
	for id, status := range e.statuses {
		tests[id] = status
	}

	return tests
}

func (e *Engine) GetRealTimeMetrics(testID string) (*core.RealTimeMetrics, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	status, exists := e.statuses[testID]
	if !exists {
		return nil, fmt.Errorf("test with ID %s not found", testID)
	}

	if status.Status != "running" {
		return nil, fmt.Errorf("test with ID %s is not running", testID)
	}

	results := e.results[testID]
	if len(results) == 0 {
		return &core.RealTimeMetrics{
			TestID:    testID,
			Timestamp: time.Now(),
		}, nil
	}

	cutoff := time.Now().Add(-10 * time.Second)
	var recentResults []core.LoadTestResult
	var totalResponseTime time.Duration
	var successfulRequests int64

	for _, result := range results {
		if result.StartTime.After(cutoff) {
			recentResults = append(recentResults, result)
			totalResponseTime += result.Duration
			if result.Success {
				successfulRequests++
			}
		}
	}

	metrics := &core.RealTimeMetrics{
		TestID:             testID,
		Timestamp:          time.Now(),
		ActiveUsers:        int(status.Progress * 100),
		TotalRequests:      int64(len(results)),
		SuccessfulRequests: successfulRequests,
		FailedRequests:     int64(len(results)) - successfulRequests,
	}

	if len(recentResults) > 0 {
		metrics.RequestsPerSecond = float64(len(recentResults)) / 10.0
		metrics.AverageResponseTime = totalResponseTime / time.Duration(len(recentResults))
		metrics.ErrorRate = float64(len(recentResults)-int(successfulRequests)) / float64(len(recentResults)) * 100
	}

	return metrics, nil
}

func (e *Engine) CleanupOldTests(olderThan time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()

	cutoff := time.Now().Add(-olderThan)

	for testID, status := range e.statuses {
		if status.Status == "completed" || status.Status == "failed" || status.Status == "cancelled" {
			if status.EndTime.Before(cutoff) {
				delete(e.statuses, testID)
				delete(e.results, testID)
				delete(e.summaries, testID)
			}
		}
	}
}
