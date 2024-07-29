package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type thesmslog struct {
	history				*widget.Label
	window      		fyne.Window
	app         		fyne.App
}
// the SMS Tab 

func NewSmslog(a fyne.App, w fyne.Window) *thesmslog {
	return &thesmslog{app: a, window: w}
}
func (s *thesmslog) buildHistory() *container.Scroll {
	s.history = &widget.Label{}
	return container.NewScroll(		
		container.NewVBox(
			s.history,
	))
}

func (s *thesmslog) tabItem() *container.TabItem {
	return &container.TabItem{Text: "SMS", Icon: theme.DocumentIcon(), Content: s.buildHistory()}
}
func (s *thesmslog) Addhistory (m string) {
	s.history.SetText(fmt.Sprintf("%s%s",s.history.Text ,m ))
}

