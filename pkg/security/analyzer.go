// Package security provides database security analysis functionality.
package security

import (
	"fmt"
	"strings"

	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
)

// SecurityAnalyzerImpl implements the SecurityAnalyzer interface.
type SecurityAnalyzerImpl struct {
	analyzer interfaces.DatabaseAnalyzer
}

// NewSecurityAnalyzer creates a new security analyzer instance.
func NewSecurityAnalyzer(analyzer interfaces.DatabaseAnalyzer) interfaces.SecurityAnalyzer {
	return &SecurityAnalyzerImpl{
		analyzer: analyzer,
	}
}

// AnalyzeSecurity performs security analysis on the database.
func (sa *SecurityAnalyzerImpl) AnalyzeSecurity() ([]types.SecurityIssue, error) {
	var issues []types.SecurityIssue

	// Get all tables for analysis
	tables, err := sa.analyzer.AnalyzeTables()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze tables: %w", err)
	}

	// Check for potential PII columns
	issues = append(issues, sa.checkPIIColumns(tables)...)

	// Check for unindexed large tables (potential for data exposure)
	issues = append(issues, sa.checkUnindexedTables(tables)...)

	// Detect other vulnerabilities
	issues = append(issues, sa.DetectVulnerabilities(tables)...)

	return issues, nil
}

// checkPIIColumns identifies columns that may contain personally identifiable information.
func (sa *SecurityAnalyzerImpl) checkPIIColumns(tables []types.TableInfo) []types.SecurityIssue {
	var issues []types.SecurityIssue

	for _, table := range tables {
		for _, column := range table.Columns {
			if sa.IsPotentialPII(column.Name, column.DataProfile.Pattern) {
				issue := types.SecurityIssue{
					Type:     "privacy",
					Severity: "high",
					Title:    "Potential PII Data",
					Description: fmt.Sprintf("Column '%s.%s' may contain personally identifiable information",
						table.Name, column.Name),
					Recommendation:  "Consider encryption, masking, or access controls",
					AffectedObjects: []string{fmt.Sprintf("%s.%s", table.Name, column.Name)},
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// checkUnindexedTables identifies large tables without proper indexing.
func (sa *SecurityAnalyzerImpl) checkUnindexedTables(tables []types.TableInfo) []types.SecurityIssue {
	var issues []types.SecurityIssue

	for _, table := range tables {
		if table.RowCount > 1000 && len(table.Indexes) == 0 {
			issue := types.SecurityIssue{
				Type:     "performance_security",
				Severity: "medium",
				Title:    "Unindexed Large Table",
				Description: fmt.Sprintf("Table '%s' has no indexes, making it vulnerable to performance attacks",
					table.Name),
				Recommendation:  "Add appropriate indexes to prevent table scanning attacks",
				AffectedObjects: []string{table.Name},
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// IsPotentialPII checks if a column might contain personally identifiable information.
func (sa *SecurityAnalyzerImpl) IsPotentialPII(columnName, pattern string) bool {
	piiPatterns := []string{"email", "phone", "ssn", "social", "address", "name", "first_name", "last_name", "password"}
	columnLower := strings.ToLower(columnName)

	for _, piiPattern := range piiPatterns {
		if strings.Contains(columnLower, piiPattern) {
			return true
		}
	}

	return pattern == "Email pattern" || pattern == "Phone number pattern"
}

// DetectVulnerabilities identifies potential security vulnerabilities.
func (sa *SecurityAnalyzerImpl) DetectVulnerabilities(tables []types.TableInfo) []types.SecurityIssue {
	var issues []types.SecurityIssue

	// Check for tables with weak naming conventions that might expose sensitive data
	for _, table := range tables {
		tableLower := strings.ToLower(table.Name)
		sensitiveTableNames := []string{"user", "account", "payment", "credit", "admin", "password"}

		for _, sensitive := range sensitiveTableNames {
			if strings.Contains(tableLower, sensitive) && len(table.Indexes) == 0 {
				issue := types.SecurityIssue{
					Type:     "access_control",
					Severity: "high",
					Title:    "Sensitive Table Without Proper Indexing",
					Description: fmt.Sprintf("Table '%s' appears to contain sensitive data but lacks proper indexing",
						table.Name),
					Recommendation:  "Implement proper indexing and access controls for sensitive data",
					AffectedObjects: []string{table.Name},
				}
				issues = append(issues, issue)
				break
			}
		}

		// Check for columns that might store passwords in plain text
		for _, column := range table.Columns {
			if sa.isPotentialPasswordColumn(column) {
				issue := types.SecurityIssue{
					Type:     "encryption",
					Severity: "critical",
					Title:    "Potential Plain Text Password Storage",
					Description: fmt.Sprintf("Column '%s.%s' might be storing passwords in plain text",
						table.Name, column.Name),
					Recommendation:  "Ensure passwords are properly hashed and salted",
					AffectedObjects: []string{fmt.Sprintf("%s.%s", table.Name, column.Name)},
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// isPotentialPasswordColumn checks if a column might be storing passwords.
func (sa *SecurityAnalyzerImpl) isPotentialPasswordColumn(column types.ColumnInfo) bool {
	columnLower := strings.ToLower(column.Name)
	passwordIndicators := []string{"password", "passwd", "pwd", "pass"}

	for _, indicator := range passwordIndicators {
		if strings.Contains(columnLower, indicator) {
			// Check if it's likely plain text (too short for hashed passwords)
			if column.MaxLength > 0 && column.MaxLength < 32 {
				return true
			}
		}
	}

	return false
}
