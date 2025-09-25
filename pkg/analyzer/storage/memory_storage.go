package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type MemoryStorage struct {
	analyses map[string]*core.AnalysisResult
	mu       sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		analyses: make(map[string]*core.AnalysisResult),
	}
}

func (ms *MemoryStorage) SaveAnalysis(ctx context.Context, result *core.AnalysisResult) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.analyses[result.ID] = result
	return nil
}

func (ms *MemoryStorage) GetAnalysis(ctx context.Context, analysisID string) (*core.AnalysisResult, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	result, exists := ms.analyses[analysisID]
	if !exists {
		return nil, fmt.Errorf("analysis not found: %s", analysisID)
	}

	return result, nil
}

func (ms *MemoryStorage) GetAnalysisHistory(ctx context.Context, limit int) ([]core.AnalysisResult, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var results []core.AnalysisResult
	count := 0

	for _, result := range ms.analyses {
		if limit > 0 && count >= limit {
			break
		}
		results = append(results, *result)
		count++
	}

	return results, nil
}

func (ms *MemoryStorage) DeleteAnalysis(ctx context.Context, analysisID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.analyses[analysisID]; !exists {
		return fmt.Errorf("analysis not found: %s", analysisID)
	}

	delete(ms.analyses, analysisID)
	return nil
}

func (ms *MemoryStorage) CleanupOldAnalyses(ctx context.Context, olderThan time.Duration) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	cutoff := time.Now().Add(-olderThan)
	var toDelete []string

	for id, result := range ms.analyses {
		if result.AnalysisTime.Before(cutoff) {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(ms.analyses, id)
	}

	return nil
}
