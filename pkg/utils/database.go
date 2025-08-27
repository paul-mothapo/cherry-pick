package utils

import "strings"

func IsNumericType(dataType string) bool {
	numericTypes := []string{"int", "integer", "bigint", "smallint", "tinyint",
		"decimal", "numeric", "float", "double", "real"}

	dataType = strings.ToLower(dataType)
	for _, numType := range numericTypes {
		if strings.Contains(dataType, numType) {
			return true
		}
	}
	return false
}

func IsStringType(dataType string) bool {
	stringTypes := []string{"varchar", "char", "text", "string"}

	dataType = strings.ToLower(dataType)
	for _, strType := range stringTypes {
		if strings.Contains(dataType, strType) {
			return true
		}
	}
	return false
}

func DetectPattern(samples []string) string {
	if len(samples) == 0 {
		return "No pattern detected"
	}

	for _, sample := range samples {
		if strings.Contains(sample, "@") {
			return "Email pattern"
		}
		if len(sample) >= 10 && IsDigitsOnly(sample) {
			return "Phone number pattern"
		}
		if strings.HasPrefix(sample, "http") {
			return "URL pattern"
		}
	}

	return "Text pattern"
}

func IsDigitsOnly(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
