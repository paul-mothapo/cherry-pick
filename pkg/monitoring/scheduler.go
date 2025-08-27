// Package monitoring provides scheduling functionality for automated analysis.
package monitoring

import (
	"fmt"
	"sync"
	"time"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
)

// SchedulerImpl implements the Scheduler interface.
type SchedulerImpl struct {
	analyzer interfaces.DatabaseAnalyzer
	reporter interfaces.ReportGenerator
	insights interfaces.InsightGenerator
	ticker   *time.Ticker
	stopChan chan bool
	running  bool
	mutex    sync.RWMutex
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(analyzer interfaces.DatabaseAnalyzer, reporter interfaces.ReportGenerator, insights interfaces.InsightGenerator) interfaces.Scheduler {
	return &SchedulerImpl{
		analyzer: analyzer,
		reporter: reporter,
		insights: insights,
		stopChan: make(chan bool),
	}
}

// ScheduleAnalysis allows for automated periodic analysis.
func (s *SchedulerImpl) ScheduleAnalysis(interval time.Duration, callback func(*types.DatabaseReport)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return fmt.Errorf("scheduler is already running")
	}

	s.ticker = time.NewTicker(interval)
	s.running = true

	go func() {
		defer func() {
			s.mutex.Lock()
			s.running = false
			s.mutex.Unlock()
		}()

		for {
			select {
			case <-s.ticker.C:
				report, err := s.generateReport()
				if err != nil {
					// Log error but continue scheduling
					fmt.Printf("Scheduled analysis failed: %v\n", err)
					continue
				}
				callback(report)
			case <-s.stopChan:
				return
			}
		}
	}()

	return nil
}

// Stop stops the scheduled analysis.
func (s *SchedulerImpl) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return fmt.Errorf("scheduler is not running")
	}

	s.ticker.Stop()
	s.stopChan <- true
	s.running = false

	return nil
}

// generateReport creates a complete database report.
func (s *SchedulerImpl) generateReport() (*types.DatabaseReport, error) {
	// Analyze all tables
	tables, err := s.analyzer.AnalyzeTables()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze tables: %w", err)
	}

	// Generate insights
	insights := s.insights.GenerateInsights(tables)

	// Generate summary
	summary := s.reporter.GenerateSummary(tables)

	// Generate recommendations
	recommendations := s.reporter.GenerateRecommendations(tables, insights)

	report := &types.DatabaseReport{
		DatabaseName:    "Scheduled Analysis",
		DatabaseType:    "Unknown", // Would need database connector
		AnalysisTime:    time.Now(),
		Summary:         summary,
		Tables:          tables,
		Insights:        insights,
		Recommendations: recommendations,
		// PerformanceMetrics would need performance analyzer
	}

	return report, nil
}
