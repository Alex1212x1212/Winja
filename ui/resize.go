package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"syscall"
)

var user32_resize = syscall.NewLazyDLL("user32.dll")
var procReleaseCaptureResize = user32_resize.NewProc("ReleaseCapture")
var procSendMessageResize = user32_resize.NewProc("SendMessageW")
var procGetActiveWindowResize = user32_resize.NewProc("GetActiveWindow")

const WM_SYSCOMMAND = 0x0112
const SC_SIZE = 0xF000

type ResizeBorder struct {
	widget.BaseWidget
	app     *WinjaApp
	edge    string // "top", "bottom", "left", "right", "topleft", "topright", "bottomleft", "bottomright"
}

func NewResizeBorder(app *WinjaApp, edge string) *ResizeBorder {
	r := &ResizeBorder{app: app, edge: edge}
	r.ExtendBaseWidget(r)
	return r
}

func (r *ResizeBorder) MinSize() fyne.Size {
	if r.edge == "left" || r.edge == "right" {
		return fyne.NewSize(8, 10)
	}
	return fyne.NewSize(10, 8)
}

func (r *ResizeBorder) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.Transparent)
	return widget.NewSimpleRenderer(bg)
}

func (r *ResizeBorder) Cursor() desktop.Cursor {
	if r.edge == "left" || r.edge == "right" {
		return desktop.HResizeCursor
	}
	return desktop.VResizeCursor
}

func (r *ResizeBorder) Tapped(e *fyne.PointEvent) {} // Just to consume

func (r *ResizeBorder) MouseDown(e *desktop.MouseEvent) {
	if e.Button == desktop.MouseButtonPrimary {
		hwnd, _, _ := procGetActiveWindowResize.Call()
		if hwnd != 0 {
			var direction uintptr = 0
			switch r.edge {
			case "left":
				direction = 1
			case "right":
				direction = 2
			case "top":
				direction = 3
			case "topleft":
				direction = 4
			case "topright":
				direction = 5
			case "bottom":
				direction = 6
			case "bottomleft":
				direction = 7
			case "bottomright":
				direction = 8
			}
			
			if direction != 0 {
				procReleaseCaptureResize.Call()
				procSendMessageResize.Call(hwnd, uintptr(WM_SYSCOMMAND), uintptr(SC_SIZE+direction), 0)
			}
		}
	}
}

func (r *ResizeBorder) MouseUp(e *desktop.MouseEvent) {}
