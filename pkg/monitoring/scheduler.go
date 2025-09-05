package monitoring

import (
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/interfaces"
	"github.com/cherry-pick/pkg/types"
)

type SchedulerImpl struct {
	analyzer interfaces.DatabaseAnalyzer
	reporter interfaces.ReportGenerator
	insights interfaces.InsightGenerator
	ticker   *time.Ticker
	stopChan chan bool
	running  bool
	mutex    sync.RWMutex
}

func NewScheduler(analyzer interfaces.DatabaseAnalyzer, reporter interfaces.ReportGenerator, insights interfaces.InsightGenerator) interfaces.Scheduler {
	return &SchedulerImpl{
		analyzer: analyzer,
		reporter: reporter,
		insights: insights,
		stopChan: make(chan bool),
	}
}

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

func (s *SchedulerImpl) generateReport() (*types.DatabaseReport, error) {
	tables, err := s.analyzer.AnalyzeTables()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze tables: %w", err)
	}

	insights := s.insights.GenerateInsights(tables)
	summary := s.reporter.GenerateSummary(tables)
	recommendations := s.reporter.GenerateRecommendations(tables, insights)

	report := &types.DatabaseReport{
		DatabaseName:    "Scheduled Analysis",
		DatabaseType:    "Unknown",
		AnalysisTime:    time.Now(),
		Summary:         summary,
		Tables:          tables,
		Insights:        insights,
		Recommendations: recommendations,
	}

	return report, nil
}
