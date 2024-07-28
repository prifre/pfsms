package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	app.Settings().SetTheme(theme.DefaultTheme())
var tabs []*container.TabItem =  []*container.TabItem{
		NewTable(app,window,&AppTable{}).tabItem(),
		NewMessages(app,window).tabItem(),
		NewEmaillog(app,window).tabItem(),
		NewSmslog(app,window).tabItem(),
		NewSettings(app, window).tabItem(),
		NewAbout(app,window).tabItem(),
	}
	return &container.AppTabs{Items:tabs}
}