package ui

import (
	"fyne.io/fyne/v2/theme"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"winja/backend"
	"winja/core"
)

func newTruncatedLabel(text string) *widget.Label { 
	lbl := widget.NewLabel(text)
	lbl.Truncation = fyne.TextTruncateEllipsis
	return lbl 
}

func (app *WinjaApp) StartScan(filePath string) {
	app.SwitchToMain()
	AppTitleLabel.Text = "Running Tasks"
	AppTitleLabel.Refresh()
	app.ShowFrame(app.buildTasksFrame())

	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		entries, _ := os.ReadDir(filePath)
		// Process up to a reasonable limit to prevent total memory/API exhaustion
		count := 0
		for _, entry := range entries {
			if !entry.IsDir() {
				app.enqueueScan(filepath.Join(filePath, entry.Name()), false)
				count++
				if count > 50 { // Safety limit for folder scanning
					break
				}
			}
		}
		return
	}
	app.enqueueScan(filePath, true)
}

func (app *WinjaApp) enqueueScan(filePath string, singleFile bool) {
	now := time.Now().Format("15:04:05")
	filename := filepath.Base(filePath)

	statusLbl := widget.NewLabel("Scanning...")

	row := container.NewGridWithColumns(4,
		newTruncatedLabel(now),
		newTruncatedLabel(filename),
		statusLbl,
		newTruncatedLabel(filePath),
	)

	
	app.tasksBox.Add(container.NewThemeOverride(row, theme.LightTheme()))
	app.tasksBox.Refresh()

	go func() {
		backend.PerformScan(filePath, core.GetAPIKey(), core.IsDoNotUpload(), func(res backend.ScanResult) {
			statusLbl.SetText(res.Status)
			
			if res.Status != "Uploading..." && res.Status != "Queued / Analyzing..." {
				if !core.IsPrivateMode() {
					_, hash1, _ := backend.GetFileHashes(filePath)
					backend.LogReportHistory(filePath, res, hash1)
					app.populateReportsBox()
				}
			}
		})
	}()
}


