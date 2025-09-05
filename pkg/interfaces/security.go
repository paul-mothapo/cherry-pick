package interfaces

import "github.com/cherry-pick/pkg/types"

type SecurityAnalyzer interface {
	AnalyzeSecurity() ([]types.SecurityIssue, error)

	IsPotentialPII(columnName, pattern string) bool

	DetectVulnerabilities(tables []types.TableInfo) []types.SecurityIssue
}
