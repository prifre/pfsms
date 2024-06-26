package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	return &container.AppTabs{Items: []*container.TabItem{
//		newSettings(app, window, bridge, appSettings).tabItem(),
		newAbout().tabItem(),
	}}}

func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About",  Content: a.buildUI()}
}


type about struct {
	icon        *canvas.Image
	nameLabel   *widget.Label
}

func newAbout() *about {
	return &about{}
}

func (a *about) buildUI() *fyne.Container {
	return container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), a.icon, layout.NewSpacer()),
		container.NewHBox(
			layout.NewSpacer(),
			a.nameLabel,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}
