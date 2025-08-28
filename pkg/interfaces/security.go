// Package interfaces defines security analysis interfaces.
package interfaces

import "github.com/cherry-pick/pkg/types"

// SecurityAnalyzer defines the interface for analyzing database security aspects.
type SecurityAnalyzer interface {
	// AnalyzeSecurity performs security analysis on the database.
	AnalyzeSecurity() ([]types.SecurityIssue, error)

	// IsPotentialPII checks if a column might contain personally identifiable information.
	IsPotentialPII(columnName, pattern string) bool

	// DetectVulnerabilities identifies potential security vulnerabilities.
	DetectVulnerabilities(tables []types.TableInfo) []types.SecurityIssue
}
