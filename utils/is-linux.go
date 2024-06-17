package utils

import "runtime"

// is linux?
func IsLinux() bool {
	return runtime.GOOS == "linux"
}
