package icons

import (
  "path"
)

var (
  AppHome string = "/home/stc/DEV/GO/src/github.com/slytomcat/YD.go"
  IconBusy [5]string
  IconError string
  IconIdle string
  IconPause string
)

func SetTheme(theme string) {

  themePath := path.Join(AppHome, "icons", theme)


  IconBusy = [5]string { path.Join(themePath, "busy1.png"),
                         path.Join(themePath, "busy2.png"),
                         path.Join(themePath, "busy3.png"),
                         path.Join(themePath, "busy4.png"),
                         path.Join(themePath, "busy5.png"),
                       }
  IconError = path.Join(themePath, "error.png")
  IconIdle = path.Join(themePath, "idle.png")
  IconPause = path.Join(themePath, "pause.png")
}

