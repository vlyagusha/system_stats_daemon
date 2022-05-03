package utils

import (
	"fmt"
	"runtime"
)

var ErrNotImplemented = fmt.Errorf("cpu statistics not implemented for: %s", runtime.GOOS)
