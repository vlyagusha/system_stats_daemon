package memorystorage

import (
	"github.com/google/uuid"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/storage"
	"sync"
)

type MemoryStorage struct {
	mu    sync.RWMutex
	stats map[uuid.UUID]app.SystemStats
}

func (m *MemoryStorage) Create(s app.SystemStats) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stats[s.ID]; ok {
		return storage.ErrObjectAlreadyExists
	}

	m.stats[s.ID] = s
	return nil
}

func (m *MemoryStorage) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stats[id]; !ok {
		return storage.ErrObjectDoesNotExist
	}

	delete(m.stats, id)
	return nil
}

func (m *MemoryStorage) FindAll() ([]app.SystemStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make([]app.SystemStats, len(m.stats))
	for _, systemStats := range m.stats {
		stats = append(stats, systemStats)
	}

	return stats, nil
}
