package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
	"github.com/cherry-pick/pkg/loadbalancer/engine"
)

type Manager struct {
	engines map[string]core.LoadTestEngine
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		engines: make(map[string]core.LoadTestEngine),
	}
}

func (m *Manager) CreateEngine(engineID string) (core.LoadTestEngine, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.engines[engineID]; exists {
		return nil, fmt.Errorf("engine with ID %s already exists", engineID)
	}

	eng := engine.NewEngine()
	m.engines[engineID] = eng
	return eng, nil
}

func (m *Manager) GetEngine(engineID string) (core.LoadTestEngine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	eng, exists := m.engines[engineID]
	if !exists {
		return nil, fmt.Errorf("engine with ID %s not found", engineID)
	}

	return eng, nil
}

func (m *Manager) DeleteEngine(engineID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.engines[engineID]; !exists {
		return fmt.Errorf("engine with ID %s not found", engineID)
	}

	delete(m.engines, engineID)
	return nil
}

func (m *Manager) ListEngines() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	engines := make([]string, 0, len(m.engines))
	for id := range m.engines {
		engines = append(engines, id)
	}

	return engines
}

func (m *Manager) GetDefaultEngine() core.LoadTestEngine {
	m.mu.Lock()
	defer m.mu.Unlock()

	if eng, exists := m.engines["default"]; exists {
		return eng
	}

	eng := engine.NewEngine()
	m.engines["default"] = eng
	return eng
}

func (m *Manager) StartLoadTest(testID string, config core.LoadTestConfig) error {
	eng := m.GetDefaultEngine()
	return eng.StartLoadTest(testID, config)
}

func (m *Manager) GetTestStatus(testID string) (*core.LoadTestStatus, error) {
	eng := m.GetDefaultEngine()
	return eng.GetTestStatus(testID)
}

func (m *Manager) GetTestSummary(testID string) (*core.LoadTestSummary, error) {
	eng := m.GetDefaultEngine()
	return eng.GetTestSummary(testID)
}

func (m *Manager) GetTestResults(testID string) ([]core.LoadTestResult, error) {
	eng := m.GetDefaultEngine()
	return eng.GetTestResults(testID)
}

func (m *Manager) CancelTest(testID string) error {
	eng := m.GetDefaultEngine()
	return eng.CancelTest(testID)
}

func (m *Manager) GetAllTests() map[string]*core.LoadTestStatus {
	eng := m.GetDefaultEngine()
	return eng.GetAllTests()
}

func (m *Manager) GetRealTimeMetrics(testID string) (*core.RealTimeMetrics, error) {
	eng := m.GetDefaultEngine()
	return eng.GetRealTimeMetrics(testID)
}

func (m *Manager) CleanupOldTests(olderThan time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, eng := range m.engines {
		eng.CleanupOldTests(olderThan)
	}
}

func (m *Manager) GetEngineStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})

	for engineID, eng := range m.engines {
		engineStats := map[string]interface{}{
			"totalTests":     0,
			"runningTests":   0,
			"completedTests": 0,
			"failedTests":    0,
			"cancelledTests": 0,
		}

		tests := eng.GetAllTests()
		engineStats["totalTests"] = len(tests)

		for _, status := range tests {
			switch status.Status {
			case "running":
				engineStats["runningTests"] = engineStats["runningTests"].(int) + 1
			case "completed":
				engineStats["completedTests"] = engineStats["completedTests"].(int) + 1
			case "failed":
				engineStats["failedTests"] = engineStats["failedTests"].(int) + 1
			case "cancelled":
				engineStats["cancelledTests"] = engineStats["cancelledTests"].(int) + 1
			}
		}

		stats[engineID] = engineStats
	}

	return stats
}
