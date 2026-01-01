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
		NewCustomers(window).tabItem(),
		NewMessages(window).tabItem(),
		NewQueue(window).tabItem(),
		NewEmail(window).tabItem(),
		NewSettings(window).tabItem(),
		NewAbout(window).tabItem(),
	}
	at := container.AppTabs{Items: tabs}
	at.OnSelected = func(t *container.TabItem) {
		switch t.Text {
		case "Customers":
			tabs[0] = NewCustomers(window).tabItem()
		case "Messages":
			tabs[1] = NewMessages(window).tabItem()
		case "Queue":
			tabs[2] = NewQueue(window).tabItem()
		case "Email":
			tabs[3] = NewEmail(window).tabItem()
		case "Settings":
			tabs[4] = NewSettings(window).tabItem()
		case "About pfsms":
			tabs[5] = NewAbout(window).tabItem()
		default:
		}
	}
	return &at
}

