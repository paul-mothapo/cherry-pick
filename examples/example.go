package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enterprise Database Intelligence System - Examples")
	fmt.Println("================================================")

	basicUsageExample()

	advancedConfigExample()

	monitoringExample()

	securityExample()

	customServiceExample()
}

func basicUsageExample() {
	fmt.Println("\n1. BASIC USAGE EXAMPLE:")
	fmt.Println("   // Create service with minimal configuration")
	fmt.Println("   service, err := intelligence.CreateSimpleService(\"mysql\", \"user:password@tcp(localhost:3306)/dbname\")")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatal(err)")
	fmt.Println("   }")
	fmt.Println("   defer service.Close()")
	fmt.Println("")
	fmt.Println("   // Analyze database")
	fmt.Println("   report, err := service.AnalyzeDatabase()")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatal(err)")
	fmt.Println("   }")
	fmt.Println("")
	fmt.Println("   // Export report")
	fmt.Println("   jsonData, _ := service.ExportReport(report, \"json\")")
	fmt.Println("   fmt.Println(string(jsonData))")
}

func advancedConfigExample() {
	fmt.Println("\n2. ADVANCED CONFIGURATION EXAMPLE:")
	fmt.Println("   // Create service with custom configuration")
	fmt.Println("   service, err := intelligence.CreateConfiguredService(")
	fmt.Println("       \"postgres\",")
	fmt.Println("       \"postgres://user:password@localhost/dbname?sslmode=disable\",")
	fmt.Println("       \"config.json\")")
	fmt.Println("")
	fmt.Println("   // Get current configuration")
	fmt.Println("   config := service.GetConfig()")
	fmt.Println("   config.AnalysisSettings.LargeTableThreshold = 500000")
	fmt.Println("   service.UpdateConfig(config)")
}

func monitoringExample() {
	fmt.Println("\n3. MONITORING SETUP EXAMPLE:")
	fmt.Println("   // Setup automated monitoring")
	fmt.Println("   service.ScheduleAnalysis(24*time.Hour, func(report *types.DatabaseReport) {")
	fmt.Println("       log.Printf(\"Analysis completed: Health Score %.2f\", report.Summary.HealthScore)")
	fmt.Println("")
	fmt.Println("       // Check for alerts")
	fmt.Println("       alerts, _ := service.CheckAlerts()")
	fmt.Println("       for _, alert := range alerts {")
	fmt.Println("           if alert.Severity == \"high\" {")
	fmt.Println("               log.Printf(\"HIGH ALERT: %s\", alert.Message)")
	fmt.Println("               // Send notification")
	fmt.Println("           }")
	fmt.Println("       }")
	fmt.Println("   })")
}

func securityExample() {
	fmt.Println("\n4. SECURITY ANALYSIS EXAMPLE:")
	fmt.Println("   // Perform security analysis")
	fmt.Println("   issues, err := service.AnalyzeSecurity()")
	fmt.Println("   if err != nil {")
	fmt.Println("       log.Fatal(err)")
	fmt.Println("   }")
	fmt.Println("")
	fmt.Println("   // Process security issues")
	fmt.Println("   for _, issue := range issues {")
	fmt.Println("       switch issue.Severity {")
	fmt.Println("       case \"critical\":")
	fmt.Println("           log.Printf(\"CRITICAL: %s - %s\", issue.Title, issue.Description)")
	fmt.Println("       case \"high\":")
	fmt.Println("           log.Printf(\"HIGH: %s - %s\", issue.Title, issue.Recommendation)")
	fmt.Println("       }")
	fmt.Println("   }")
}

func customServiceExample() {
	fmt.Println("\n5. CUSTOM SERVICE BUILDING EXAMPLE:")
	fmt.Println("   // Build service with custom components")
	fmt.Println("   service, err := intelligence.NewServiceBuilder(\"sqlite3\", \"./database.db\")")
	fmt.Println("       .WithConfig(\"custom-config.json\")")
	fmt.Println("       .Build()")
	fmt.Println("")
	fmt.Println("   // Use all features")
	fmt.Println("   report, _ := service.AnalyzeDatabase()")
	fmt.Println("   lineage, _ := service.TrackLineage()")
	fmt.Println("   suggestion, _ := service.OptimizeQuery(\"SELECT * FROM users WHERE active = 1\")")
	fmt.Println("")
	fmt.Println("   // Compare reports over time")
	fmt.Println("   // oldReport := loadPreviousReport()")
	fmt.Println("   // comparison := service.CompareReports(oldReport, report)")
	fmt.Println("   // analyzeChanges(comparison)")
}
