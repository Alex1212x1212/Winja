package main

import (
	"testing"
	"fyne.io/fyne/v2/test"
	"winja/ui"
)

func TestAppDoesNotCrash(t *testing.T) {
	app := test.NewApp()
	w := app.NewWindow("Test")
	
	mainApp := ui.NewApp(app, w)
	w.SetContent(mainApp.BuildContent())
	
	// Test switching to settings
	mainApp.SwitchToSettings()
	// Test switching to tools
	mainApp.SwitchToTools()
	// Test switching to main
	mainApp.SwitchToMain()
}
