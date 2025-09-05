package interfaces

import "github.com/cherry-pick/pkg/types"

type QueryOptimizer interface {
	AnalyzeQuery(query string) (*types.OptimizationSuggestion, error)

	ValidateQuery(query string) error
}
