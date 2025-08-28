package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cherry-pick/pkg/intelligence"
)

func main() {
	fmt.Println("=== MongoDB Database Analysis Example ===")

	connectionString := os.Getenv("MONGODB_URL")
	if connectionString == "" {
		fmt.Println("To run this example with a real MongoDB database:")
		fmt.Println("1. Set the MONGODB_URL environment variable")
		fmt.Println("2. Examples:")
		fmt.Println("   Local:     MONGODB_URL=\"mongodb://localhost:27017/db\"")
		fmt.Println("   Atlas:     MONGODB_URL=\"mongodb+srv://user:password@cluster.mongodb.net/db\"")
		fmt.Println("   Auth:      MONGODB_URL=\"mongodb://user:password@localhost:27017/db\"")
		fmt.Println("\n3. Then run: go run examples/mongodb-example.go")
		return
	}

	fmt.Printf("Connecting to MongoDB...\n")

	service, err := intelligence.CreateSimpleService("mongodb", connectionString)
	if err != nil {
		log.Fatalf("Failed to create MongoDB intelligence service: %v", err)
	}
	defer service.Close()

	fmt.Println("Performing comprehensive MongoDB analysis...")
	report, err := service.AnalyzeDatabase()
	if err != nil {
		log.Fatalf("Failed to analyze MongoDB database: %v", err)
	}

	fmt.Printf("\n=== Analysis Results for: %s ===\n", report.DatabaseName)
	fmt.Printf("Database Type: %s\n", report.DatabaseType)
	fmt.Printf("Analysis Time: %s\n", report.AnalysisTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total Collections: %d\n", report.Summary.TotalTables)
	fmt.Printf("Total Fields: %d\n", report.Summary.TotalColumns)
	fmt.Printf("Total Documents: %d\n", report.Summary.TotalRows)
	fmt.Printf("Total Size: %s\n", report.Summary.TotalSize)
	fmt.Printf("Health Score: %.2f/1.0\n", report.Summary.HealthScore)
	fmt.Printf("Complexity Score: %.2f\n", report.Summary.ComplexityScore)

	if len(report.Tables) > 0 {
		fmt.Printf("\n=== Collection Details ===\n")
		for _, table := range report.Tables {
			fmt.Printf("• %s: %d documents, %d fields, %d indexes\n",
				table.Name, table.RowCount, len(table.Columns), len(table.Indexes))
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

	fmt.Println("\nTesting MongoDB query optimization...")
	suggestion, err := service.OptimizeQuery("db.users.find({})")
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
		fmt.Printf("✓ Tracked lineage for %d collections\n", len(lineage))
	}

	fmt.Println("\nChecking alerts...")
	alerts, err := service.CheckAlerts()
	if err != nil {
		log.Printf("Alert checking failed: %v", err)
	} else {
		if len(alerts) > 0 {
			fmt.Printf("=== Alerts (%d triggered) ===\n", len(alerts))
			for _, alert := range alerts {
				fmt.Printf("• [%s] %s: %s\n", alert.Severity, alert.Name, alert.Message)
			}
		} else {
			fmt.Println("✓ No alerts triggered!")
		}
	}

	fmt.Println("\nExporting full report to 'data/mongodb-analysis-report.json'...")
	jsonData, err := service.ExportReport(report, "json")
	if err != nil {
		log.Printf("Failed to export report: %v", err)
	} else {
		err = os.MkdirAll("data", 0755)
		if err != nil {
			log.Printf("Failed to create data directory: %v", err)
		} else {
			err = os.WriteFile("data/mongodb-analysis-report.json", jsonData, 0644)
			if err != nil {
				log.Printf("Failed to write report file: %v", err)
			} else {
				fmt.Println("✓ Full report saved to 'data/mongodb-analysis-report.json'")
			}
		}
	}

	fmt.Println("\n=== MongoDB Analysis Complete ===")
	fmt.Println("Key MongoDB Features Analyzed:")
	fmt.Println("• Collection statistics and document counts")
	fmt.Println("• Schema analysis and field type detection")
	fmt.Println("• Index usage and recommendations")
	fmt.Println("• Security analysis for PII detection")
	fmt.Println("• Performance optimization suggestions")
	fmt.Println("• Data lineage and reference tracking")
}
