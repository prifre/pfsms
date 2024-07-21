package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppMessages struct {
	// Theme holds the current theme
	Theme string
}

type theform struct {
	form			*widget.Form
	phone	 		*widget.Entry
	reference 		*widget.Entry
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
	s.reference = &widget.Entry{}
	s.reference.Text="test1"
	s.message = &widget.Entry{}
	s.message=widget.NewMultiLineEntry()
	s.message.Text=""
	s.form = &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Phone", Widget: s.phone},
			{Text: "Reference", Widget: s.reference},
			{Text: "Messagetext", Widget: s.message},
		},
		OnSubmit: func() { // optional, handle form submission
			s.HandleSubmit( s.phone.Text, s.reference.Text, s.message.Text)
		},
		// ariasms.SendMessage([]string{s.phone.Text},s.message.Text)
		}
	return container.NewScroll(		
		container.NewVBox(
			s.form,
	))
}

func (s *theform) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Messages", Icon: theme.MailSendIcon(), Content: s.buildForm()}
}

func (s *theform) HandleSubmit(p,t,m string) {
	txt :=""
	txt += fmt.Sprintln("Phone:", s.phone.Text)
	txt += fmt.Sprintln("Reference:", s.reference.Text)
	txt += fmt.Sprintln("Message:", s.message.Text)
	dialog.ShowInformation("Version",txt, s.window)
}
