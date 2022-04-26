//go:build linux
// +build linux

package load

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func Avg() (float64, error) {
	stat, err := fileAvg()
	if err != nil {
		stat, err = sysInfoAvg()
	}
	return stat, err
}

func sysInfoAvg() (float64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return 0, err
	}

	return float64(info.Loads[0]) / float64(1<<16), nil
}

func fileAvg() (float64, error) {
	values, err := readLoadAvgFromFile()
	if err != nil {
		return 0, err
	}

	load, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return 0, err
	}

	return load, nil
}

func readLoadAvgFromFile() ([]string, error) {
	procDir := os.Getenv("HOST_PROC")
	if procDir == "" {
		procDir = "/proc"
	}

	line, err := ioutil.ReadFile(filepath.Join(procDir, "/loadavg"))
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(line)), nil
}
