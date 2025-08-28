package intelligence

import (
	"context"
	"fmt"
	"strings"

	"github.com/cherry-pick/pkg/analyzer"
	"github.com/cherry-pick/pkg/config"
	"github.com/cherry-pick/pkg/connector"
	"github.com/cherry-pick/pkg/insights"
	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/monitoring"
	"github.com/cherry-pick/pkg/optimization"
	"github.com/cherry-pick/pkg/security"
)

type ServiceBuilder struct {
	driverName     string
	dataSourceName string
	configPath     string
}

func NewServiceBuilder(driverName, dataSourceName string) *ServiceBuilder {
	return &ServiceBuilder{
		driverName:     driverName,
		dataSourceName: dataSourceName,
	}
}

func (sb *ServiceBuilder) WithConfig(configPath string) *ServiceBuilder {
	sb.configPath = configPath
	return sb
}

func (sb *ServiceBuilder) Build() (*Service, error) {
	if strings.ToLower(sb.driverName) == "mongodb" {
		return sb.buildMongoService()
	}

	dbConnector := connector.NewDatabaseConnector(sb.driverName, sb.dataSourceName)
	if err := dbConnector.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbAnalyzer := analyzer.NewDatabaseAnalyzer(dbConnector.GetDB(), sb.driverName)

	insightGenerator := insights.NewInsightGenerator()
	reportGenerator := insights.NewReportGenerator()
	securityAnalyzer := security.NewSecurityAnalyzer(dbAnalyzer)
	queryOptimizer := optimization.NewQueryOptimizer()
	alertManager := monitoring.NewAlertManager()
	comparisonEngine := monitoring.NewComparisonEngine()
	lineageTracker := monitoring.NewDataLineageTracker(dbAnalyzer)
	scheduler := monitoring.NewScheduler(dbAnalyzer, reportGenerator, insightGenerator)

	configManager := config.NewConfigManager()
	if err := configManager.LoadConfig(sb.configPath); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	service := NewService(
		dbConnector,
		dbAnalyzer,
		insightGenerator,
		reportGenerator,
		securityAnalyzer,
		queryOptimizer,
		alertManager,
		comparisonEngine,
		lineageTracker,
		scheduler,
		configManager,
		nil,
	)

	return service, nil
}

func (sb *ServiceBuilder) buildMongoService() (*Service, error) {
	databaseName := "test"
	if strings.Contains(sb.dataSourceName, "/") {
		parts := strings.Split(sb.dataSourceName, "/")
		if len(parts) > 3 {
			dbPart := parts[3]
			if strings.Contains(dbPart, "?") {
				databaseName = strings.Split(dbPart, "?")[0]
			} else {
				databaseName = dbPart
			}
		}
	}

	mongoConnector := connector.NewMongoConnector(sb.dataSourceName, databaseName)
	ctx := context.Background()
	if err := mongoConnector.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	mongoAnalyzer := analyzer.NewMongoAnalyzer(mongoConnector)

	service := NewMongoService(
		mongoConnector,
		mongoAnalyzer,
		sb.configPath,
	)

	return service, nil
}

func CreateSimpleService(driverName, dataSourceName string) (*Service, error) {
	return NewServiceBuilder(driverName, dataSourceName).Build()
}

func CreateConfiguredService(driverName, dataSourceName, configPath string) (*Service, error) {
	return NewServiceBuilder(driverName, dataSourceName).
		WithConfig(configPath).
		Build()
}

type MockServiceBuilder struct {
	connector   interfaces.DatabaseConnector
	analyzer    interfaces.DatabaseAnalyzer
	insights    interfaces.InsightGenerator
	reporter    interfaces.ReportGenerator
	security    interfaces.SecurityAnalyzer
	optimizer   interfaces.QueryOptimizer
	alerts      interfaces.AlertManager
	comparison  interfaces.ComparisonEngine
	lineage     interfaces.DataLineageTracker
	scheduler   interfaces.Scheduler
	config      interfaces.ConfigManager
	performance interfaces.PerformanceAnalyzer
}

func NewMockServiceBuilder() *MockServiceBuilder {
	return &MockServiceBuilder{}
}

func (msb *MockServiceBuilder) WithConnector(connector interfaces.DatabaseConnector) *MockServiceBuilder {
	msb.connector = connector
	return msb
}

func (msb *MockServiceBuilder) WithAnalyzer(analyzer interfaces.DatabaseAnalyzer) *MockServiceBuilder {
	msb.analyzer = analyzer
	return msb
}

func (msb *MockServiceBuilder) WithInsights(insights interfaces.InsightGenerator) *MockServiceBuilder {
	msb.insights = insights
	return msb
}

func (msb *MockServiceBuilder) WithReporter(reporter interfaces.ReportGenerator) *MockServiceBuilder {
	msb.reporter = reporter
	return msb
}

func (msb *MockServiceBuilder) WithSecurity(security interfaces.SecurityAnalyzer) *MockServiceBuilder {
	msb.security = security
	return msb
}

func (msb *MockServiceBuilder) WithOptimizer(optimizer interfaces.QueryOptimizer) *MockServiceBuilder {
	msb.optimizer = optimizer
	return msb
}

func (msb *MockServiceBuilder) WithAlerts(alerts interfaces.AlertManager) *MockServiceBuilder {
	msb.alerts = alerts
	return msb
}

func (msb *MockServiceBuilder) WithComparison(comparison interfaces.ComparisonEngine) *MockServiceBuilder {
	msb.comparison = comparison
	return msb
}

func (msb *MockServiceBuilder) WithLineage(lineage interfaces.DataLineageTracker) *MockServiceBuilder {
	msb.lineage = lineage
	return msb
}

func (msb *MockServiceBuilder) WithScheduler(scheduler interfaces.Scheduler) *MockServiceBuilder {
	msb.scheduler = scheduler
	return msb
}

func (msb *MockServiceBuilder) WithConfig(config interfaces.ConfigManager) *MockServiceBuilder {
	msb.config = config
	return msb
}

func (msb *MockServiceBuilder) WithPerformance(performance interfaces.PerformanceAnalyzer) *MockServiceBuilder {
	msb.performance = performance
	return msb
}

func (msb *MockServiceBuilder) Build() *Service {
	return NewService(
		msb.connector,
		msb.analyzer,
		msb.insights,
		msb.reporter,
		msb.security,
		msb.optimizer,
		msb.alerts,
		msb.comparison,
		msb.lineage,
		msb.scheduler,
		msb.config,
		msb.performance,
	)
}
