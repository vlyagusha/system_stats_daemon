package app

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type SystemStats struct {
	ID          uuid.UUID
	CollectedAt time.Time
	Load        LoadStats
	CPU         CpuStats
	Disk        DiskStats
}

func (s SystemStats) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ID: %s ", s.ID))
	builder.WriteString(fmt.Sprintf("Load average: %f ", s.Load))
	builder.WriteString(fmt.Sprintf("CPU usage: %s ", s.CPU))
	builder.WriteString(fmt.Sprintf("Disk: %s ", s.Disk))

	return builder.String()
}

type LoadStats struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func (l LoadStats) String() string {
	return fmt.Sprintf("1m: %2f 5m: %2f 15m: %2f", l.Load1, l.Load5, l.Load15)
}

type CpuStats struct {
	User   int
	System int
	Idle   int
}

func (c CpuStats) String() string {
	return fmt.Sprintf("User: %d System: %d Idle: %d", c.User, c.System, c.Idle)
}

type DiskStats struct {
	KBt float64
	TPS int
	MBs float64
}

func (d DiskStats) String() string {
	return fmt.Sprintf("KB/t: %2f tps: %d MB/s: %2f", d.KBt, d.TPS, d.MBs)
}
