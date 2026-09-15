package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func (app *WinjaApp) buildNetworkSettingsFrame() fyne.CanvasObject {
	// Header bar: Light grey background, icon, and title
	headerBg := canvas.NewRectangle(color.NRGBA{245, 245, 245, 255})
	
	iconMonitor := canvas.NewImageFromResource(IconMonitor)
	if iconMonitor.Resource == nil {
		iconMonitor = canvas.NewImageFromResource(theme.ComputerIcon())
	}
	iconMonitor.FillMode = canvas.ImageFillContain
	iconMonitor.SetMinSize(fyne.NewSize(24, 24))
	
	// Invert icon if it's white? Since we force light theme, Fyne icons auto-invert, 
	// but image resource might not. Assuming it's fine.
	
	lblTitle := canvas.NewText("CONFIGURE PROXY", color.NRGBA{120, 130, 140, 255})
	lblTitle.TextSize = 16
	lblTitle.TextStyle = fyne.TextStyle{Bold: true}
	
	headerHBox := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(20, 1)), // Left padding
		container.NewCenter(iconMonitor),
		container.NewGridWrap(fyne.NewSize(10, 1)), // gap
		container.NewCenter(lblTitle),
	)
	
	headerBar := container.NewStack(
		headerBg,
		container.NewPadded(headerHBox),
	)

	// Radio buttons
	radProxy := widget.NewRadioGroup([]string{
		"Please don't try to connect through any proxy server",
		"Use default settings (Recommended)",
		"Configure a proxy manually",
	}, nil)
	radProxy.SetSelected("Please don't try to connect through any proxy server")
	
	lblHost := canvas.NewText("Proxy Host / DNS :", color.NRGBA{160, 160, 160, 255})
	lblHost.TextSize = 14
	lblHost.TextStyle = fyne.TextStyle{Bold: true}
	
	entryHost := widget.NewEntry()
	// Create a soft grey rounded background for the entry
	entryHost.Wrapping = fyne.TextTruncate
	
	lblPort := canvas.NewText("Proxy Port :", color.NRGBA{160, 160, 160, 255})
	lblPort.TextSize = 14
	lblPort.TextStyle = fyne.TextStyle{Bold: true}
	
	entryPort := widget.NewEntry()
	
	// Constrain the width of the entries
	entryHostWrap := container.NewGridWrap(fyne.NewSize(300, 36), entryHost)
	entryPortWrap := container.NewGridWrap(fyne.NewSize(150, 36), entryPort)
	
	// Form layout: labels above entries
	formVBox := container.NewVBox(
		lblHost,
		entryHostWrap,
		container.NewGridWrap(fyne.NewSize(1, 10)), // spacing
		lblPort,
		entryPortWrap,
	)
	
	// Indent the radio buttons and form
	indentSize := fyne.NewSize(40, 1)
	
	contentLayout := container.NewVBox(
		headerBar,
		container.NewGridWrap(fyne.NewSize(1, 10)),
		container.NewHBox(container.NewGridWrap(indentSize), radProxy),
		container.NewGridWrap(fyne.NewSize(1, 10)),
		container.NewHBox(container.NewGridWrap(indentSize), formVBox),
		layout.NewSpacer(),
	)
	
	// Force the entire panel to be White
	mainBg := canvas.NewRectangle(color.White)
	
	// Wrap in LightTheme so the radio buttons and entries render correctly
	themedContent := container.NewThemeOverride(
		container.NewStack(mainBg, container.NewVScroll(contentLayout)),
		theme.LightTheme(),
	)
	
	return themedContent
}
