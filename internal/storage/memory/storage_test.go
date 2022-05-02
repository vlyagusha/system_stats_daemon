package memorystorage

import (
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/storage"
)

func TestStorage(t *testing.T) { //nolint:funlen,gocognit,nolintlint
	store := New()

	t.Run("saving test", func(t *testing.T) {
		collectedAt, err := time.Parse("2006-01-02 15:04:05", "2022-05-01 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}
		stat := app.SystemStats{
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

		err = store.Create(stat)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err := store.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, stat, saved[0])

		err = store.Delete(stat.ID)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err = store.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 0)
	})

	t.Run("test errors", func(t *testing.T) {
		collectedAt, err := time.Parse("2006-01-02 15:04:05", "2022-05-01 12:00:00")
		statID := uuid.New()
		stat := app.SystemStats{
			ID:          statID,
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

		err = store.Create(stat)
		require.Nil(t, err)

		err = store.Create(stat)
		require.ErrorIs(t, err, storage.ErrObjectAlreadyExists)

		err = store.Delete(statID)
		require.Nil(t, err)

		err = store.Delete(statID)
		require.ErrorIs(t, err, storage.ErrObjectDoesNotExist)
	})

	t.Run("test get avg simple", func(t *testing.T) {
		stats := []app.SystemStats{
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597150"),
				CollectedAt: time.Now().Add(-10 * time.Second),
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
					KBt: 5,
					TPS: 10,
					MBs: 7,
				},
			},
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597151"),
				CollectedAt: time.Now().Add(-20 * time.Second),
				Load: &app.LoadStats{
					Load1:  2,
					Load5:  10,
					Load15: 30,
				},
				CPU: &app.CPUStats{
					User:   15,
					System: 15,
					Idle:   70,
				},
				Disk: &app.DiskStats{
					KBt: 50,
					TPS: 20,
					MBs: 7,
				},
			},
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597152"),
				CollectedAt: time.Now().Add(-30 * time.Second),
				Load: &app.LoadStats{
					Load1:  2,
					Load5:  15,
					Load15: 25,
				},
				CPU: &app.CPUStats{
					User:   15,
					System: 25,
					Idle:   60,
				},
				Disk: &app.DiskStats{
					KBt: 50,
					TPS: 5,
					MBs: 70,
				},
			},
		}

		for _, e := range stats {
			err := store.Create(e)
			if err != nil {
				t.FailNow()
				return
			}
		}

		avg, err := store.FindAvg(60 * time.Second)
		require.Nil(t, err)
		require.Equal(t, math.Round(avg.Load1*100)/100, 1.67)
		require.Equal(t, math.Round(avg.Load5*100)/100, 10.0)
		require.Equal(t, math.Round(avg.Load15*100)/100, 23.33)
		require.Equal(t, math.Round(avg.User*100)/100, 13.33)
		require.Equal(t, math.Round(avg.System*100)/100, 20.0)
		require.Equal(t, math.Round(avg.Idle*100)/100, 66.67)
		require.Equal(t, math.Round(avg.KBt*100)/100, 35.0)
		require.Equal(t, math.Round(avg.TPS*100)/100, 11.67)
		require.Equal(t, math.Round(avg.MBs*100)/100, 28.0)
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
