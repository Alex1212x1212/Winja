package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	
	"fyne.io/fyne/v2/driver/desktop"
	"winja/ui"
)

func main() {
	var initialFileToScan string
	os.Setenv("FYNE_THEME", "dark")
	os.Setenv("FYNE_FONT", "Oswald.ttf")
	
	defer func() {
		if r := recover(); r != nil {
			os.WriteFile("crash.log", []byte(fmt.Sprint(r)), 0644)
		}
	}()

	if len(os.Args) > 1 {
		// If trying to open a file, send to running instance if possible
		conn, err := net.Dial("tcp", "127.0.0.1:49201")
		if err == nil {
			conn.Write([]byte(os.Args[1]))
			conn.Close()
			os.Exit(0)
		}
		initialFileToScan = os.Args[1]
	}

	a := app.NewWithID("com.winja.clone")
	var w fyne.Window
	if drv, ok := a.Driver().(desktop.Driver); ok {
		w = drv.CreateSplashWindow()
	} else {
		w = a.NewWindow("Winja")
	}
	w.Resize(fyne.NewSize(1024, 600))

	basePath, err := os.Getwd()
	if err != nil {
		ex, _ := os.Executable()
		basePath = filepath.Dir(ex)
	}
	ui.LoadResources(basePath)

	logoRes, err := fyne.LoadResourceFromPath(filepath.Join(basePath, "assets", "logo.png"))
	if err == nil {
		a.SetIcon(logoRes)
		w.SetIcon(logoRes)
	}

	mainApp := ui.NewApp(a, w)

	// Start IPC server
	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:49201")
		if err != nil {
			fmt.Println("IPC Listen error:", err)
			return
		}
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			buf := make([]byte, 4096)
			n, _ := conn.Read(buf)
			if n > 0 {
				path := string(buf[:n])
				mainApp.StartScan(path)
			}
			conn.Close()
		}
	}()

	w.SetContent(mainApp.BuildContent())

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Winja",
			fyne.NewMenuItem("Show Winja", func() { w.Show() }),
			fyne.NewMenuItem("Quit Winja", func() { a.Quit() }),
		)
		desk.SetSystemTrayMenu(m)
		if logoRes != nil { desk.SetSystemTrayIcon(logoRes) }
	}
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	w.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			path := uris[0].Path()
			mainApp.StartScan(path)
		}
	})

	if initialFileToScan != "" {
		go mainApp.StartScan(initialFileToScan)
	}
	
	w.ShowAndRun()
}

















