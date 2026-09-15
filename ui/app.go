package ui

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type WinjaApp struct {
	App    fyne.App
	Window fyne.Window

	// Sidebars
	mainSidebar     *fyne.Container
	sidebarContainer *fyne.Container
	settingsSidebar *fyne.Container// Content area
	contentArea *fyne.Container

	// Task container
	tasksBox *fyne.Container
	reportsBox *fyne.Container
}

func NewApp(a fyne.App, w fyne.Window) *WinjaApp {
	return &WinjaApp{
		App:    a,
		Window: w,
		tasksBox: container.NewVBox(),
		reportsBox: container.NewVBox(),
	} 
}

func (app *WinjaApp) BuildContent() fyne.CanvasObject {
	app.buildMainSidebar()
	app.buildSettingsSidebar()

	// Initial content is home
	app.contentArea = container.NewMax(app.buildHomeFrame())

	// Split sidebars and content
	app.sidebarContainer = container.NewMax(app.mainSidebar) // start with main

	titleBar := BuildCustomTitleBar(app)
	
	// No gap, flush edges
	mainLayout := container.NewBorder(titleBar, nil, app.sidebarContainer, nil, app.contentArea)

	borderTop := NewResizeBorder(app, "top")
	borderBottom := NewResizeBorder(app, "bottom")
	borderLeft := NewResizeBorder(app, "left")
	borderRight := NewResizeBorder(app, "right")

	return container.NewBorder(borderTop, borderBottom, borderLeft, borderRight, mainLayout)
}

func (app *WinjaApp) SwitchToSettings() {
	app.sidebarContainer.RemoveAll(); app.sidebarContainer.Add(app.settingsSidebar)
	AppTitleLabel.Text = "General"; AppTitleLabel.Refresh(); app.ShowFrame(app.buildGeneralSettingsFrame())
}


func (app *WinjaApp) SwitchToMain() {
	app.sidebarContainer.RemoveAll(); app.sidebarContainer.Add(app.mainSidebar); AppTitleLabel.Text = ""; AppTitleLabel.Refresh(); app.ShowFrame(app.buildHomeFrame())
	app.ShowFrame(app.buildHomeFrame())
}



func (app *WinjaApp) ShowFrame(frame fyne.CanvasObject) {
	app.contentArea.RemoveAll()
	app.contentArea.Add(frame)
}































