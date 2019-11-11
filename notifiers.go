package notifier

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"text/template"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	uuid "github.com/nu7hatch/gouuid"
)

type WindowsToaster struct {
	AppID    string
	Title    string
	Message  string
	Duration string
}

const (
	LINUX   = "linux"
	WINDOWS = "windows"
	OSX     = "darwin"
)

var windowToastPlaceholder *template.Template

func init() {
	windowToastPlaceholder = template.New("toast")
	windowToastPlaceholder.Parse(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}Windows App{{end}}'

$template = @"
<toast activationType="protocol" launch="" duration="short">
    <visual>
        <binding template="ToastGeneric">
            {{if .Title}}
            <text><![CDATA[{{.Title}}]]></text>
            {{end}}
            {{if .Message}}
            <text><![CDATA[{{.Message}}]]></text>
      			{{end}}
         </binding>
    </visual>
	<audio silent="true" />
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
    `)
}

func setDuration(level string) string {
	switch lvl := level; lvl {
	case LOW:
		return "short"
	case NORMAL:
		return "long"
	case CRITICAL:
		return "long"
	default:
		return "short"
	}
}

func parseWindowsToaster(applicationName string, title string, text string, level string) WindowsToaster {
	return WindowsToaster{
		AppID:    applicationName,
		Title:    title,
		Message:  text,
		Duration: setDuration(level),
	}
}

func (toastData *WindowsToaster) generateWindowsXML() (string, error) {
	var toasts bytes.Buffer
	error := windowToastPlaceholder.Execute(&toasts, toastData)
	if error != nil {
		return "", error
	}
	return toasts.String(), nil
}

func initiateWindowsNotification(toast string) error {
	id, _ := uuid.NewV4()
	file := filepath.Join(os.TempDir(), id.String()+".ps1")
	defer os.Remove(file)
	err := ioutil.WriteFile(file, []byte(toast), 0600)
	if err != nil {
		return err
	}
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", "-File", file)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (toastData *WindowsToaster) Notify() error {
	windowsXML, error := toastData.generateWindowsXML()
	if error != nil {
		return error
	}

	return initiateWindowsNotification(windowsXML)
}

func notSupported() {
	fmt.Println("Support for this OS is not available")
}

func notifyLinux(title string, text string, level string) {
	cmd := exec.Command("notify-send", title, text, "-u", level)
	cmd.Run()
}

func notifyWindows(applicationName string, title string, text string, level string) {
	windowToast := parseWindowsToaster(applicationName, title, text, level)
	windowToast.Notify()
}

func notifyDarwin(title string, text string, level string) {
	note := gosxnotifier.NewNotification(text)
	note.Title = title
	err := note.Push()
	if err != nil {
		log.Println("Error while pushing notification")
	}
}

func Notify(applicationName string, title string, text string, level string) {
	switch os := runtime.GOOS; os {
	case LINUX:
		notifyLinux(title, text, level)
	case WINDOWS:
		notifyWindows(applicationName, title, text, level)
	case OSX:
		notifyDarwin(title, text, level)
	default:
		notSupported()
	}
}
