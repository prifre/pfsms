package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppHistory struct {
	// Theme holds the current theme
	Theme string
}

type thehistory struct {
	history				*widget.Label
	window      		fyne.Window
	appHistory 			*AppHistory
	app         		fyne.App
}

func NewHistory(a fyne.App, w fyne.Window,  el *AppHistory) *thehistory {
	return &thehistory{app: a, window: w,  appHistory: el}
}
func (s *thehistory) buildHistory() *container.Scroll {
	s.history = &widget.Label{}
	return container.NewScroll(		
		container.NewVBox(
			s.history,
	))
}

func (s *thehistory) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Mail Log", Icon: theme.DocumentIcon(), Content: s.buildHistory()}
}
func (s *thehistory) Addhistory (m string) {
	s.history.SetText(fmt.Sprintf("%s%s",s.history.Text ,m ))
}

