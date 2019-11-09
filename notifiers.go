package notifier

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	gosxnotifier "github.com/deckarep/gosx-notifier"
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

func notifyDarwin(title string, text string, level string) {
	note := gosxnotifier.NewNotification(text)
	note.Title = title
	err := note.Push()
	if err != nil {
		log.Println("Error while pushing notification")
	}
}

func Notify(title string, text string, level string) {
	switch os := runtime.GOOS; os {
	case LINUX:
		notifyLinux(title, text, level)
	case WINDOWS:
		notSupported()
	case OSX:
		notifyDarwin(title, text, level)
	default:
		notSupported()
	}
}
