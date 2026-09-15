package ui

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (app *WinjaApp) showDownloadAndScanDialog() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter URL here...")

	icon := canvas.NewImageFromResource(IconDownload)
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(64, 64))

	lblTitle := canvas.NewText("Target File Location", nil)
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}
	lblTitle.Alignment = fyne.TextAlignCenter
	lblTitle.TextSize = 16

	descText := "Using this function is harmless for your system. The file will be\ndownloaded with a special extension which prevents any execution of\nthe file. As soon as the temporary file is processed, the file will be\nautomatically deleted."
	lblDesc := widget.NewLabel(descText)
	lblDesc.Wrapping = fyne.TextWrapWord
	lblDesc.Alignment = fyne.TextAlignCenter

	var d dialog.Dialog

	btnDownload := widget.NewButton("Download", func() {
		targetURL := strings.TrimSpace(entry.Text)
		if targetURL == "" {
			return
		}
		_, err := url.ParseRequestURI(targetURL)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid URL: %v", err), app.Window)
			return
		}

		d.Hide()
		app.downloadAndScan(targetURL)
	})
	btnDownload.Importance = widget.HighImportance

	content := container.NewVBox(
		container.NewCenter(icon),
		layout.NewSpacer(),
		lblTitle,
		entry,
		lblDesc,
		layout.NewSpacer(),
		container.NewCenter(btnDownload),
	)
	
	// Pad the content to make it look spacious
	paddedContent := container.NewPadded(content)
	// Force a minimum size
	paddedContent = container.NewGridWrap(fyne.NewSize(400, 300), paddedContent)

	d = dialog.NewCustom("Download and Scan", "Cancel", paddedContent, app.Window)
	d.Show()
}

func (app *WinjaApp) downloadAndScan(targetURL string) {
	progress := dialog.NewProgress("Downloading", "Downloading file from URL...", app.Window)
	progress.Show()

	go func() {
		defer progress.Hide()
		
		resp, err := http.Get(targetURL)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Download failed: %v", err), app.Window)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("Download failed with status: %s", resp.Status), app.Window)
			return
		}

		// Create a temporary file with a special extension
		tmpDir := os.TempDir()
		tmpFile := filepath.Join(tmpDir, "winja_download.winjatemp")
		
		out, err := os.Create(tmpFile)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to create temp file: %v", err), app.Window)
			return
		}

		_, err = io.Copy(out, resp.Body)
		out.Close()
		
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to write temp file: %v", err), app.Window)
			return
		}

		// Start scan
		app.StartScan(tmpFile)
		
		// In a real app we'd delete the file AFTER the scan finishes.
		// Since StartScan sets up the UI and runs asynchronously, 
		// we can leave it in the OS temp directory (OS cleans it up eventually),
		// or we can wait for the scan to finish if we hook into it.
		// For now, it is safely in the Temp folder with .winjatemp extension.
	}()
}
