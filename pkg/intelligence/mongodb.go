package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cherry-pick/pkg/config"
	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
)

type MongoService struct {
	connector interfaces.MongoConnector
	analyzer  interfaces.MongoAnalyzer
	config    interfaces.ConfigManager
}

func NewMongoService(
	connector interfaces.MongoConnector,
	analyzer interfaces.MongoAnalyzer,
	configPath string,
) *Service {
	configManager := config.NewConfigManager()
	if err := configManager.LoadConfig(configPath); err != nil {
		fmt.Printf("Warning: Failed to load configuration: %v\n", err)
	}

	mongoService := &MongoService{
		connector: connector,
		analyzer:  analyzer,
		config:    configManager,
	}

	return &Service{
		connector:    nil,
		analyzer:     nil,
		insights:     nil,
		reporter:     nil,
		security:     nil,
		optimizer:    nil,
		alerts:       nil,
		comparison:   nil,
		lineage:      nil,
		scheduler:    nil,
		config:       configManager,
		performance:  nil,
		mongoService: mongoService,
	}
}

func (ms *MongoService) GetConnector() interfaces.MongoConnector {
	return ms.connector
}

func (ms *MongoService) AnalyzeDatabase(ctx context.Context) (*types.DatabaseReport, error) {
	return ms.analyzer.AnalyzeDatabase(ctx)
}

func (ms *MongoService) AnalyzeSecurity(ctx context.Context) ([]types.SecurityIssue, error) {
	collections, err := ms.analyzer.AnalyzeCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze collections for security: %w", err)
	}

	var issues []types.SecurityIssue

	for _, coll := range collections {
		for _, field := range coll.Fields {
			if ms.isPotentialPII(field.Name) {
				issue := types.SecurityIssue{
					Type:            "privacy",
					Severity:        "high",
					Title:           "Potential PII Data",
					Description:     fmt.Sprintf("Field '%s.%s' may contain personally identifiable information", coll.Name, field.Name),
					Recommendation:  "Consider encryption, field-level security, or access controls",
					AffectedObjects: []string{fmt.Sprintf("%s.%s", coll.Name, field.Name)},
				}
				issues = append(issues, issue)
			}
		}

		if len(coll.Indexes) <= 1 && coll.DocumentCount > 1000 {
			issue := types.SecurityIssue{
				Type:            "performance_security",
				Severity:        "medium",
				Title:           "Unindexed Large Collection",
				Description:     fmt.Sprintf("Collection '%s' has no indexes, making it vulnerable to performance attacks", coll.Name),
				Recommendation:  "Add appropriate indexes to prevent collection scanning attacks",
				AffectedObjects: []string{coll.Name},
			}
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

func (ms *MongoService) OptimizeQuery(query string) (*types.OptimizationSuggestion, error) {
	suggestion := &types.OptimizationSuggestion{
		OriginalQuery: query,
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))

	if strings.Contains(queryLower, "find()") && !strings.Contains(queryLower, "limit(") {
		suggestion.OptimizedQuery = query + ".limit(100)"
		suggestion.Explanation = "Query lacks limit() which may return excessive documents"
		suggestion.ExpectedGain = "20-70% performance improvement"
		suggestion.Confidence = 0.8
		return suggestion, nil
	}

	if strings.Contains(queryLower, "find({})") {
		suggestion.OptimizedQuery = strings.Replace(query, "find({})", "find({specificField: value})", 1)
		suggestion.Explanation = "Empty find query returns all documents, consider adding filter criteria"
		suggestion.ExpectedGain = "50-90% performance improvement"
		suggestion.Confidence = 0.9
		return suggestion, nil
	}

	suggestion.OptimizedQuery = query
	suggestion.Explanation = "MongoDB query appears to be well-structured"
	suggestion.ExpectedGain = "No optimization needed"
	suggestion.Confidence = 0.5

	return suggestion, nil
}

func (ms *MongoService) CheckAlerts(ctx context.Context) ([]types.MonitoringAlert, error) {
	report, err := ms.AnalyzeDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze database for alerts: %w", err)
	}

	var alerts []types.MonitoringAlert

	for _, table := range report.Tables {
		if table.RowCount > 1000000 {
			alert := types.MonitoringAlert{
				ID:          "large_collection_" + table.Name,
				Name:        "Large Collection Alert",
				Condition:   "document_count > 1000000",
				Threshold:   1000000,
				Severity:    "medium",
				Triggered:   true,
				LastTrigger: time.Now(),
				Message:     fmt.Sprintf("Collection %s has %d documents", table.Name, table.RowCount),
			}
			alerts = append(alerts, alert)
		}
	}

	for _, table := range report.Tables {
		if len(table.Indexes) <= 1 && table.RowCount > 10000 {
			alert := types.MonitoringAlert{
				ID:          "missing_indexes_" + table.Name,
				Name:        "Missing Indexes Alert",
				Condition:   "index_count <= 1 AND document_count > 10000",
				Threshold:   10000,
				Severity:    "high",
				Triggered:   true,
				LastTrigger: time.Now(),
				Message:     fmt.Sprintf("Collection %s needs indexes (%d documents, %d indexes)", table.Name, table.RowCount, len(table.Indexes)),
			}
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

func (ms *MongoService) TrackLineage(ctx context.Context) (map[string]types.DataLineage, error) {
	collections, err := ms.analyzer.AnalyzeCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze collections for lineage: %w", err)
	}

	lineage := make(map[string]types.DataLineage)

	for _, coll := range collections {
		collLineage := types.DataLineage{
			TableName:   coll.Name,
			LastUpdated: time.Now(),
		}

		for _, field := range coll.Fields {
			if strings.HasSuffix(field.Name, "_id") && field.Name != "_id" {
				refCollection := strings.TrimSuffix(field.Name, "_id")
				dep := types.LineageDependency{
					TableName: refCollection,
					Type:      "reference",
				}
				collLineage.UpstreamDeps = append(collLineage.UpstreamDeps, dep)
			}
		}

		lineage[coll.Name] = collLineage
	}

	return lineage, nil
}

func (ms *MongoService) ExportReport(report *types.DatabaseReport, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(report, "", "  ")
	case "summary":
		return ms.generateTextSummary(report), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

func (ms *MongoService) GetConfig() *types.Config {
	return ms.config.GetConfig()
}

func (ms *MongoService) UpdateConfig(config *types.Config) error {
	return ms.config.UpdateConfig(config)
}

func (ms *MongoService) Close(ctx context.Context) error {
	return ms.connector.Close(ctx)
}

func (ms *MongoService) isPotentialPII(fieldName string) bool {
	piiPatterns := []string{"email", "phone", "ssn", "social", "address", "name", "firstname", "lastname", "password"}
	fieldLower := strings.ToLower(fieldName)

	for _, pattern := range piiPatterns {
		if strings.Contains(fieldLower, pattern) {
			return true
		}
	}
	return false
}

func (ms *MongoService) generateTextSummary(report *types.DatabaseReport) []byte {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("MongoDB Analysis Report - %s\n", report.DatabaseName))
	summary.WriteString(fmt.Sprintf("Analysis Date: %s\n\n", report.AnalysisTime.Format("2006-01-02 15:04:05")))

	summary.WriteString("SUMMARY\n")
	summary.WriteString("=======\n")
	summary.WriteString(fmt.Sprintf("Collections: %d\n", report.Summary.TotalTables))
	summary.WriteString(fmt.Sprintf("Total Fields: %d\n", report.Summary.TotalColumns))
	summary.WriteString(fmt.Sprintf("Total Documents: %d\n", report.Summary.TotalRows))
	summary.WriteString(fmt.Sprintf("Health Score: %.2f/1.0\n", report.Summary.HealthScore))
	summary.WriteString(fmt.Sprintf("Complexity Score: %.2f\n\n", report.Summary.ComplexityScore))

	summary.WriteString("KEY INSIGHTS\n")
	summary.WriteString("============\n")
	for _, insight := range report.Insights {
		summary.WriteString(fmt.Sprintf("• [%s] %s: %s\n",
			strings.ToUpper(insight.Severity), insight.Title, insight.Description))
	}

	summary.WriteString("\nRECOMMENDATIONS\n")
	summary.WriteString("===============\n")
	for _, rec := range report.Recommendations {
		summary.WriteString(fmt.Sprintf("• %s\n", rec))
	}

	return []byte(summary.String())
}
