package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppMessages struct {
	// Theme holds the current theme
	Theme string
}

type theform struct {
	phone	 		*widget.Entry
	message	 		*widget.Entry
	appMessages 	*AppMessages
	window      	fyne.Window
	app         	fyne.App
}

func NewMessages(a fyne.App, w fyne.Window,  am *AppMessages) *theform {
	return &theform{app: a, window: w,  appMessages: am}
}
func (s *theform) buildForm() *container.Scroll {
	s.phone = &widget.Entry{}
	s.phone.Text="00"
	s.message = &widget.Entry{}
	s.message.Text="00"

	return container.NewScroll(container.NewVBox(
		s.phone,s.message,
	))
}
func (s *theform) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Messages", Icon: theme.MailSendIcon(), Content: s.buildForm()}
}
