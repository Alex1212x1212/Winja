package main

import (
	"io/ioutil"
	"fyne.io/fyne/v2"
	"github.com/getlantern/systray"
)

func SetupCustomTray(a fyne.App, w fyne.Window) {
	go systray.Run(func() {
		b, err := ioutil.ReadFile("assets/logo.png")
		if err == nil {
			systray.SetIcon(b)
		}
		systray.SetTooltip("Winja")
		mShow := systray.AddMenuItem("Show Winja", "")
		mQuit := systray.AddMenuItem("Quit Winja", "")
		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					w.Show()
				case <-mQuit.ClickedCh:
					systray.Quit()
					a.Quit()
				}
			}
		}()
	}, func() {})
}
