package memorystorage

import (
	"github.com/google/uuid"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	mu    sync.RWMutex
	stats map[uuid.UUID]app.SystemStats
}

func New() *Storage {
	return &Storage{
		stats: make(map[uuid.UUID]app.SystemStats),
	}
}

func (m *Storage) Create(s app.SystemStats) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stats[s.ID]; ok {
		return storage.ErrObjectAlreadyExists
	}

	m.stats[s.ID] = s
	return nil
}

func (m *Storage) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stats[id]; !ok {
		return storage.ErrObjectDoesNotExist
	}

	delete(m.stats, id)
	return nil
}

func (m *Storage) FindAll() ([]app.SystemStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make([]app.SystemStats, len(m.stats))
	for _, systemStats := range m.stats {
		stats = append(stats, systemStats)
	}

	return stats, nil
}

func (m *Storage) FindAvg(duration int32) (*app.SystemStatsAvg, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	var load1, load5, load15, user, system, idle, kbt, tps, mbs float64
	for _, systemStats := range m.stats {
		if now.Sub(systemStats.CollectedAt) <= (time.Duration(duration) * time.Second) {
			load1 = load1 + systemStats.Load.Load1
			load5 = load5 + systemStats.Load.Load5
			load15 = load15 + systemStats.Load.Load15
			user = user + float64(systemStats.CPU.User)
			system = system + float64(systemStats.CPU.System)
			idle = idle + float64(systemStats.CPU.Idle)
			kbt = kbt + systemStats.Disk.KBt
			tps = tps + float64(systemStats.Disk.TPS)
			mbs = mbs + systemStats.Disk.MBs
		}
	}

	totalItems := float64(len(m.stats))

	return &app.SystemStatsAvg{
		Load1:  load1 / totalItems,
		Load5:  load5 / totalItems,
		Load15: load15 / totalItems,
		User:   user / totalItems,
		System: system / totalItems,
		Idle:   idle / totalItems,
		KBt:    kbt / totalItems,
		TPS:    tps / totalItems,
		MBs:    mbs / totalItems,
	}, nil
}
