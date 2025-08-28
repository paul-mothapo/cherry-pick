package optimization

import (
	"strings"

	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
)

type QueryOptimizerImpl struct{}

func NewQueryOptimizer() interfaces.QueryOptimizer {
	return &QueryOptimizerImpl{}
}

func (qo *QueryOptimizerImpl) AnalyzeQuery(query string) (*types.OptimizationSuggestion, error) {
	suggestion := &types.OptimizationSuggestion{
		OriginalQuery: query,
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))

	if strings.HasPrefix(queryLower, "select") && !strings.Contains(queryLower, "where") {
		suggestion.OptimizedQuery = query + "\n-- Add WHERE clause to limit results"
		suggestion.Explanation = "Query lacks WHERE clause which may result in full table scan"
		suggestion.ExpectedGain = "50-90% performance improvement"
		suggestion.Confidence = 0.8
		return suggestion, nil
	}

	if strings.Contains(queryLower, "select *") {
		suggestion.OptimizedQuery = strings.Replace(query, "SELECT *", "SELECT specific_columns", 1)
		suggestion.Explanation = "SELECT * returns all columns, consider specifying only needed columns"
		suggestion.ExpectedGain = "10-30% performance improvement"
		suggestion.Confidence = 0.9
		return suggestion, nil
	}

	if strings.HasPrefix(queryLower, "select") && !strings.Contains(queryLower, "limit") {
		suggestion.OptimizedQuery = query + "\nLIMIT 100 -- Add appropriate limit"
		suggestion.Explanation = "Query lacks LIMIT clause which may return excessive data"
		suggestion.ExpectedGain = "20-70% performance improvement"
		suggestion.Confidence = 0.7
		return suggestion, nil
	}

	if strings.Contains(queryLower, "join") && !strings.Contains(queryLower, "on") {
		suggestion.OptimizedQuery = query + "\n-- Ensure proper JOIN conditions are specified"
		suggestion.Explanation = "JOIN without proper ON conditions may result in cartesian product"
		suggestion.ExpectedGain = "80-95% performance improvement"
		suggestion.Confidence = 0.95
		return suggestion, nil
	}

	suggestion.OptimizedQuery = query
	suggestion.Explanation = "Query appears to be well-structured"
	suggestion.ExpectedGain = "No optimization needed"
	suggestion.Confidence = 0.5

	return suggestion, nil
}

func (qo *QueryOptimizerImpl) ValidateQuery(query string) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return &QueryValidationError{Message: "Query cannot be empty"}
	}

	queryLower := strings.ToLower(query)
	validStarts := []string{"select", "insert", "update", "delete", "create", "alter", "drop"}

	isValid := false
	for _, start := range validStarts {
		if strings.HasPrefix(queryLower, start) {
			isValid = true
			break
		}
	}

	if !isValid {
		return &QueryValidationError{Message: "Query must start with a valid SQL keyword"}
	}

	return nil
}

type QueryValidationError struct {
	Message string
}

func (e *QueryValidationError) Error() string {
	return e.Message
}
