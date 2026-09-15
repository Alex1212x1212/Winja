package main

import (
	"fmt"
	"syscall"
	"time"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procGetActiveWindow = user32.NewProc("GetActiveWindow")
var procGetWindowLong = user32.NewProc("GetWindowLongW")
var procSetWindowLong = user32.NewProc("SetWindowLongW")
var procSetWindowPos = user32.NewProc("SetWindowPos")

var GWL_STYLE int = -16
var WS_THICKFRAME uintptr = 0x00040000
var WS_CAPTION uintptr = 0x00C00000

const SWP_FRAMECHANGED = 0x0020
const SWP_NOMOVE = 0x0002
const SWP_NOSIZE = 0x0001
const SWP_NOZORDER = 0x0004

func MakeResizable() {
	go func() {
		// Try a few times to get the active window
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			hwnd, _, _ := procGetActiveWindow.Call()
			if hwnd != 0 {
				style, _, _ := procGetWindowLong.Call(hwnd, uintptr(GWL_STYLE))
				
				// Only apply if it doesn't already have WS_THICKFRAME
				if style & WS_THICKFRAME == 0 {
					style = style | WS_THICKFRAME
					style = style &^ WS_CAPTION
					procSetWindowLong.Call(hwnd, uintptr(GWL_STYLE), style)
					procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED)
					fmt.Println("Applied WS_THICKFRAME to HWND", hwnd)
					return
				}
			}
		}
	}()
}
