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

func buildHeaderBar(title string, iconRes fyne.Resource) *fyne.Container {
	bg := canvas.NewRectangle(color.NRGBA{245, 245, 245, 255})
	
	img := widget.NewIcon(iconRes)
	imgContainer := container.NewGridWrap(fyne.NewSize(24, 24), img)
	
	lbl := canvas.NewText(title, color.NRGBA{120, 130, 140, 255})
	lbl.TextSize = 16
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	
	hbox := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(20, 1)),
		container.NewCenter(imgContainer),
		container.NewGridWrap(fyne.NewSize(10, 1)),
		container.NewCenter(lbl),
	)
	
	return container.NewStack(bg, container.NewPadded(hbox))
}

func (app *WinjaApp) buildNotificationSettingsFrame() fyne.CanvasObject {
	// --- DESKTOP WIDGET ---
	headerWidget := buildHeaderBar("DESKTOP WIDGET", theme.ComputerIcon())
	
	chkFixPos, txtFixPos, rowFixPos := newCustomCheck("Widget position is fixed. You can't move the widget position.", core.IsWidgetFixed(), func(checked bool) { core.SetWidgetFixed(checked) })
	chkTopMost, txtTopMost, rowTopMost := newCustomCheck("Widget is always on top of other windows", core.IsWidgetTopMost(), func(checked bool) { core.SetWidgetTopMost(checked) })
	chkOpacity, txtOpacity, rowOpacity := newCustomCheck("Widget has opacity, you can see through the Widget window.", core.IsWidgetOpacity(), func(checked bool) { core.SetWidgetOpacity(checked) })
	
	updateState := func(checked bool) {
		if checked {
			chkFixPos.Enable(); txtFixPos.Color = color.NRGBA{100, 100, 100, 255}; txtFixPos.Refresh()
			chkTopMost.Enable(); txtTopMost.Color = color.NRGBA{100, 100, 100, 255}; txtTopMost.Refresh()
			chkOpacity.Enable(); txtOpacity.Color = color.NRGBA{100, 100, 100, 255}; txtOpacity.Refresh()
		} else {
			chkFixPos.Disable(); txtFixPos.Color = color.NRGBA{200, 200, 200, 255}; txtFixPos.Refresh()
			chkTopMost.Disable(); txtTopMost.Color = color.NRGBA{200, 200, 200, 255}; txtTopMost.Refresh()
			chkOpacity.Disable(); txtOpacity.Color = color.NRGBA{200, 200, 200, 255}; txtOpacity.Refresh()
		}
	}
	
	_, _, rowShowWidget := newCustomCheck("When Winja main window is minimized, display the Widget", core.IsWidgetShow(), func(checked bool) {
		core.SetWidgetShow(checked)
		updateState(checked)
	})
	
	updateState(core.IsWidgetShow())
	
	indentSize := fyne.NewSize(40, 1)
	subIndentSize := fyne.NewSize(60, 1)
	
	vboxWidget := container.NewVBox(
		headerWidget,
		container.NewGridWrap(fyne.NewSize(1, 10)),
		container.NewHBox(container.NewGridWrap(indentSize), rowShowWidget),
		container.NewHBox(container.NewGridWrap(subIndentSize), rowFixPos),
		container.NewHBox(container.NewGridWrap(subIndentSize), rowTopMost),
		container.NewHBox(container.NewGridWrap(subIndentSize), rowOpacity),
		container.NewGridWrap(fyne.NewSize(1, 20)),
	)

	// --- DESKTOP NOTIFICATIONS ---
	headerNotif := buildHeaderBar("DESKTOP NOTIFICATIONS", theme.MailComposeIcon())
	
	greyText := color.NRGBA{160, 160, 160, 255}
	
	_, _, rowSticky := newCustomCheck("Notifications are sticky. It doesn't vanish automatically after a little amount of time. You must manually close the notification.", core.IsNotifSticky(), func(checked bool) { core.SetNotifSticky(checked) })
	
	lblAutoClose := canvas.NewText("Close notification automatically after :", greyText)
	lblAutoClose.TextSize = 14
	lblAutoClose.TextStyle = fyne.TextStyle{Bold: true}
	
	entryTime := widget.NewEntry()
	entryTime.SetText(core.GetNotifCloseTime())
	entryTime.OnChanged = func(s string) { core.SetNotifCloseTime(s) }
	entryTimeWrap := container.NewGridWrap(fyne.NewSize(100, 36), entryTime)
	
	lblSec := canvas.NewText("Seconds (2 - 120)", greyText)
	lblSec.TextSize = 14
	
	rowTimeInputs := container.NewHBox(entryTimeWrap, lblSec)
	
	lblPlace := canvas.NewText("Notification placement on desktop", greyText)
	lblPlace.TextSize = 14
	lblPlace.TextStyle = fyne.TextStyle{Bold: true}
	
	selPlace := widget.NewSelect([]string{"Bottom Centered", "Top Right", "Bottom Right"}, func(s string) { core.SetNotifPlacement(s) })
	selPlace.SetSelected(core.GetNotifPlacement())
	selPlaceWrap := container.NewGridWrap(fyne.NewSize(200, 36), selPlace)
	
	btnTry := widget.NewButton("Try", func() {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("Test Notification", "This is a test notification from Winja!"))
	})
	
	rowPlaceInputs := container.NewHBox(selPlaceWrap, btnTry)
	
	vboxNotif := container.NewVBox(
		headerNotif,
		container.NewGridWrap(fyne.NewSize(1, 10)),
		container.NewHBox(container.NewGridWrap(indentSize), rowSticky),
		container.NewGridWrap(fyne.NewSize(1, 5)),
		container.NewHBox(container.NewGridWrap(indentSize), lblAutoClose),
		container.NewHBox(container.NewGridWrap(indentSize), rowTimeInputs),
		container.NewGridWrap(fyne.NewSize(1, 10)),
		container.NewHBox(container.NewGridWrap(indentSize), lblPlace),
		container.NewHBox(container.NewGridWrap(indentSize), rowPlaceInputs),
		layout.NewSpacer(),
	)
	
	contentLayout := container.NewVBox(vboxWidget, vboxNotif)
	mainBg := canvas.NewRectangle(color.White)
	
	themedContent := container.NewThemeOverride(
		container.NewStack(mainBg, container.NewVScroll(contentLayout)),
		theme.LightTheme(),
	)
	
	return themedContent
}
