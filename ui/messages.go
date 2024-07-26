package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/ariasms"
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
	s.reference = &widget.Entry{}
	s.reference.Text="test1"
	s.phone = &widget.Entry{}
	s.phone=widget.NewMultiLineEntry()
	s.phone.Wrapping=fyne.TextWrap(fyne.TextWrapWord)
	s.phone.SetMinRowsVisible(8)
	s.phone.Text=""
	s.message = &widget.Entry{}
	s.message=widget.NewMultiLineEntry()
	s.message.Wrapping=fyne.TextWrap(fyne.TextWrapWord)
	s.message.SetMinRowsVisible(12)
	s.message.Text=""
	s.form = &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Reference", Widget: s.reference},
			{Text: "Phone", Widget: s.phone},
			{Text: "Messagetext", Widget: s.message},
		},
		OnSubmit: func() { // optional, handle form submission
			s.HandleSubmit( s.phone.Text, s.reference.Text, s.message.Text)
		},}
		// sms:=new(ariasms.SMStype)
		// sms.SendMessage([]string{s.phone.Text},s.message.Text)		
	return container.NewScroll(		
		container.NewVBox(
			widget.NewLabel("To use multiple mobile numbers, separate them with commas or Enter."),
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
	// split phone into \r\n and ","
	ph := s.phone.Text
	ph = strings.Replace(ph,"\r",",",-1)
	ph = strings.Replace(ph,"\n",",",-1)
	ph = strings.Replace(ph,",,",",",-1)
	ph = strings.Replace(ph,",,",",",-1)
	p1:=strings.Split(ph,",")
	var sms ariasms.SMStype =*new(ariasms.SMStype)
	sms.SendMessage(p1,m)
	for i:=0;i<len(p1);i++ {
		Appendtotextfile("smshistory.txt","Sending "+s.message.Text+" to "+p1[i])
	}
}
