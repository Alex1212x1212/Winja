package ui

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"winja/core"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (app *WinjaApp) buildGeneralSettingsFrame() fyne.CanvasObject {
	// Helper for section titles with icons
	sectionTitle := func(title string, icon fyne.Resource) fyne.CanvasObject {
		img := canvas.NewImageFromResource(icon)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(24, 24))
		lbl := canvas.NewText(title, color.NRGBA{120, 120, 120, 255})
		lbl.TextStyle = fyne.TextStyle{Bold: true}
		lbl.TextSize = 16
		return container.NewHBox(img, lbl)
	}

	// --- GENERAL SETTINGS ---
	titleGen := sectionTitle("GENERAL SETTINGS", IconSettings)
	chkStartWin := widget.NewCheck("Start Winja with Windows. (Register to system startup)", func(checked bool) { core.SetRunAtStartup(checked) })
	chkStartWin.SetChecked(core.IsRunAtStartup())
	chkRegExp := widget.NewCheck("Register Winja Explorer Extension (Contextual Scan)", func(checked bool) { core.SetExplorerExtension(checked) })
	chkRegExp.SetChecked(core.IsExplorerExtensionEnabled())
	chkMinTray := widget.NewCheck("Minimize Winja to System Tray (Little icon close to your clock)", func(checked bool) { core.SetMinimizeToTray(checked) })
	chkMinTray.SetChecked(core.IsMinimizeToTray())
	chkIncSize := widget.NewCheck("Increase the size of list items (Useful for touchscreen)", nil)

	groupGen := container.NewVBox(
		titleGen,
		container.NewPadded(container.NewVBox(chkStartWin, chkRegExp, chkMinTray, chkIncSize)),
		widget.NewSeparator(),
	)

	// --- FILE SCAN SETTINGS ---
	titleScan := sectionTitle("FILE SCAN SETTINGS", IconMonitor)
	chkPrivate := widget.NewCheck("Private Mode: Scan reports are not stored. They won't be listed in report history", func(checked bool) { core.SetPrivateMode(checked) })
	chkPrivate.SetChecked(core.IsPrivateMode())
	chkNoUp := widget.NewCheck("Do not upload unknown samples to VirusTotal Servers. (Only check if a fingerprint exists)", func(checked bool) { core.SetDoNotUpload(checked) })
	chkNoUp.SetChecked(core.IsDoNotUpload())
	chkOpenWeb := widget.NewCheck("As soon as a report is available, open VirusTotal website on your default browser", func(checked bool) { core.SetOpenWeb(checked) })
	chkOpenWeb.SetChecked(core.IsOpenWeb())
	
	opts := []string{
		"Always, whether or not a report mentions a clean or infected file (Recommended)",
		"Only when a report mentions a dangerous or infected file",
		"Only when a report mentions a clean file (Not recommended)",
	}
	
	radNotifOpts := widget.NewRadioGroup(opts, func(selected string) {
		for i, v := range opts {
			if v == selected {
				core.SetNotifLevel(i)
				break
			}
		}
	})
	
	level := core.GetNotifLevel()
	if level >= 0 && level < len(opts) {
		radNotifOpts.SetSelected(opts[level])
	}
	
	chkNotif := widget.NewCheck("Display a desktop notification when :", func(checked bool) {
		core.SetNotifEnabled(checked)
		if checked {
			radNotifOpts.Enable()
		} else {
			radNotifOpts.Disable()
		}
	})
	chkNotif.SetChecked(core.IsNotifEnabled())
	if !core.IsNotifEnabled() {
		radNotifOpts.Disable()
	}
	
	chkCritFail := widget.NewCheck("On any VirusTotal transaction critical failure (Ex: Network Error)", func(checked bool) { core.SetCritFail(checked) })
	chkCritFail.SetChecked(core.IsCritFail())
	
	notifOptsContainer := container.NewPadded(container.NewVBox(
		radNotifOpts,
		chkCritFail,
	))

	groupScan := container.NewVBox(
		titleScan,
		container.NewPadded(container.NewVBox(chkPrivate, chkNoUp, chkOpenWeb, chkNotif, notifOptsContainer)),
		widget.NewSeparator(),
	)

	// --- CUSTOM VIRUSTOTAL API KEY ---
	titleApi := sectionTitle("CUSTOM VIRUSTOTAL API KEY", IconToolsC)
	lblApiKey := widget.NewLabel("API Key :")
	entryApi := widget.NewEntry()
	entryApi.SetText(core.GetAPIKey())
	entryApi.OnChanged = func(s string) { core.SetAPIKey(s) }
	
	btnClear := widget.NewButton("X", func() {
		entryApi.SetText("")
	})
	
	btnHelp := widget.NewButton("?", func() {})
	
	// Layout for the API row
	entrySized := container.NewGridWrap(fyne.NewSize(500, 36), entryApi)
	apiInputContainer := container.NewHBox(entrySized, btnClear, btnHelp)
	apiRow := container.NewHBox(lblApiKey, apiInputContainer)
	
	lblApiDesc := widget.NewLabel("Turn off the upload size limitation (>32MB). Using your personal VirusTotal API v3 key allows you to upload files up to the 650 MB limit.")
	lblApiDesc.Wrapping = fyne.TextWrapWord
	
	groupApi := container.NewVBox(
		titleApi,
		container.NewPadded(container.NewVBox(apiRow, lblApiDesc)),
	)

	content := container.NewVBox(
		groupGen,
		groupScan,
		groupApi,
	)

	bg := canvas.NewRectangle(color.White)
	return container.NewStack(bg, container.NewThemeOverride(container.NewVScroll(container.NewPadded(content)), theme.LightTheme()))
}
















