package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"winja/core"
)

// Helper to create a check with explicit text color
func newCustomCheck(label string, checked bool, onChanged func(bool)) (*widget.Check, *canvas.Text, fyne.CanvasObject) {
	chk := widget.NewCheck("", onChanged)
	chk.SetChecked(checked)
	txt := canvas.NewText(label, color.NRGBA{100, 100, 100, 255})
	txt.TextSize = 13
	
	// Clickable wrapper over text
	btn := widget.NewButton("", func() {
		if !chk.Disabled() {
			chk.SetChecked(!chk.Checked)
		}
	})
	btn.Importance = widget.LowImportance
	textStack := container.NewStack(container.NewCenter(txt), btn)
	
	return chk, txt, container.NewHBox(chk, textStack)
}

func (app *WinjaApp) buildWebcamSettingsFrame() fyne.CanvasObject {
	chkActive, _, rowActive := newCustomCheck("Active Webcam usage detection (Monitor process using any camera devices)", core.IsWebcamActive(), func(checked bool) { core.SetWebcamActive(checked) })
	
	chkNotifStart, txtNotifStart, rowNotifStart := newCustomCheck("Display a notification as soon as a process uses a webcam device", core.IsWebcamNotifStart(), func(checked bool) { core.SetWebcamNotifStart(checked) })
	chkNotifStop, txtNotifStop, rowNotifStop := newCustomCheck("Display a notification as soon as a process stops using a camera device", core.IsWebcamNotifStop(), func(checked bool) { core.SetWebcamNotifStop(checked) })
	chkReqRep, txtReqRep, rowReqRep := newCustomCheck("As soon as a process starts using a camera device, request a report from VirusTotal", core.IsWebcamRequestReport(), func(checked bool) { core.SetWebcamRequestReport(checked) })
	
	lblAction := canvas.NewText("Action attached to the notification button (Only if notifications are used)", color.NRGBA{140, 140, 140, 255})
	lblAction.TextSize = 13
	
	// Custom radio group
	radActionVal := core.GetWebcamAction()
	
	rad1 := widget.NewRadioGroup([]string{""}, nil)
	txtRad1 := canvas.NewText("No action (Action button will not be present in the desktop notification)", color.NRGBA{100, 100, 100, 255})
	txtRad1.TextSize = 13
	btnRad1 := widget.NewButton("", func() {
		if !rad1.Disabled() {
			core.SetWebcamAction(0)
			rad1.SetSelected("")
			// We need a second radio to clear
		}
	})
	btnRad1.Importance = widget.LowImportance
	
	rad2 := widget.NewRadioGroup([]string{""}, nil)
	txtRad2 := canvas.NewText("Attempt to close the process detected as using a webcam device", color.NRGBA{100, 100, 100, 255})
	txtRad2.TextSize = 13
	btnRad2 := widget.NewButton("", func() {
		if !rad2.Disabled() {
			core.SetWebcamAction(1)
			rad2.SetSelected("")
		}
	})
	btnRad2.Importance = widget.LowImportance

	// Link the two radios manually
	rad1.OnChanged = func(s string) {
		if s != "" {
			core.SetWebcamAction(0)
			rad2.SetSelected("UNSELECTED") // fake value to clear
		}
	}
	rad2.OnChanged = func(s string) {
		if s != "" {
			core.SetWebcamAction(1)
			rad1.SetSelected("UNSELECTED")
		}
	}
	
	if radActionVal == 0 {
		rad1.SetSelected("")
	} else {
		rad2.SetSelected("")
	}
	
	rowRad1 := container.NewHBox(rad1, container.NewStack(container.NewCenter(txtRad1), btnRad1))
	rowRad2 := container.NewHBox(rad2, container.NewStack(container.NewCenter(txtRad2), btnRad2))
	
	updateState := func() {
		if chkActive.Checked {
			chkNotifStart.Enable()
			txtNotifStart.Color = color.NRGBA{100, 100, 100, 255}
			chkNotifStop.Enable()
			txtNotifStop.Color = color.NRGBA{100, 100, 100, 255}
			chkReqRep.Enable()
			txtReqRep.Color = color.NRGBA{100, 100, 100, 255}
			rad1.Enable()
			txtRad1.Color = color.NRGBA{100, 100, 100, 255}
			rad2.Enable()
			txtRad2.Color = color.NRGBA{100, 100, 100, 255}
		} else {
			chkNotifStart.Disable()
			txtNotifStart.Color = color.NRGBA{200, 200, 200, 255}
			chkNotifStop.Disable()
			txtNotifStop.Color = color.NRGBA{200, 200, 200, 255}
			chkReqRep.Disable()
			txtReqRep.Color = color.NRGBA{200, 200, 200, 255}
			rad1.Disable()
			txtRad1.Color = color.NRGBA{200, 200, 200, 255}
			rad2.Disable()
			txtRad2.Color = color.NRGBA{200, 200, 200, 255}
		}
		txtNotifStart.Refresh()
		txtNotifStop.Refresh()
		txtReqRep.Refresh()
		txtRad1.Refresh()
		txtRad2.Refresh()
	}

	chkActive.OnChanged = func(checked bool) {
		core.SetWebcamActive(checked)
		updateState()
	}
	updateState()

	// Separator lines
	makeSeparator := func() fyne.CanvasObject {
		rect := canvas.NewRectangle(color.NRGBA{235, 235, 235, 255})
		rect.SetMinSize(fyne.NewSize(1, 1))
		return rect
	}

	indentSize := fyne.NewSize(40, 1)

	contentLayout := container.NewVBox(
		container.NewGridWrap(fyne.NewSize(1, 20)),
		container.NewHBox(container.NewGridWrap(fyne.NewSize(20, 1)), rowActive),
		container.NewGridWrap(fyne.NewSize(1, 5)),
		makeSeparator(),
		container.NewGridWrap(fyne.NewSize(1, 5)),
		
		container.NewHBox(container.NewGridWrap(indentSize), rowNotifStart),
		container.NewHBox(container.NewGridWrap(indentSize), rowNotifStop),
		container.NewHBox(container.NewGridWrap(indentSize), rowReqRep),
		
		container.NewGridWrap(fyne.NewSize(1, 20)),
		container.NewHBox(container.NewGridWrap(indentSize), lblAction),
		container.NewGridWrap(fyne.NewSize(1, 5)),
		makeSeparator(),
		container.NewGridWrap(fyne.NewSize(1, 5)),
		
		container.NewHBox(container.NewGridWrap(indentSize), rowRad1),
		container.NewHBox(container.NewGridWrap(indentSize), rowRad2),
		
		layout.NewSpacer(),
	)
	
	mainBg := canvas.NewRectangle(color.White)
	
	themedContent := container.NewThemeOverride(
		container.NewStack(mainBg, container.NewVScroll(contentLayout)),
		theme.LightTheme(),
	)
	
	return themedContent
}

