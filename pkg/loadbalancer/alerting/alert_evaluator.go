package alerting

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type AlertEvaluatorImpl struct{}

func NewAlertEvaluator() *AlertEvaluatorImpl {
	return &AlertEvaluatorImpl{}
}

func (ae *AlertEvaluatorImpl) EvaluateAlert(alert *core.Alert, metrics *core.RealTimeMetrics) (bool, float64, error) {
	value, err := ae.extractMetricValue(alert.Metric, metrics)
	if err != nil {
		return false, 0, err
	}

	triggered, err := ae.evaluateCondition(value, alert.Operator, alert.Threshold)
	if err != nil {
		return false, 0, err
	}

	return triggered, value, nil
}

func (ae *AlertEvaluatorImpl) extractMetricValue(metric string, metrics *core.RealTimeMetrics) (float64, error) {
	switch strings.ToLower(metric) {
	case "error_rate":
		return metrics.ErrorRate, nil
	case "response_time", "average_response_time":
		return float64(metrics.AverageResponseTime.Nanoseconds()) / 1e6, nil
	case "requests_per_second", "throughput":
		return metrics.RequestsPerSecond, nil
	case "total_requests":
		return float64(metrics.TotalRequests), nil
	case "successful_requests":
		return float64(metrics.SuccessfulRequests), nil
	case "failed_requests":
		return float64(metrics.FailedRequests), nil
	case "percentile_50", "p50":
		return float64(metrics.Percentile50.Nanoseconds()) / 1e6, nil
	case "percentile_95", "p95":
		return float64(metrics.Percentile95.Nanoseconds()) / 1e6, nil
	case "percentile_99", "p99":
		return float64(metrics.Percentile99.Nanoseconds()) / 1e6, nil
	case "bandwidth":
		return metrics.Bandwidth, nil
	case "min_response_time":
		return float64(metrics.MinResponseTime.Nanoseconds()) / 1e6, nil
	case "max_response_time":
		return float64(metrics.MaxResponseTime.Nanoseconds()) / 1e6, nil
	case "standard_deviation":
		return float64(metrics.StandardDeviation.Nanoseconds()) / 1e6, nil
	case "variance":
		return metrics.Variance, nil
	default:
		return 0, fmt.Errorf("unknown metric: %s", metric)
	}
}

func (ae *AlertEvaluatorImpl) evaluateCondition(value float64, operator string, threshold float64) (bool, error) {
	switch operator {
	case ">":
		return value > threshold, nil
	case "<":
		return value < threshold, nil
	case ">=":
		return value >= threshold, nil
	case "<=":
		return value <= threshold, nil
	case "==", "=":
		return value == threshold, nil
	case "!=":
		return value != threshold, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func (ae *AlertEvaluatorImpl) ParseCondition(condition string) (metric, operator string, threshold float64, err error) {
	condition = strings.TrimSpace(condition)
	operators := []string{">=", "<=", "!=", "==", ">", "<", "="}
	var foundOperator string
	var operatorIndex int = -1
	for _, op := range operators {
		if index := strings.Index(condition, op); index != -1 {
			foundOperator = op
			operatorIndex = index
			break
		}
	}
	if operatorIndex == -1 {
		return "", "", 0, fmt.Errorf("no valid operator found in condition: %s", condition)
	}
	metric = strings.TrimSpace(condition[:operatorIndex])
	thresholdStr := strings.TrimSpace(condition[operatorIndex+len(foundOperator):])
	threshold, err = strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid threshold value: %s", thresholdStr)
	}
	if foundOperator == "=" {
		foundOperator = "=="
	}
	return metric, foundOperator, threshold, nil
}

func (ae *AlertEvaluatorImpl) ValidateCondition(condition string) error {
	_, _, _, err := ae.ParseCondition(condition)
	return err
}

func (ae *AlertEvaluatorImpl) GetSupportedMetrics() []string {
	return []string{
		"error_rate",
		"response_time",
		"average_response_time",
		"requests_per_second",
		"throughput",
		"total_requests",
		"successful_requests",
		"failed_requests",
		"percentile_50",
		"p50",
		"percentile_95",
		"p95",
		"percentile_99",
		"p99",
		"bandwidth",
		"min_response_time",
		"max_response_time",
		"standard_deviation",
		"variance",
	}
}

func (ae *AlertEvaluatorImpl) GetSupportedOperators() []string {
	return []string{">", "<", ">=", "<=", "==", "!="}
}
