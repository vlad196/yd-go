// Copyleft 2017-2021 Sly_tom_cat (slytomcat@mail.ru)
// License: GPL v.3

//go:generate gotext -srclang=en update -out=catalog.go -lang=en,ru

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/slytomcat/llog"
	"github.com/slytomcat/yd-go/icons"
	"github.com/slytomcat/yd-go/tools"
	"github.com/slytomcat/ydisk"
	"golang.org/x/text/message"
)

var (
	version = "local build"
	// msg is the Localization printer
	msg      *message.Printer
	statusTr map[string]string
)

const about = `yd-go is the GTK-based panel indicator for Yandex.Disk daemon.

	Version: %s

Copyleft 2017-%s Sly_tom_cat (slytomcat@mail.ru)

	License: GPL v.3

`

func notifySend(title, body string) {
	llog.Debug("Message:", title, ":", body)
	err := exec.Command("notify-send", "-i", icons.NotifyIcon, title, body).Start()
	if err != nil {
		llog.Error(err)
	}
}

type menu struct {
	status *systray.MenuItem
	size1  *systray.MenuItem
	size2  *systray.MenuItem
	last   *systray.MenuItem
	lastM  [10]*systray.MenuItem
	lastP  [10]string
	start  *systray.MenuItem
	stop   *systray.MenuItem
	out    *systray.MenuItem
	path   *systray.MenuItem
	site   *systray.MenuItem
	help   *systray.MenuItem
	about  *systray.MenuItem
	don    *systray.MenuItem
	quit   *systray.MenuItem
}

// Icon is the busy icon animation helper
type Icon struct {
	current int
	Tick    *time.Timer
}

// Next returns next busy icon index and resets the timer
func (i *Icon) Next() int {
	i.current = (i.current + 1) % 5
	i.Tick.Reset(333 * time.Millisecond)
	return i.current
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	// Initialize application and receive the application configuration
	AppCfg := tools.AppInit("yd-go")

	// Initialize translations
	lng := os.Getenv("LANG")
	if len(lng) > 2 {
		lng = lng[:2]
	}

	llog.Infof("Local language is: %v", lng)
	msg = message.NewPrinter(message.MatchLanguage(lng))

	YD, err := ydisk.NewYDisk(AppCfg.Conf)
	if err != nil {
		llog.Critical("Fatal error:", err)
	}

	// Initialize icon theme
	icons.SelectTheme(AppCfg.Theme)
	// Initialize systray icon
	systray.SetIcon(icons.PauseIcon)

	// Initialize status localization
	statusTr = map[string]string{
		"idle":   msg.Sprintf("idle"),
		"index":  msg.Sprintf("index"),
		"busy":   msg.Sprintf("busy"),
		"none":   msg.Sprintf("none"),
		"paused": msg.Sprintf("paused"),
	}

	m := new(menu)
	systray.SetTitle("yd-go indicator")
	// Initialize systray menu
	m.status = systray.AddMenuItem("", "")
	m.status.Disable()
	m.size1 = systray.AddMenuItem("", "")
	m.size1.Disable()
	m.size2 = systray.AddMenuItem("", "")
	m.size2.Disable()
	systray.AddSeparator()
	m.last = systray.AddMenuItem(msg.Sprintf("Last synchronized"), "")
	for i := 0; i < 10; i++ {
		m.lastM[i] = m.last.AddSubMenuItem("", "")
		m.lastM[i].Hide()
	}
	systray.AddSeparator()
	m.start = systray.AddMenuItem(msg.Sprintf("Start daemon"), "")
	m.start.Hide()
	m.stop = systray.AddMenuItem(msg.Sprintf("Stop daemon"), "")
	m.stop.Hide()
	systray.AddSeparator()
	m.out = systray.AddMenuItem(msg.Sprintf("Show daemon output"), "")
	m.path = systray.AddMenuItem(msg.Sprintf("Open: %s", YD.Path), "")
	m.site = systray.AddMenuItem(msg.Sprintf("Open YandexDisk in browser"), "")
	systray.AddSeparator()
	m.help = systray.AddMenuItem(msg.Sprintf("Help"), "")
	m.about = systray.AddMenuItem(msg.Sprintf("About"), "")
	m.don = systray.AddMenuItem(msg.Sprintf("Donations"), "")
	systray.AddSeparator()
	m.quit = systray.AddMenuItem(msg.Sprintf("Quit"), "")

	// Start handler
	go eventHandler(m, AppCfg, YD)
}

// menuHandler handles UI events
func eventHandler(m *menu, cfg *tools.Config, YD *ydisk.YDisk) {
	llog.Debug("event handler started")
	defer llog.Debug("event handler exited.")
	if cfg.StartDaemon {
		go YD.Start()
	}
	defer func() {
		if cfg.StopDaemon {
			YD.Stop()
		}
		YD.Close()
		systray.Quit()
	}()

	// Prepare the staff for icon animation
	icon := Icon{
		current: 0,
		Tick:    time.NewTimer(333 * time.Millisecond),
	}
	defer icon.Tick.Stop()
	currentStatus := ""
loop:
	for {
		select {
		case <-m.lastM[0].ClickedCh:
			tools.XdgOpen(m.lastP[0])
		case <-m.lastM[1].ClickedCh:
			tools.XdgOpen(m.lastP[1])
		case <-m.lastM[2].ClickedCh:
			tools.XdgOpen(m.lastP[2])
		case <-m.lastM[3].ClickedCh:
			tools.XdgOpen(m.lastP[3])
		case <-m.lastM[4].ClickedCh:
			tools.XdgOpen(m.lastP[4])
		case <-m.lastM[5].ClickedCh:
			tools.XdgOpen(m.lastP[5])
		case <-m.lastM[6].ClickedCh:
			tools.XdgOpen(m.lastP[6])
		case <-m.lastM[7].ClickedCh:
			tools.XdgOpen(m.lastP[7])
		case <-m.lastM[8].ClickedCh:
			tools.XdgOpen(m.lastP[8])
		case <-m.lastM[9].ClickedCh:
			tools.XdgOpen(m.lastP[9])
		case <-m.start.ClickedCh:
			go YD.Start()
		case <-m.stop.ClickedCh:
			go YD.Stop()
		case <-m.out.ClickedCh:
			notifySend(msg.Sprintf("Yandex.Disk daemon output"), YD.Output())
		case <-m.path.ClickedCh:
			tools.XdgOpen(YD.Path)
		case <-m.site.ClickedCh:
			tools.XdgOpen("https://disk.yandex.com")
		case <-m.help.ClickedCh:
			tools.XdgOpen("https://github.com/slytomcat/YD.go/wiki/FAQ&SUPPORT")
		case <-m.about.ClickedCh:
			notifySend("yd-go", msg.Sprintf(about, version, time.Now().Format("2006")))
		case <-m.don.ClickedCh:
			tools.XdgOpen("https://github.com/slytomcat/yd-go/wiki/Donations")
		case <-m.quit.ClickedCh:
			llog.Debug("Exit requested.")
			break loop
		case <-icon.Tick.C: //  Icon timer event
			if currentStatus == "busy" || currentStatus == "index" {
				systray.SetIcon(icons.BusyIcons[icon.Next()])
			}
		case yds := <-YD.Changes: // YDisk change event
			currentStatus = yds.Stat
			updateMenu(m, yds, icon, cfg.Notifications, YD.Path)
		}
	}
}

func updateMenu(m *menu, yds ydisk.YDvals, icon Icon, note bool, path string) {
	st := strings.Join([]string{statusTr[yds.Stat], yds.Prog, yds.Err, tools.ShortName(yds.ErrP, 30)}, " ")
	m.status.SetTitle(msg.Sprintf("Status: %s", st))
	if yds.Stat == "error" {
		m.status.SetTooltip(fmt.Sprintf("%s\nPath: %s", yds.Err, yds.ErrP))
	} else {
		m.status.SetTooltip("")
	}
	m.size1.SetTitle(msg.Sprintf("Used: %s/%s", yds.Used, yds.Total))
	m.size2.SetTitle(msg.Sprintf("Free: %s Trash: %s", yds.Free, yds.Trash))
	if yds.ChLast { // last synchronized list changed
		for i, p := range yds.Last {
			m.lastP[i] = filepath.Join(path, p)
			m.lastM[i].SetTitle(tools.ShortName(p, 40))
			if tools.NotExists(m.lastP[i]) {
				m.lastM[i].Disable()
			} else {
				m.lastM[i].Enable()
			}
			m.lastM[i].Show() // show list items
		}
		for i := len(yds.Last); i < 10; i++ {
			m.lastM[i].Hide() // hide the rest of list
		}
		if len(yds.Last) == 0 {
			m.last.Disable()
		} else {
			m.last.Enable()
		}
		llog.Debug("Last synchronized length:", len(yds.Last))
	}
	if yds.Stat != yds.Prev { // status changed
		// change indicator icon
		switch yds.Stat {
		case "idle":
			systray.SetIcon(icons.IdleIcon)
		case "busy", "index":
			systray.SetIcon(icons.BusyIcons[icon.Next()])
		case "none", "paused":
			systray.SetIcon(icons.PauseIcon)
		default:
			systray.SetIcon(icons.ErrorIcon)
		}
		// handle Start/Stop menu items
		if yds.Stat == "none" {
			m.start.Show()
			m.stop.Hide()
			m.out.Disable()
		} else {
			m.stop.Show()
			m.start.Hide()
			m.out.Enable()
		}
		if note { // handle notifications
			go handleNotifications(yds)
		}
	}
	llog.Debug("Change handled")
}

func handleNotifications(yds ydisk.YDvals) {
	switch {
	case yds.Stat == "none" && yds.Prev != "unknown":
		notifySend("Yandex.Disk", msg.Sprintf("Daemon stopped"))
	case yds.Prev == "none":
		notifySend("Yandex.Disk", msg.Sprintf("Daemon started"))
	case (yds.Stat == "busy" || yds.Stat == "index") &&
		(yds.Prev != "busy" && yds.Prev != "index"):
		notifySend("Yandex.Disk", msg.Sprintf("Synchronization started"))
	case (yds.Stat == "idle" || yds.Stat == "error") &&
		(yds.Prev == "busy" || yds.Prev == "index"):
		notifySend("Yandex.Disk", msg.Sprintf("Synchronization finished"))
	}
}

func onExit() {
	icons.CleanUp()
	llog.Debug("All done. Bye!")
}
