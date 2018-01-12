package icons

import (
	"path"
)

var (
	IconBusy   [5]string
	IconError  string
	IconIdle   string
	IconPause  string
	IconNotify string
)

func SetTheme(appHome, theme string) {

	themePath := path.Join(appHome, theme)

	IconNotify = path.Join(appHome, "yd-128.png")

	IconBusy = [5]string{
		path.Join(themePath, "busy1.png"),
		path.Join(themePath, "busy2.png"),
		path.Join(themePath, "busy3.png"),
		path.Join(themePath, "busy4.png"),
		path.Join(themePath, "busy5.png"),
	}
	IconError = path.Join(themePath, "error.png")
	IconIdle = path.Join(themePath, "idle.png")
	IconPause = path.Join(themePath, "pause.png")
}
