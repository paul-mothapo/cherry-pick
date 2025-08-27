package main

import (
	"fmt"
	"log"
	"os"

	"github.com/intelligent-algorithm/pkg/intelligence"
)

func main() {
	fmt.Println("=== Real Database Analysis Example ===")

	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		fmt.Println("To run this example with a real database:")
		fmt.Println("1. Set the DATABASE_URL environment variable")
		fmt.Println("2. Examples:")
		fmt.Println("   MySQL:      DATABASE_URL=\"user:password@tcp(localhost:3306)/dbname\"")
		fmt.Println("   PostgreSQL: DATABASE_URL=\"postgres://user:password@localhost/dbname?sslmode=disable\"")
		fmt.Println("   SQLite:     DATABASE_URL=\"./database.db\"")
		fmt.Println("\n3. Then run: go run examples/real-database-example.go")
		return
	}

	var dbType string
	switch {
	case contains(connectionString, "mysql") || contains(connectionString, "@tcp"):
		dbType = "mysql"
	case contains(connectionString, "postgres://"):
		dbType = "postgres"
	case contains(connectionString, ".db"):
		dbType = "sqlite3"
	default:
		log.Fatal("Could not determine database type from connection string")
	}

	fmt.Printf("Connecting to %s database...\n", dbType)

	service, err := intelligence.CreateSimpleService(dbType, connectionString)
	if err != nil {
		log.Fatalf("Failed to create database intelligence service: %v", err)
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
	fmt.Printf("Total Tables: %d\n", report.Summary.TotalTables)
	fmt.Printf("Total Columns: %d\n", report.Summary.TotalColumns)
	fmt.Printf("Total Rows: %d\n", report.Summary.TotalRows)
	fmt.Printf("Health Score: %.2f/1.0\n", report.Summary.HealthScore)
	fmt.Printf("Complexity Score: %.2f\n", report.Summary.ComplexityScore)

	if len(report.Insights) > 0 {
		fmt.Printf("\n=== Key Insights (%d found) ===\n", len(report.Insights))
		for i, insight := range report.Insights {
			if i >= 5 {
				fmt.Printf("... and %d more insights\n", len(report.Insights)-5)
				break
			}
			fmt.Printf("• [%s] %s: %s\n",
				insight.Severity, insight.Title, insight.Description)
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
					issue.Severity, issue.Title, issue.Description)
			}
		} else {
			fmt.Println("✓ No security issues found!")
		}
	}

	fmt.Println("\nTesting query optimization...")
	suggestion, err := service.OptimizeQuery("SELECT * FROM users")
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
		fmt.Printf("✓ Tracked lineage for %d tables\n", len(lineage))
	}

	fmt.Println("\nExporting full report to 'data/database-analysis-report.json'...")
	jsonData, err := service.ExportReport(report, "json")
	if err != nil {
		log.Printf("Failed to export report: %v", err)
	} else {
		err = os.MkdirAll("data", 0755)
		if err != nil {
			log.Printf("Failed to create data directory: %v", err)
		} else {
			err = os.WriteFile("data/database-analysis-report.json", jsonData, 0644)
			if err != nil {
				log.Printf("Failed to write report file: %v", err)
			} else {
				fmt.Println("✓ Full report saved to 'data/database-analysis-report.json'")
			}
		}
	}

	fmt.Println("\n=== Analysis Complete ===")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
