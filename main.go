package main

import (
	"fmt"
	"strings"

	"github.com/cherry-pick/pkg/intelligence"
)

func NewDatabaseIntelligence(driverName, dataSourceName string) (*intelligence.Service, error) {
	return intelligence.CreateSimpleService(driverName, dataSourceName)
}

func main() {
	fmt.Println("Enterprise Database Intelligence System")
	fmt.Println("=====================================")

	// Example connection strings - replace with your actual database
	// MySQL: "user:password@tcp(localhost:3306)/dbname"
	// PostgreSQL: "postgres://user:password@localhost/dbname?sslmode=disable"
	// SQLite: "./database.db"
	// MongoDB: "mongodb://localhost:27017/dbname" or "mongodb+srv://user:password@cluster.mongodb.net/dbname"

	// Uncomment and modify the following lines to use with your database:
	/*
		service, err := NewDatabaseIntelligence("mysql", "user:password@tcp(localhost:3306)/dbname")
		if err != nil {
			log.Fatal("Failed to create database intelligence service:", err)
		}
		defer service.Close()

		// Perform comprehensive analysis
		report, err := service.AnalyzeDatabase()
		if err != nil {
			log.Fatal("Failed to analyze database:", err)
		}

		// Convert to JSON for easy viewing
		jsonReport, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			log.Fatal("Failed to marshal report:", err)
		}

		fmt.Println(string(jsonReport))

		// Example of other features

		// Security analysis
		securityIssues, err := service.AnalyzeSecurity()
		if err != nil {
			log.Printf("Security analysis failed: %v", err)
		} else {
			fmt.Printf("Found %d security issues\n", len(securityIssues))
		}

		// Query optimization
		suggestion, err := service.OptimizeQuery("SELECT * FROM users")
		if err != nil {
			log.Printf("Query optimization failed: %v", err)
			} else {
			fmt.Printf("Query optimization suggestion: %s\n", suggestion.Explanation)
		}

		// Check alerts
		alerts, err := service.CheckAlerts()
			if err != nil {
			log.Printf("Alert checking failed: %v", err)
		} else {
			fmt.Printf("Found %d triggered alerts\n", len(alerts))
		}

		// Data lineage tracking
		lineage, err := service.TrackLineage()
			if err != nil {
			log.Printf("Lineage tracking failed: %v", err)
			} else {
			fmt.Printf("Tracked lineage for %d tables\n", len(lineage))
		}

		// Schedule periodic analysis (commented out for demo)
		// service.ScheduleAnalysis(24*time.Hour, func(report *types.DatabaseReport) {
		//     log.Printf("Scheduled analysis completed at %s", report.AnalysisTime)
		// })
	*/

	// For demonstration, show the new architecture
	demonstrateNewArchitecture()
}

func demonstrateNewArchitecture() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("NEW MODULAR ARCHITECTURE")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n1. BASIC USAGE:")
	fmt.Println("   service, err := intelligence.CreateSimpleService(\"mysql\", connectionString)")
	fmt.Println("   report, err := service.AnalyzeDatabase()")

	fmt.Println("\n2. ADVANCED USAGE WITH CONFIGURATION:")
	fmt.Println("   service, err := intelligence.CreateConfiguredService(\"mysql\", connectionString, \"config.json\")")

	fmt.Println("\n3. CUSTOM SERVICE BUILDING:")
	fmt.Println("   service, err := intelligence.NewServiceBuilder(\"mysql\", connectionString)")
	fmt.Println("       .WithConfig(\"config.json\")")
	fmt.Println("       .Build()")

	fmt.Println("\n4. MONITORING SETUP:")
	fmt.Println("   service.ScheduleAnalysis(24*time.Hour, func(report *types.DatabaseReport) {")
	fmt.Println("       alerts := service.CheckAlerts()")
	fmt.Println("       // Handle triggered alerts")
	fmt.Println("   })")

	fmt.Println("\n5. SECURITY ANALYSIS:")
	fmt.Println("   issues, err := service.AnalyzeSecurity()")

	fmt.Println("\n6. QUERY OPTIMIZATION:")
	fmt.Println("   suggestion, err := service.OptimizeQuery(\"SELECT * FROM users\")")

	fmt.Println("\n7. DATA LINEAGE:")
	fmt.Println("   lineage, err := service.TrackLineage()")

	fmt.Println("\n8. REPORT COMPARISON:")
	fmt.Println("   comparison := service.CompareReports(oldReport, newReport)")

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("KEY IMPROVEMENTS:")
	fmt.Println("• Proper separation of concerns")
	fmt.Println("• Interface-based design for testability")
	fmt.Println("• Dependency injection")
	fmt.Println("• Factory pattern for easy instantiation")
	fmt.Println("• Modular package structure")
	fmt.Println("• Configuration management")
	fmt.Println("• Error handling best practices")
	fmt.Println("")
	fmt.Println("SUPPORTED DATABASES: MySQL, PostgreSQL, SQLite, MongoDB")
	fmt.Println("FEATURES: Analysis, Monitoring, Optimization, Security, Lineage")
	fmt.Println(strings.Repeat("=", 60))
}
