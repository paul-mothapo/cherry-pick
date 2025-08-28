// Package interfaces defines query optimization interfaces.
package interfaces

import "github.com/cherry-pick/pkg/types"

// QueryOptimizer defines the interface for query optimization suggestions.
type QueryOptimizer interface {
	// AnalyzeQuery analyzes a query and provides optimization suggestions.
	AnalyzeQuery(query string) (*types.OptimizationSuggestion, error)

	// ValidateQuery validates the syntax and structure of a query.
	ValidateQuery(query string) error
}
