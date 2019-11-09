package notifier

import (
	"fmt"
	"os/exec"
	"runtime"
)

const (
	LINUX   = "linux"
	WINDOWS = "windows"
	OSX     = "darwin"
)

func notSupported() {
	fmt.Println("Support for this OS is not available")
}

func notifyLinux(title string, text string, level string) {
	cmd := exec.Command("notify-send", title, text, "-u", level)
	cmd.Run()
}

func Notify(title string, text string, level string) {
	switch os := runtime.GOOS; os {
	case LINUX:
		notifyLinux(title, text, level)
	case WINDOWS:
		notSupported()
	case OSX:
		notSupported()
	default:
		notSupported()
	}
}
