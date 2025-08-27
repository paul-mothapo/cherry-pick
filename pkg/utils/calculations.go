package utils

import "math"

func CalculateImpact(changeType string, oldValue, newValue interface{}) string {
	switch changeType {
	case "table_count":
		oldCount, okOld := oldValue.(int)
		newCount, okNew := newValue.(int)
		if !okOld || !okNew || oldCount == 0 {
			return "medium"
		}
		change := float64(newCount-oldCount) / float64(oldCount)
		if math.Abs(change) > 0.5 {
			return "high"
		} else if math.Abs(change) > 0.2 {
			return "medium"
		}
		return "low"
	default:
		return "medium"
	}
}

func CalculateRowCountImpact(oldCount, newCount int64) string {
	if oldCount == 0 {
		return "high"
	}

	change := float64(newCount-oldCount) / float64(oldCount)
	absChange := math.Abs(change)

	if absChange > 1.0 {
		return "high"
	} else if absChange > 0.3 {
		return "medium"
	}
	return "low"
}

func CalculatePercentageChange(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		if newValue == 0 {
			return 0
		}
		return 100
	}
	return ((newValue - oldValue) / oldValue) * 100
}
