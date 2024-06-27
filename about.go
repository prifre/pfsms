package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type about struct {
}

func newAbout() *about {
	return &about{}
}

func (a *about) buildUI() *fyne.Container {
	return container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(		
			layout.NewSpacer(),
			container.NewVBox(		
				newBoldLabel("PFSMS"), 
				layout.NewSpacer(),
				newBoldLabel(version),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}
func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About", Content: a.buildUI()}
}
