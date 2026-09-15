package ui

import (
	"image/color"
	"syscall"
	"winja/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procReleaseCapture = user32.NewProc("ReleaseCapture")
var procSendMessage = user32.NewProc("SendMessageW")
var procGetActiveWindow = user32.NewProc("GetActiveWindow")

// DragArea is a transparent widget that forwards mouse drags to the Windows OS
// to drag frameless windows.
type DragArea struct {
	widget.BaseWidget
}

func NewDragArea() *DragArea {
	d := &DragArea{}
	d.ExtendBaseWidget(d)
	return d
}

func (d *DragArea) CreateRenderer() fyne.WidgetRenderer {
	c := canvas.NewRectangle(color.Transparent)
	return widget.NewSimpleRenderer(c)
}

func (d *DragArea) MouseDown(e *desktop.MouseEvent) {
	hwnd, _, _ := procGetActiveWindow.Call()
	if hwnd != 0 {
		procReleaseCapture.Call()
		procSendMessage.Call(hwnd, 0x00A1, 2, 0) // WM_NCLBUTTONDOWN, HTCAPTION
	}
}

func (d *DragArea) MouseUp(e *desktop.MouseEvent) {}
func (d *DragArea) MouseMoved(e *desktop.MouseEvent) {}

var AppTitleLabel *canvas.Text

func BuildCustomTitleBar(app *WinjaApp) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.NRGBA{R: 12, G: 16, B: 22, A: 255})
	
	// Left: WINJA block
	winjaBg := canvas.NewRectangle(color.NRGBA{R: 0, G: 200, B: 255, A: 255})
	winjaBg.CornerRadius = 12
	winjaBg.SetMinSize(fyne.NewSize(67, 37))
	
	winjaText := canvas.NewText("WINJA", color.Black)
	winjaText.TextSize = 16
	winjaText.TextStyle = fyne.TextStyle{Bold: true}
	winjaText.Alignment = fyne.TextAlignCenter
	
	winjaBlock := container.NewPadded(container.NewStack(winjaBg, container.NewCenter(winjaText)))

	btnMin := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		if core.IsMinimizeToTray() {
			app.Window.Hide()
			ShowDropTarget(app.App, app.Window, app)
		} else {
			// Normal minimize using win32
			hwnd, _, _ := procGetActiveWindow.Call()
			if hwnd != 0 {
				user32.NewProc("ShowWindow").Call(hwnd, 6) // SW_MINIMIZE
			}
		}
	})
	
	// Maximize can be tricky in frameless, placeholder for now
	btnMax := widget.NewButtonWithIcon("", theme.ViewFullScreenIcon(), func() {
		if app.Window.FullScreen() {
			app.Window.SetFullScreen(false)
		} else {
			app.Window.SetFullScreen(true)
		}
	})
	
	btnClose := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		app.Window.Close()
	})
	
	controlsHBox := container.NewHBox(btnMin, btnMax, btnClose)
	controlsBg := canvas.NewRectangle(color.NRGBA{0, 0, 0, 255})
	controlsBg.CornerRadius = 8
	controlsPill := container.NewStack(controlsBg, controlsHBox)
	controls := container.NewStack(container.NewGridWrap(fyne.NewSize(100, 30), controlsPill))

	// Middle: Drag area
	drag := NewDragArea()
	AppTitleLabel = canvas.NewText("", color.White)
	AppTitleLabel.Alignment = fyne.TextAlignCenter
	AppTitleLabel.TextSize = 14

	content := container.NewBorder(nil, nil, winjaBlock, controls, container.NewStack(drag, container.NewCenter(AppTitleLabel)))
	
	return container.NewStack(bg, content)
}














