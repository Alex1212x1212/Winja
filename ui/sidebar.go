package ui

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SidebarButton struct {
	widget.BaseWidget
	Icon     fyne.Resource
	OnTapped func()
	Active   bool
	bg       *canvas.Rectangle
	activeLine *canvas.Rectangle
	icon     *widget.Icon
}

func NewSidebarButton(icon fyne.Resource, tapped func()) *SidebarButton {
	b := &SidebarButton{Icon: icon, OnTapped: tapped}
	b.bg = canvas.NewRectangle(color.Transparent)
	b.bg.SetMinSize(fyne.NewSize(75, 60)) // Make buttons big squares (80x80)
	b.activeLine = canvas.NewRectangle(color.NRGBA{0, 200, 255, 255})
	b.activeLine.Hide()
		b.bg.FillColor = color.Transparent
	b.icon = widget.NewIcon(icon)
	b.ExtendBaseWidget(b)
	return b
}

func (b *SidebarButton) CreateRenderer() fyne.WidgetRenderer {
	// activeLine in left of Border stretches full height!
	b.activeLine.SetMinSize(fyne.NewSize(4, 10)) // just need width now, height stretches
	c := container.NewStack(b.bg, container.NewBorder(nil, nil, b.activeLine, nil, container.NewCenter(container.NewGridWrap(fyne.NewSize(28, 28), b.icon))))
	return widget.NewSimpleRenderer(c)
}

func (b *SidebarButton) Tapped(pe *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *SidebarButton) SetActive(active bool) {
	b.Active = active
	if active {
		b.activeLine.Show()
		b.bg.FillColor = color.NRGBA{R: 12, G: 16, B: 22, A: 255} // Dark background for active tab
	} else {
		b.activeLine.Hide()
		b.bg.FillColor = color.Transparent
	}
	b.activeLine.Refresh()
		b.bg.Refresh()
}

func (app *WinjaApp) buildMainSidebar() {
	var buttons []*SidebarButton

	setActive := func(activeBtn *SidebarButton) {
		for _, b := range buttons {
			b.SetActive(b == activeBtn)
		}
	}

	btnHome := NewSidebarButton(IconStar, nil)
	btnTasks := NewSidebarButton(IconHourglass, nil)
	btnReports := NewSidebarButton(IconUser, nil)
	btnSettings := NewSidebarButton(IconSettings, nil)
	btnHeart := NewSidebarButton(IconHeart, nil)

	btnHome.OnTapped = func() { setActive(btnHome); AppTitleLabel.Text = ""; AppTitleLabel.Refresh(); app.ShowFrame(app.buildHomeFrame()) }
	btnTasks.OnTapped = func() { setActive(btnTasks); AppTitleLabel.Text = "Running Tasks"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildTasksFrame()) }
	btnReports.OnTapped = func() { setActive(btnReports); AppTitleLabel.Text = "Report History"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildReportsFrame()) }
	btnSettings.OnTapped = func() { app.SwitchToSettings() }

	buttons = append(buttons, btnHome, btnTasks, btnReports, btnSettings, btnHeart)
	setActive(btnHome) // default active

	bg := canvas.NewRectangle(color.NRGBA{R: 12, G: 16, B: 22, A: 255})
	bg.SetMinSize(fyne.NewSize(75, 0)) // Enforce sidebar width to match WINJA block
	
	// The rest of the buttons in a VBox
	vbox := container.NewVBox(
		btnHome,
		btnTasks,
		btnReports,
		btnSettings,
		btnHeart,
		layout.NewSpacer(),
	)

	// Since we moved WINJA block to titlebar, just stack vbox on bg
	// Add some padding at the top if needed? The original sidebar icons
	// start directly under the title bar, so no padding needed.
	app.mainSidebar = container.NewStack(bg, container.NewPadded(vbox))
}

func (app *WinjaApp) buildSettingsSidebar() {
	var buttons []*SidebarButton

	setActive := func(activeBtn *SidebarButton) {
		for _, b := range buttons {
			b.SetActive(b == activeBtn)
		}
	}

	btnBack := NewSidebarButton(IconStar, nil)
	btnGenSet := NewSidebarButton(IconToolsC, nil)
	btnNetSet := NewSidebarButton(IconMonitor, nil)
	btnNotifSet := NewSidebarButton(IconChat, nil)
	btnWebcamSet := NewSidebarButton(IconTarget, nil)
	btnReportHist := NewSidebarButton(IconClock, nil)
	btnHeart := NewSidebarButton(IconHeart, nil)

	btnBack.OnTapped = func() { app.SwitchToMain() }
	btnGenSet.OnTapped = func() { setActive(btnGenSet); AppTitleLabel.Text = "General"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildGeneralSettingsFrame()) }
	btnNetSet.OnTapped = func() { setActive(btnNetSet); AppTitleLabel.Text = "Network Configuration"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildNetworkSettingsFrame()) }
	btnNotifSet.OnTapped = func() { setActive(btnNotifSet); AppTitleLabel.Text = "Notification / Widget"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildNotificationSettingsFrame()) }
	// webcam and reports can just show placeholders for now if not implemented
	btnWebcamSet.OnTapped = func() { setActive(btnWebcamSet); AppTitleLabel.Text = "Webcam Monitor"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildWebcamSettingsFrame()) }
	btnReportHist.OnTapped = func() { setActive(btnReportHist); AppTitleLabel.Text = "Report History"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildReportsFrame()) }

	buttons = append(buttons, btnBack, btnGenSet, btnNetSet, btnNotifSet, btnWebcamSet, btnReportHist, btnHeart)
	setActive(btnGenSet)

	bg := canvas.NewRectangle(color.NRGBA{R: 12, G: 16, B: 22, A: 255})
	vbox := container.NewVBox(btnBack, btnGenSet, btnNetSet, btnNotifSet, btnWebcamSet, btnReportHist, layout.NewSpacer(), btnHeart)
	app.settingsSidebar = container.NewStack(bg, vbox)
}











