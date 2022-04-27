//go:build darwin
// +build darwin

package cpu

import (
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"os/exec"
	"strconv"
	"strings"
)

func Get() (*app.CpuStats, error) {
	cmd := exec.Command("iostat")
	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	resFields := strings.Fields(string(res))
	user, err := strconv.Atoi(resFields[16])
	if err != nil {
		return nil, err
	}

	system, err := strconv.Atoi(resFields[17])
	if err != nil {
		return nil, err
	}

	idle, err := strconv.Atoi(resFields[18])
	if err != nil {
		return nil, err
	}

	return &app.CpuStats{
		User:   user,
		System: system,
		Idle:   idle,
	}, nil
}
