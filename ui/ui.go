package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	// appTable := &AppTable{}
	appMessages := &AppMessages{}
	appEmail := &AppEmail{}
	appHistory := &AppHistory{}
	appSettings := &AppSettings{}
	appSettings.Theme = checkTheme(app.Preferences().StringWithFallback("Theme", "Adaptive (requires restart)"), app)
var tabs []*container.TabItem =  []*container.TabItem{
		NewTable(app,window,&AppTable{}).tabItem(),
		NewMessages(app,window,appMessages).tabItem(),
		NewEmaillog(app,window,appEmail).tabItem(),
		NewHistory(app,window,appHistory).tabItem(),
		NewSettings(app, window,  appSettings).tabItem(),
		NewAbout().tabItem(),
	}
	return &container.AppTabs{Items:tabs}
}