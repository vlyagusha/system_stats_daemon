//go:build darwin
// +build darwin

package disk

import (
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"os/exec"
	"strconv"
	"strings"
)

func Get() (*app.DiskStats, error) {
	cmd := exec.Command("iostat")
	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	resFields := strings.Fields(string(res))

	kbt, err := strconv.ParseFloat(resFields[13], 64)
	if err != nil {
		return nil, err
	}

	tps, err := strconv.Atoi(resFields[14])
	if err != nil {
		return nil, err
	}

	mbs, err := strconv.ParseFloat(resFields[15], 64)
	if err != nil {
		return nil, err
	}

	return &app.DiskStats{
		KBt: kbt,
		TPS: tps,
		MBs: mbs,
	}, nil
}
