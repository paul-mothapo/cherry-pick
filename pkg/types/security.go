package types

type SecurityIssue struct {
	Type            string   `json:"type"`
	Severity        string   `json:"severity"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Recommendation  string   `json:"recommendation"`
	AffectedObjects []string `json:"affected_objects"`
}
