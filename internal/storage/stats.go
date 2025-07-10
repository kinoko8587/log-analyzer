package storage

import (
	"sync"
	"time"

	"github.com/pinjung/log-analyzer/pkg/analyzer"
)

type Stats struct {
	mu         sync.RWMutex
	errorCount int64
	totalCount int64
	lastUpdate time.Time
}

type StatsSnapshot struct {
	ErrorCount int64     `json:"error_count"`
	TotalCount int64     `json:"total_count"`
	LastUpdate time.Time `json:"last_update"`
}

func NewStats() *Stats {
	return &Stats{
		lastUpdate: time.Now(),
	}
}

func (s *Stats) RecordLog(log *analyzer.Log) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.totalCount++
	if log.IsError() {
		s.errorCount++
	}
	s.lastUpdate = time.Now()
}

func (s *Stats) GetErrorCount() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errorCount
}

func (s *Stats) GetTotalCount() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalCount
}

func (s *Stats) GetSnapshot() StatsSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return StatsSnapshot{
		ErrorCount: s.errorCount,
		TotalCount: s.totalCount,
		LastUpdate: s.lastUpdate,
	}
}

func (s *Stats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.errorCount = 0
	s.totalCount = 0
	s.lastUpdate = time.Now()
}