package memorystorage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
)

func TestStorage(t *testing.T) { //nolint:funlen,gocognit,nolintlint
	storage := New()

	t.Run("saving test", func(t *testing.T) {
		collectedAt, err := time.Parse("2006-01-02 15:04:05", "2022-05-01 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}
		stats := app.SystemStats{
			ID:          uuid.New(),
			CollectedAt: collectedAt,
			Load: &app.LoadStats{
				Load1:  1,
				Load5:  5,
				Load15: 15,
			},
			CPU: &app.CPUStats{
				User:   10,
				System: 20,
				Idle:   70,
			},
			Disk: &app.DiskStats{
				KBt: 5.6,
				TPS: 12,
				MBs: 9.3,
			},
		}

		err = storage.Create(stats)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err := storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, stats, saved[0])

		err = storage.Delete(stats.ID)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 0)
	})

	t.Run("test get avg simple", func(t *testing.T) {
		stats := []app.SystemStats{
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597150"),
				CollectedAt: parseDate(t, "2022-05-01T12:00:00Z"),
				Load:        nil,
				CPU:         nil,
				Disk:        nil,
			},
		}

		for _, e := range stats {
			err := storage.Create(e)
			if err != nil {
				t.FailNow()
				return
			}
		}
	})
}

func parseUUID(t *testing.T, str string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(str)
	if err != nil {
		t.Errorf("failed to parse UUID: %s", err)
	}
	return id
}

func parseDate(t *testing.T, str string) time.Time {
	t.Helper()
	dt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Errorf("failed to parse date: %s", err)
	}
	return dt
}
