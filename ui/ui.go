package ui

import (
	"log"
	"time"

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
		NewSettings(app, window).tabItem(),
		NewAbout(app,window).tabItem(),
	}
	Setuplog()
	log.Printf("%s started!",time.Now().Format("2006-01-02 15:04:05"))
	return &container.AppTabs{Items:tabs}
}
