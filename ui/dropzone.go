package ui

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var dropZoneWin fyne.Window

func ShowDropTarget(app fyne.App, mainWin fyne.Window, mainApp *WinjaApp) {
	if dropZoneWin != nil {
		dropZoneWin.Show()
		return
	}

	drv, ok := app.Driver().(desktop.Driver)
	if !ok {
		return
	}

	dropZoneWin = drv.CreateSplashWindow()
	dropZoneWin.Resize(fyne.NewSize(150, 150))
	
	// Custom Title Bar for Drop Zone
	bg := canvas.NewRectangle(color.NRGBA{R: 12, G: 16, B: 22, A: 255})
	bg.SetMinSize(fyne.NewSize(0, 30))
	
	winjaBg := canvas.NewRectangle(color.NRGBA{R: 0, G: 200, B: 255, A: 255})
	winjaBg.SetMinSize(fyne.NewSize(50, 30))
	winjaText := canvas.NewText("WINJA", color.Black)
	winjaText.TextSize = 10
	winjaText.TextStyle = fyne.TextStyle{Bold: true}
	winjaText.Alignment = fyne.TextAlignCenter
	winjaBlock := container.NewStack(winjaBg, container.NewCenter(winjaText))

	// Buttons
	btnMax := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		dropZoneWin.Hide()
		mainWin.Show()
	})
	btnClose := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		// Close drop zone but keep app running if tray is active, or close completely?
		// Usually close drop zone just hides it, or we quit if main is hidden.
		app.Quit() 
	})
	
	// Make buttons smaller if possible
	controls := container.NewHBox(btnMax, btnClose)

	drag := NewDragArea()
	titleBar := container.NewStack(bg, container.NewBorder(nil, nil, winjaBlock, controls, drag))

	// Content area
	contentBg := canvas.NewRectangle(color.White)
	
	// The dashed outline of a star
	starIcon := canvas.NewImageFromResource(IconStar)
	starIcon.FillMode = canvas.ImageFillContain
	starIcon.SetMinSize(fyne.NewSize(80, 80))
	
	// Make the icon semi-transparent to look like a placeholder
	starIcon.Translucency = 0.8
	
	content := container.NewStack(contentBg, container.NewCenter(starIcon))

	dropZoneWin.SetContent(container.NewBorder(titleBar, nil, nil, nil, content))

	dropZoneWin.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			path := uris[0].Path()
			// Send it to the main app!
			mainApp.StartScan(path)
			// Switch to main window automatically?
			dropZoneWin.Hide()
			mainWin.Show()
			mainApp.ShowFrame(mainApp.buildTasksFrame())
			// Must update sidebar state visually if needed, but this works
		}
	})

	dropZoneWin.Show()
}
