package desktop

import (
	"errors"
	"os/exec"
	"runtime"
)

// ErrDesktopNotSupported is returned when the current platform isn't supported by lighthouse desktop
var ErrDesktopNotSupported = errors.New("Current platform is not supported by lighthouse desktop, try running lighthouse with the -no-gui flag")

// LaunchDesktopApp launches the lighthouse desktop app
func LaunchDesktopApp() error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("./build/lighthouse-desktop.exe")
		return cmd.Start()
	}

	return ErrDesktopNotSupported
}
