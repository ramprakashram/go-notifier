package notifier

import (
	"os/exec"
)

func Notify(title string, text string, level string) {
	cmd := exec.Command("notify-send", title, text, "-u", level)
	cmd.Run()
}
