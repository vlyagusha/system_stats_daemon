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
	Main        mainStats
	Disk        diskStats
}

func (s SystemStats) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ID: %s ", s.ID))
	builder.WriteString(fmt.Sprintf("Load average: %f ", s.Main.Load))

	return builder.String()
}

type mainStats struct {
	Load float64
	CPU  string
}

type diskStats struct {
	TPS string
	KBS string
}
