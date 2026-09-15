package ui

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"winja/core"
)

type BigButton struct {
	widget.BaseWidget
	Title    string
	Icon     fyne.Resource
	OnTapped func()
	bg       *canvas.Rectangle
}

func NewBigButton(title string, icon fyne.Resource, tapped func()) *BigButton {
	b := &BigButton{Title: title, Icon: icon, OnTapped: tapped}
	b.bg = canvas.NewRectangle(color.NRGBA{250, 250, 250, 255})
	b.bg.CornerRadius = 15
	b.ExtendBaseWidget(b)
	return b
}

func (b *BigButton) MinSize() fyne.Size {
	return fyne.NewSize(280, 140)
}

func (b *BigButton) CreateRenderer() fyne.WidgetRenderer {
	// Use image instead of widget.Icon to preserve original colors!
	img := canvas.NewImageFromResource(b.Icon)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(48, 48))
	
	label := canvas.NewText(b.Title, color.NRGBA{100, 100, 100, 255})
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.TextSize = 13
	label.Alignment = fyne.TextAlignCenter
	
	c := container.NewStack(b.bg, container.NewCenter(container.NewVBox(img, layout.NewSpacer(), label)))
	return widget.NewSimpleRenderer(c)
}

func (b *BigButton) Tapped(pe *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (app *WinjaApp) buildHomeFrame() fyne.CanvasObject {
	headerTitle := canvas.NewText("WELCOME TO WINJA 7.1", color.White)
	headerTitle.TextSize = 26
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}
	
	headerSub := canvas.NewText("\"Once we accept our limits, we go beyond them\" ~Albert Einstein", color.White)
	headerSub.TextSize = 14
	
	logoImg := canvas.NewImageFromResource(IconStar)
	if logoRes, err := fyne.LoadResourceFromPath("assets/logo.png"); err == nil {
		logoImg = canvas.NewImageFromResource(logoRes)
	}
	logoImg.FillMode = canvas.ImageFillContain
	logoImg.SetMinSize(fyne.NewSize(80, 80))

	headerTextVBox := container.NewVBox(layout.NewSpacer(), headerTitle, headerSub, layout.NewSpacer())
	
	// Left-aligned with some padding
	indent := canvas.NewRectangle(color.Transparent)
	indent.SetMinSize(fyne.NewSize(30, 1))
	
	headerContent := container.NewHBox(
		indent,
		logoImg,
		indent,
		headerTextVBox,
	)

	headerBg := canvas.NewRectangle(color.NRGBA{R: 12, G: 16, B: 22, A: 255})
	headerBg.SetMinSize(fyne.NewSize(100, 110)) // Make it a bit taller
	
	header := container.NewStack(headerBg, container.NewPadded(headerContent))

	btnBrowse := NewBigButton("BROWSE FOR FILE", IconBrowse, func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err == nil && r != nil {
				path := r.URI().Path()
				app.StartScan(path)
			}
		}, app.Window)
	})
	btnDownload := NewBigButton("DOWNLOAD AND SCAN", IconDownload, func() { app.showDownloadAndScanDialog() })	

	grid := container.NewGridWithColumns(2,
		btnBrowse,
		btnDownload,
		
	)

	starCenter := canvas.NewImageFromResource(IconStar)
	starCenter.FillMode = canvas.ImageFillContain
	starCenter.SetMinSize(fyne.NewSize(24, 24))
	centerOverlay := container.NewCenter(starCenter)
	
	gridWithStar := container.NewStack(grid, centerOverlay)

	// Center content below header
	contentVBox := container.NewVBox(
		layout.NewSpacer(),
		container.NewGridWrap(fyne.NewSize(1, 20)), // Spacer
		container.NewCenter(gridWithStar),
		layout.NewSpacer(),
	)

	apiStatusText := "VirusTotal API: NOT CONFIGURED (Scanning Disabled)"
	apiStatusColor := color.NRGBA{R:200, G:50, B:50, A:255}
	if core.GetAPIKey() != "" {
		apiStatusText = "VirusTotal API: ACTIVE (Scanning Ready)"
		apiStatusColor = color.NRGBA{R:50, G:200, B:50, A:255}
	}
	apiStatusLbl := canvas.NewText(apiStatusText, apiStatusColor)
	apiStatusLbl.TextStyle = fyne.TextStyle{Bold: true}
	apiStatusLbl.Alignment = fyne.TextAlignCenter
	statusContainer := container.NewPadded(apiStatusLbl)
	
	contentStack := container.NewStack(contentVBox, container.NewVBox(layout.NewSpacer(), statusContainer))

	bg := canvas.NewRectangle(color.White) // Main background is WHITE

	return container.NewStack(bg, container.NewBorder(
		header, nil, nil, nil,
		contentStack,
	))
}


func (app *WinjaApp) buildTasksFrame() fyne.CanvasObject {
	// Table Headers
	lblTime := canvas.NewText("Request Time", color.NRGBA{80, 80, 80, 255})
	lblTime.TextStyle = fyne.TextStyle{Bold: true}
	lblName := canvas.NewText("Target Name", color.NRGBA{80, 80, 80, 255})
	lblName.TextStyle = fyne.TextStyle{Bold: true}
	lblStatus := canvas.NewText("Status", color.NRGBA{80, 80, 80, 255})
	lblStatus.TextStyle = fyne.TextStyle{Bold: true}
	lblLocation := canvas.NewText("Location", color.NRGBA{80, 80, 80, 255})
	lblLocation.TextStyle = fyne.TextStyle{Bold: true}

	headerRow := container.NewGridWithColumns(4,
		lblTime,
		lblName,
		lblStatus,
		lblLocation,
	)

	separator := canvas.NewRectangle(color.NRGBA{230, 230, 230, 255}) // Faint grey line
	separator.SetMinSize(fyne.NewSize(10, 1))

	

	scroll := container.NewVScroll(app.tasksBox)

	content := container.NewBorder(
		container.NewVBox(
			container.NewPadded(headerRow),
			separator,
		),
		nil, nil, nil,
		scroll,
	)
	
	// Ensure background is pure white
	bg := canvas.NewRectangle(color.White)
	return container.NewThemeOverride(container.NewStack(bg, content), theme.LightTheme())
}

func (app *WinjaApp) buildReportsFrame() fyne.CanvasObject {
	// Table Headers
	lblDate := canvas.NewText("Report Date", color.NRGBA{80, 80, 80, 255})
	lblDate.TextStyle = fyne.TextStyle{Bold: true}
	lblName := canvas.NewText("File Name", color.NRGBA{80, 80, 80, 255})
	lblName.TextStyle = fyne.TextStyle{Bold: true}
	lblStatus := canvas.NewText("Status", color.NRGBA{80, 80, 80, 255})
	lblStatus.TextStyle = fyne.TextStyle{Bold: true}
	lblRate := canvas.NewText("Detection Rate", color.NRGBA{80, 80, 80, 255})
	lblRate.TextStyle = fyne.TextStyle{Bold: true}
	lblScanDate := canvas.NewText("Last Scanned Date", color.NRGBA{80, 80, 80, 255})
	lblScanDate.TextStyle = fyne.TextStyle{Bold: true}
	lblSha := canvas.NewText("SHA1", color.NRGBA{80, 80, 80, 255})
	lblSha.TextStyle = fyne.TextStyle{Bold: true}

	headerRow := container.NewGridWithColumns(6,
		lblDate,
		lblName,
		lblStatus,
		lblRate,
		lblScanDate,
		lblSha,
	)

	// Custom separators matching screenshot: faint grey for first 5, cyan for last
	faintLine := func() fyne.CanvasObject {
		r := canvas.NewRectangle(color.NRGBA{230, 230, 230, 255})
		r.SetMinSize(fyne.NewSize(10, 2))
		return r
	}
	cyanLine := canvas.NewRectangle(color.NRGBA{0, 200, 255, 255})
	cyanLine.SetMinSize(fyne.NewSize(10, 2))

	separatorRow := container.NewGridWithColumns(6,
		faintLine(), faintLine(), faintLine(), faintLine(), faintLine(), cyanLine,
	)

	app.populateReportsBox()
	scroll := container.NewVScroll(app.reportsBox)

	content := container.NewBorder(
		container.NewVBox(
			container.NewPadded(headerRow),
			separatorRow,
		),
		nil, nil, nil,
		scroll,
	)
	
	bg := canvas.NewRectangle(color.White)
	return container.NewThemeOverride(container.NewStack(bg, content), theme.LightTheme())
}


















