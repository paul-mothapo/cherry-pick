package intelligence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
)

type Service struct {
	connector    interfaces.DatabaseConnector
	analyzer     interfaces.DatabaseAnalyzer
	insights     interfaces.InsightGenerator
	reporter     interfaces.ReportGenerator
	security     interfaces.SecurityAnalyzer
	optimizer    interfaces.QueryOptimizer
	alerts       interfaces.AlertManager
	comparison   interfaces.ComparisonEngine
	lineage      interfaces.DataLineageTracker
	scheduler    interfaces.Scheduler
	config       interfaces.ConfigManager
	performance  interfaces.PerformanceAnalyzer
	mongoService *MongoService
}

func NewService(
	connector interfaces.DatabaseConnector,
	analyzer interfaces.DatabaseAnalyzer,
	insights interfaces.InsightGenerator,
	reporter interfaces.ReportGenerator,
	security interfaces.SecurityAnalyzer,
	optimizer interfaces.QueryOptimizer,
	alerts interfaces.AlertManager,
	comparison interfaces.ComparisonEngine,
	lineage interfaces.DataLineageTracker,
	scheduler interfaces.Scheduler,
	config interfaces.ConfigManager,
	performance interfaces.PerformanceAnalyzer,
) *Service {
	return &Service{
		connector:   connector,
		analyzer:    analyzer,
		insights:    insights,
		reporter:    reporter,
		security:    security,
		optimizer:   optimizer,
		alerts:      alerts,
		comparison:  comparison,
		lineage:     lineage,
		scheduler:   scheduler,
		config:      config,
		performance: performance,
	}
}

func (s *Service) AnalyzeDatabase() (*types.DatabaseReport, error) {
	if s.mongoService != nil {
		ctx := context.Background()
		return s.mongoService.AnalyzeDatabase(ctx)
	}

	log.Println("Starting comprehensive database analysis...")

	dbName, err := s.connector.GetDatabaseName()
	if err != nil {
		log.Printf("Warning: Could not determine database name: %v", err)
		dbName = "Unknown"
	}

	dbType := s.connector.GetDatabaseType()

	tables, err := s.analyzer.AnalyzeTables()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze tables: %w", err)
	}

	insights := s.insights.GenerateInsights(tables)
	summary := s.reporter.GenerateSummary(tables)
	recommendations := s.reporter.GenerateRecommendations(tables, insights)

	var performanceMetrics types.PerformanceMetrics
	if s.performance != nil {
		performanceMetrics = s.performance.AnalyzePerformance()
	}

	report := &types.DatabaseReport{
		DatabaseName:       dbName,
		DatabaseType:       dbType,
		AnalysisTime:       time.Now(),
		Summary:            summary,
		Tables:             tables,
		Insights:           insights,
		Recommendations:    recommendations,
		PerformanceMetrics: performanceMetrics,
	}

	log.Println("Database analysis completed successfully")
	return report, nil
}

func (s *Service) AnalyzeSecurity() ([]types.SecurityIssue, error) {
	if s.mongoService != nil {
		ctx := context.Background()
		return s.mongoService.AnalyzeSecurity(ctx)
	}
	return s.security.AnalyzeSecurity()
}

func (s *Service) OptimizeQuery(query string) (*types.OptimizationSuggestion, error) {
	if s.mongoService != nil {
		return s.mongoService.OptimizeQuery(query)
	}
	return s.optimizer.AnalyzeQuery(query)
}

func (s *Service) CheckAlerts() ([]types.MonitoringAlert, error) {
	if s.mongoService != nil {
		ctx := context.Background()
		return s.mongoService.CheckAlerts(ctx)
	}

	report, err := s.AnalyzeDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze database for alerts: %w", err)
	}

	return s.alerts.CheckAlerts(report), nil
}

func (s *Service) CompareReports(oldReport, newReport *types.DatabaseReport) *types.ComparisonReport {
	return s.comparison.CompareReports(oldReport, newReport)
}

func (s *Service) TrackLineage() (map[string]types.DataLineage, error) {
	if s.mongoService != nil {
		ctx := context.Background()
		return s.mongoService.TrackLineage(ctx)
	}
	return s.lineage.TrackLineage()
}

func (s *Service) ScheduleAnalysis(interval time.Duration, callback func(*types.DatabaseReport)) error {
	return s.scheduler.ScheduleAnalysis(interval, callback)
}

func (s *Service) StopScheduledAnalysis() error {
	return s.scheduler.Stop()
}

func (s *Service) ExportReport(report *types.DatabaseReport, format string) ([]byte, error) {
	if s.mongoService != nil {
		return s.mongoService.ExportReport(report, format)
	}
	return s.reporter.ExportReport(report, format)
}

func (s *Service) GetConfig() *types.Config {
	return s.config.GetConfig()
}

func (s *Service) UpdateConfig(config *types.Config) error {
	return s.config.UpdateConfig(config)
}

func (s *Service) Close() error {
	if s.mongoService != nil {
		ctx := context.Background()
		return s.mongoService.Close(ctx)
	}

	if s.scheduler != nil {
		if err := s.scheduler.Stop(); err != nil {
			log.Printf("Warning: Failed to stop scheduler: %v", err)
		}
	}

	if s.connector != nil {
		if err := s.connector.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}

	return nil
}
