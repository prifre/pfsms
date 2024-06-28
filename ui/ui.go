package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	appSettings := &AppSettings{}
	appTable := &AppTable{}
	appSettings.Theme = checkTheme(app.Preferences().StringWithFallback("Theme", "Adaptive (requires restart)"), app)

	return &container.AppTabs{Items: []*container.TabItem{
		NewSettings(app, window,  appSettings).tabItem(),
		NewAbout().tabItem(),
		NewTable(app,window,appTable).tabItem(),
	}}
}

