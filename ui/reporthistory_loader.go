package ui

import (
	"fyne.io/fyne/v2/theme"
	"sync"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"winja/backend"
)

var reportsMutex sync.Mutex

func (app *WinjaApp) populateReportsBox() {
	reportsMutex.Lock()
	defer reportsMutex.Unlock()
	
	app.reportsBox.Objects = nil // clear existing
	
	history := backend.GetReportHistory()
	for _, h := range history {
		shortHash := h.SHA1
		if len(shortHash) > 16 {
			shortHash = shortHash[:8] + "..." + shortHash[len(shortHash)-8:]
		}
		
		vtURL, _ := url.Parse("https://www.virustotal.com/gui/file/" + h.SHA1 + "/detection")
		hashLink := widget.NewHyperlink(shortHash, vtURL)
		hashLink.Truncation = fyne.TextTruncateEllipsis
		
		row := container.NewGridWithColumns(6,
			newTruncatedLabel(h.ReportDate),
			newTruncatedLabel(h.FileName),
			newTruncatedLabel(h.Status),
			newTruncatedLabel(h.DetectionRate),
			newTruncatedLabel(h.LastScannedDate),
			hashLink,
		)
		app.reportsBox.Add(container.NewThemeOverride(row, theme.LightTheme()))
	}
	app.reportsBox.Refresh()
}



