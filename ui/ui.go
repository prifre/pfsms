package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// Create will stitch together all ui components
func Create(window fyne.Window) *container.AppTabs {
	fyne.CurrentApp().Settings().SetTheme(theme.DefaultTheme())
	Setupfiles()
	var tabs []*container.TabItem = []*container.TabItem{
		NewTable(window).tabItem(),
		NewMessages(window).tabItem(),
		// NewEmail(window).tabItem(),
		NewSettings(window).tabItem(),
		NewAbout(window).tabItem(),
	}
	at := container.AppTabs{Items: tabs}
	at.OnSelected = func(t *container.TabItem) {
		switch t.Text {
		case "Customers":
			tabs[0] = NewTable(window).tabItem()
		case "Messages":
			tabs[1] = NewMessages(window).tabItem()
		// case "Email":
		// 	tabs[2]=NewEmail(window).tabItem()
		case "Settings":
			tabs[2] = NewSettings(window).tabItem()
		case "About pfsms":
			tabs[3] = NewAbout(window).tabItem()
		default:
		}
	}
	return &at
}
