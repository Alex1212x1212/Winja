package ui

import (
	"fyne.io/fyne/v2"
	"path/filepath"
)

var (
	IconStar      fyne.Resource
	IconHourglass fyne.Resource
	IconUser      fyne.Resource
	IconSettings  fyne.Resource
	IconHeart     fyne.Resource
	IconShield    fyne.Resource
	IconBrowse    fyne.Resource
	IconDownload  fyne.Resource
	IconRocket    fyne.Resource
	IconTools     fyne.Resource
	IconToolsC    fyne.Resource
	IconMonitor   fyne.Resource
	IconChat      fyne.Resource
	IconTarget    fyne.Resource
	IconClock     fyne.Resource
)

func LoadResources(basePath string) {
	load := func(name string) fyne.Resource {
		res, err := fyne.LoadResourceFromPath(filepath.Join(basePath, "assets", name))
		if err != nil {
			return nil // Fallback to nil if missing
		}
		return res
	}

	IconStar = load("star.png")
	IconHourglass = load("hourglass.png")
	IconUser = load("user.png")
	IconSettings = load("settings.png")
	IconHeart = load("heart.png")
	IconShield = load("shield.png")
	IconBrowse = load("browse.png")
	IconDownload = load("download.png")
	IconRocket = load("rocket.png")
	IconTools = load("tools.png")
	IconToolsC = load("tools_cross.png")
	IconMonitor = load("monitor.png")
	IconChat = load("chat.png")
	IconTarget = load("target.png")
	IconClock = load("clock.png")
}
