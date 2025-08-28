package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cherry-pick/pkg/config"
	"github.com/cherry-pick/pkg/intelligence"
)

func main() {
	fmt.Println("=== Environment-Based Database Analysis ===")

	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		fmt.Printf("No database configuration found: %v\n\n", err)
		config.PrintConfigHelp()
		return
	}

	fmt.Printf("Detected database type: %s\n", dbConfig.Type)
	fmt.Printf("Connection URL: %s\n", maskPassword(dbConfig.URL))
	fmt.Println("Connecting to database...")

	service, err := intelligence.CreateSimpleService(dbConfig.Type, dbConfig.URL)
	if err != nil {
		log.Fatalf("Failed to create intelligence service: %v", err)
	}
	defer service.Close()

	fmt.Println("Performing comprehensive database analysis...")
	report, err := service.AnalyzeDatabase()
	if err != nil {
		log.Fatalf("Failed to analyze database: %v", err)
	}

	fmt.Printf("\n=== Analysis Results for: %s ===\n", report.DatabaseName)
	fmt.Printf("Database Type: %s\n", report.DatabaseType)
	fmt.Printf("Analysis Time: %s\n", report.AnalysisTime.Format("2006-01-02 15:04:05"))

	if report.DatabaseType == "mongodb" {
		fmt.Printf("Total Collections: %d\n", report.Summary.TotalTables)
		fmt.Printf("Total Fields: %d\n", report.Summary.TotalColumns)
		fmt.Printf("Total Documents: %d\n", report.Summary.TotalRows)
	} else {
		fmt.Printf("Total Tables: %d\n", report.Summary.TotalTables)
		fmt.Printf("Total Columns: %d\n", report.Summary.TotalColumns)
		fmt.Printf("Total Rows: %d\n", report.Summary.TotalRows)
	}

	fmt.Printf("Total Size: %s\n", report.Summary.TotalSize)
	fmt.Printf("Health Score: %.2f/1.0\n", report.Summary.HealthScore)
	fmt.Printf("Complexity Score: %.2f\n", report.Summary.ComplexityScore)

	if len(report.Tables) > 0 {
		tableWord := "Tables"
		if report.DatabaseType == "mongodb" {
			tableWord = "Collections"
		}

		fmt.Printf("\n=== %s Details (showing first 10) ===\n", tableWord)
		for i, table := range report.Tables {
			if i >= 10 {
				fmt.Printf("... and %d more %s\n", len(report.Tables)-10, strings.ToLower(tableWord))
				break
			}
			rowWord := "rows"
			if report.DatabaseType == "mongodb" {
				rowWord = "documents"
			}
			fmt.Printf("• %s: %d %s, %d columns, %d indexes\n",
				table.Name, table.RowCount, rowWord, len(table.Columns), len(table.Indexes))
		}
	}

	if len(report.Insights) > 0 {
		fmt.Printf("\n=== Key Insights (%d found) ===\n", len(report.Insights))
		for i, insight := range report.Insights {
			if i >= 5 {
				fmt.Printf("... and %d more insights\n", len(report.Insights)-5)
				break
			}
			fmt.Printf("• [%s] %s: %s\n",
				strings.ToUpper(insight.Severity), insight.Title, insight.Description)
		}
	}

	if len(report.Recommendations) > 0 {
		fmt.Printf("\n=== Recommendations (%d found) ===\n", len(report.Recommendations))
		for i, rec := range report.Recommendations {
			if i >= 3 {
				fmt.Printf("... and %d more recommendations\n", len(report.Recommendations)-3)
				break
			}
			fmt.Printf("• %s\n", rec)
		}
	}

	fmt.Println("\nPerforming security analysis...")
	securityIssues, err := service.AnalyzeSecurity()
	if err != nil {
		log.Printf("Security analysis failed: %v", err)
	} else {
		if len(securityIssues) > 0 {
			fmt.Printf("=== Security Issues (%d found) ===\n", len(securityIssues))
			for i, issue := range securityIssues {
				if i >= 3 {
					fmt.Printf("... and %d more security issues\n", len(securityIssues)-3)
					break
				}
				fmt.Printf("• [%s] %s: %s\n",
					strings.ToUpper(issue.Severity), issue.Title, issue.Description)
			}
		} else {
			fmt.Println("✓ No security issues found!")
		}
	}

	fmt.Println("\nTesting query optimization...")
	var testQuery string
	if report.DatabaseType == "mongodb" {
		testQuery = "db.users.find({})"
	} else {
		testQuery = "SELECT * FROM users"
	}

	suggestion, err := service.OptimizeQuery(testQuery)
	if err != nil {
		log.Printf("Query optimization failed: %v", err)
	} else {
		fmt.Printf("Query Optimization Suggestion: %s\n", suggestion.Explanation)
		fmt.Printf("Expected Gain: %s\n", suggestion.ExpectedGain)
	}

	fmt.Println("\nTracking data lineage...")
	lineage, err := service.TrackLineage()
	if err != nil {
		log.Printf("Lineage tracking failed: %v", err)
	} else {
		tableWord := "tables"
		if report.DatabaseType == "mongodb" {
			tableWord = "collections"
		}
		fmt.Printf("✓ Tracked lineage for %d %s\n", len(lineage), tableWord)
	}

	fmt.Println("\nChecking alerts...")
	alerts, err := service.CheckAlerts()
	if err != nil {
		log.Printf("Alert checking failed: %v", err)
	} else {
		if len(alerts) > 0 {
			fmt.Printf("=== Alerts (%d triggered) ===\n", len(alerts))
			for _, alert := range alerts {
				fmt.Printf("• [%s] %s: %s\n",
					strings.ToUpper(alert.Severity), alert.Name, alert.Message)
			}
		} else {
			fmt.Println("✓ No alerts triggered!")
		}
	}

	reportFileName := fmt.Sprintf("data/%s-analysis-report.json", report.DatabaseType)
	fmt.Printf("\nExporting full report to '%s'...\n", reportFileName)
	jsonData, err := service.ExportReport(report, "json")
	if err != nil {
		log.Printf("Failed to export report: %v", err)
	} else {
		err = os.MkdirAll("data", 0755)
		if err != nil {
			log.Printf("Failed to create data directory: %v", err)
		} else {
			err = os.WriteFile(reportFileName, jsonData, 0644)
			if err != nil {
				log.Printf("Failed to write report file: %v", err)
			} else {
				fmt.Printf("✓ Full report saved to '%s'\n", reportFileName)
			}
		}
	}

	fmt.Printf("\n=== %s Analysis Complete ===\n", strings.Title(report.DatabaseType))
	fmt.Println("Environment Configuration Used:")
	fmt.Printf("• Database Type: %s\n", dbConfig.Type)
	fmt.Printf("• Connection: %s\n", maskPassword(dbConfig.URL))
}

func maskPassword(url string) string {
	if strings.Contains(url, "://") && strings.Contains(url, "@") {
		parts := strings.Split(url, "@")
		if len(parts) >= 2 {
			beforeAt := parts[0]
			afterAt := strings.Join(parts[1:], "@")

			if strings.Contains(beforeAt, ":") {
				protocolAndAuth := strings.Split(beforeAt, "://")
				if len(protocolAndAuth) == 2 {
					protocol := protocolAndAuth[0]
					auth := protocolAndAuth[1]
					if strings.Contains(auth, ":") {
						authParts := strings.Split(auth, ":")
						user := authParts[0]
						return fmt.Sprintf("%s://%s:****@%s", protocol, user, afterAt)
					}
				}
			}
		}
	}

	return url
}
