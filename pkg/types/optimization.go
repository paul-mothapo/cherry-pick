package types

type OptimizationSuggestion struct {
	OriginalQuery   string   `json:"original_query"`
	OptimizedQuery  string   `json:"optimized_query"`
	Explanation     string   `json:"explanation"`
	ExpectedGain    string   `json:"expected_gain"`
	Confidence      float64  `json:"confidence"`
	AffectedTables  []string `json:"affected_tables"`
}
