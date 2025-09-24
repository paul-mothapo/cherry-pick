package utils

import (
	"sort"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type MetricsCalculator struct{}

func NewMetricsCalculator() *MetricsCalculator {
	return &MetricsCalculator{}
}

func (mc *MetricsCalculator) CalculateAdvancedMetrics(results []core.LoadTestResult) *core.RealTimeMetrics {
	if len(results) == 0 {
		return &core.RealTimeMetrics{}
	}

	responseTimes := make([]time.Duration, len(results))
	var totalResponseTime time.Duration
	var totalBytes int64
	var successfulRequests int64
	var failedRequests int64

	for i, result := range results {
		responseTimes[i] = result.Duration
		totalResponseTime += result.Duration
		totalBytes += result.ResponseSize

		if result.Success {
			successfulRequests++
		} else {
			failedRequests++
		}
	}

	sort.Slice(responseTimes, func(i, j int) bool {
		return responseTimes[i] < responseTimes[j]
	})

	totalRequests := int64(len(results))
	averageResponseTime := totalResponseTime / time.Duration(len(results))
	errorRate := float64(failedRequests) / float64(totalRequests) * 100

	percentile50 := mc.calculatePercentile(responseTimes, 50)
	percentile95 := mc.calculatePercentile(responseTimes, 95)
	percentile99 := mc.calculatePercentile(responseTimes, 99)

	minResponseTime := responseTimes[0]
	maxResponseTime := responseTimes[len(responseTimes)-1]

	standardDeviation, variance := mc.calculateStandardDeviation(responseTimes, averageResponseTime)

	throughput := mc.calculateThroughput(results)

	bandwidth := mc.calculateBandwidth(results)

	return &core.RealTimeMetrics{
		TestID:              results[0].RequestID[:8],
		Timestamp:           time.Now(),
		ActiveUsers:         0,
		RequestsPerSecond:   throughput,
		AverageResponseTime: averageResponseTime,
		ErrorRate:           errorRate,
		TotalRequests:       totalRequests,
		SuccessfulRequests:  successfulRequests,
		FailedRequests:      failedRequests,
		Percentile50:        percentile50,
		Percentile95:        percentile95,
		Percentile99:        percentile99,
		Throughput:          throughput,
		Bandwidth:           bandwidth,
		MinResponseTime:     minResponseTime,
		MaxResponseTime:     maxResponseTime,
		StandardDeviation:   standardDeviation,
		Variance:            variance,
	}
}

func (mc *MetricsCalculator) calculatePercentile(responseTimes []time.Duration, percentile int) time.Duration {
	if len(responseTimes) == 0 {
		return 0
	}

	index := int(float64(len(responseTimes)) * float64(percentile) / 100.0)
	if index >= len(responseTimes) {
		index = len(responseTimes) - 1
	}

	return responseTimes[index]
}

func (mc *MetricsCalculator) calculateStandardDeviation(responseTimes []time.Duration, mean time.Duration) (time.Duration, float64) {
	if len(responseTimes) == 0 {
		return 0, 0
	}

	var sumSquaredDiffs float64
	meanFloat := float64(mean.Nanoseconds())

	for _, rt := range responseTimes {
		diff := float64(rt.Nanoseconds()) - meanFloat
		sumSquaredDiffs += diff * diff
	}

	variance := sumSquaredDiffs / float64(len(responseTimes))
	standardDeviation := time.Duration(int64(variance))

	return standardDeviation, variance
}

func (mc *MetricsCalculator) calculateThroughput(results []core.LoadTestResult) float64 {
	if len(results) == 0 {
		return 0
	}

	startTime := results[0].StartTime
	endTime := results[0].EndTime

	for _, result := range results {
		if result.StartTime.Before(startTime) {
			startTime = result.StartTime
		}
		if result.EndTime.After(endTime) {
			endTime = result.EndTime
		}
	}

	duration := endTime.Sub(startTime)
	if duration.Seconds() == 0 {
		return 0
	}

	return float64(len(results)) / duration.Seconds()
}

func (mc *MetricsCalculator) calculateBandwidth(results []core.LoadTestResult) float64 {
	if len(results) == 0 {
		return 0
	}

	var totalBytes int64
	startTime := results[0].StartTime
	endTime := results[0].EndTime

	for _, result := range results {
		totalBytes += result.ResponseSize
		if result.StartTime.Before(startTime) {
			startTime = result.StartTime
		}
		if result.EndTime.After(endTime) {
			endTime = result.EndTime
		}
	}

	duration := endTime.Sub(startTime)
	if duration.Seconds() == 0 {
		return 0
	}

	return float64(totalBytes) / duration.Seconds()
}

func (mc *MetricsCalculator) CalculateResponseTimeDistribution(results []core.LoadTestResult) map[string]int64 {
	distribution := map[string]int64{
		"<100ms":    0,
		"100-500ms": 0,
		"500ms-1s":  0,
		"1-2s":      0,
		">2s":       0,
	}

	for _, result := range results {
		duration := result.Duration
		switch {
		case duration < 100*time.Millisecond:
			distribution["<100ms"]++
		case duration < 500*time.Millisecond:
			distribution["100-500ms"]++
		case duration < 1000*time.Millisecond:
			distribution["500ms-1s"]++
		case duration < 2000*time.Millisecond:
			distribution["1-2s"]++
		default:
			distribution[">2s"]++
		}
	}

	return distribution
}

func (mc *MetricsCalculator) CalculateStatusCodesDistribution(results []core.LoadTestResult) map[int]int64 {
	distribution := make(map[int]int64)

	for _, result := range results {
		distribution[result.StatusCode]++
	}

	return distribution
}
