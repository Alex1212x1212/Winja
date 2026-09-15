package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"winja/ui"
)

func generateIPCSecret() string {
	appData, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	secretDir := filepath.Join(appData, "Winja")
	if err := os.MkdirAll(secretDir, 0700); err != nil {
		return ""
	}

	secretFile := filepath.Join(secretDir, "ipc_secret")
	if data, err := os.ReadFile(secretFile); err == nil {
		secret := strings.TrimSpace(string(data))
		if secret != "" {
			return secret
		}
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	secret := base64.RawURLEncoding.EncodeToString(buf)
	if err := os.WriteFile(secretFile, []byte(secret), 0600); err != nil {
		return ""
	}
	return secret
}

func validateIPCRequest(payload, secret string) (string, bool) {
	msg := strings.TrimSpace(payload)
	parts := strings.SplitN(msg, "\n", 2)
	if len(parts) != 2 {
		return "", false
	}
	if parts[0] != secret {
		return "", false
	}
	path := strings.TrimSpace(parts[1])
	if path == "" {
		return "", false
	}
	return path, true
}

func main() {
	var initialFileToScan string
	os.Setenv("FYNE_THEME", "dark")
	os.Setenv("FYNE_FONT", "Oswald.ttf")
	ipcSecret := generateIPCSecret()
	
	defer func() {
		if r := recover(); r != nil {
			os.WriteFile("crash.log", []byte(fmt.Sprint(r)), 0644)
		}
	}()

	if len(os.Args) > 1 {
		if ipcSecret != "" {
			conn, err := net.Dial("tcp", "127.0.0.1:49201")
			if err == nil {
				_, _ = conn.Write([]byte(ipcSecret + "\n" + os.Args[1]))
				conn.Close()
				os.Exit(0)
			}
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

	if ipcSecret != "" {
		// Start IPC server with a per-user secret so only this app instance can receive file requests.
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
				buf, err := io.ReadAll(conn)
				conn.Close()
				if err != nil || len(buf) == 0 {
					continue
				}
				if path, ok := validateIPCRequest(string(buf), ipcSecret); ok {
					mainApp.StartScan(path)
				}
			}
		}()
	}

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

















