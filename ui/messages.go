package ui

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"
	"github.com/prifre/pfsms/pfmobile"
)

type theform struct {
	form           *widget.Form
	phone          *widget.Entry
	groupname      *widget.Entry
	message        *widget.Entry
	btnSaveGroup   *widget.Button
	btnDeleteGroup *widget.Button
	btnSubmit      *widget.Button
	groupSelect    *widget.Select
	dataAllGroups  [][]string
	logtext        *widget.Label
	window         fyne.Window
}

func NewMessages(w fyne.Window) *theform {
	return &theform{window: w}
}
func (s *theform) buildUI() *container.Scroll {
	s.groupname = &widget.Entry{}
	s.groupname.Text = fyne.CurrentApp().Preferences().StringWithFallback("groupname","")
	s.groupname.OnChanged = func(v string) {
		fyne.CurrentApp().Preferences().SetString("groupname",s.groupname.Text)
	}
	s.phone = &widget.Entry{}
	s.phone = widget.NewMultiLineEntry()
	s.phone.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
	s.phone.SetMinRowsVisible(6)
	s.phone.Text = fyne.CurrentApp().Preferences().String("messagephone")
	s.phone.OnChanged = func(v string) {
		 fyne.CurrentApp().Preferences().SetString("messagephone",s.phone.Text)
	}
	s.dataAllGroups = new(pfdatabase.DBtype).ShowAllGroups()
	s.groupSelect = &widget.Select{Options: new(pfdatabase.DBtype).ShowGroups(),
		Selected: "",
		OnChanged: func(v string) {
			s.phone.Text = s.Getphonesforgroup(v)
			fyne.CurrentApp().Preferences().SetString("messagephone",s.phone.Text)
			s.phone.Refresh()
			s.groupname.Text = v
			fyne.CurrentApp().Preferences().SetString("groupname",s.groupname.Text)
			s.groupname.Refresh()
		},
	}
	s.btnSaveGroup = &widget.Button{Text: "Save Group", OnTapped: func() {
		// whe4n pasting, ensure so onlyu "0123456789+, " survives...
		if len(s.phone.Text) < 5 {
			return
		}
		z := s.phone.Text
		z = strings.Replace(z, "\r", ", ", -1)
		z2 := ""
		for i := 0; i < len(z); i++ {
			if strings.Contains("0123456789+, ", string(z[i])) {
				z2 += string(z[i])
			}
		}
		z = z2
		for strings.Contains(" ,", string(z[len(z)-1])) {
			z = z[:len(z)-2]
		}
		for strings.Contains(" ,", string(z[0])) {
			z = z[1:]
		}
		s.phone.Text = z
		if strings.Contains(s.phone.Text, ",") && len(s.groupname.Text) > 1 {
			new(pfdatabase.DBtype).SaveGroup(s.groupname.Text, s.phone.Text)
		}
		s.dataAllGroups = new(pfdatabase.DBtype).ShowAllGroups()
		s.groupSelect.Options = new(pfdatabase.DBtype).ShowGroups()
		s.groupSelect.SetSelected(s.groupname.Text)
		s.groupSelect.Refresh()
		s.phone.Text = s.Getphonesforgroup(s.groupSelect.Selected)
		s.phone.Refresh()
	}}
	s.btnDeleteGroup = &widget.Button{Text: "Delete Group", OnTapped: func() {
		new(pfdatabase.DBtype).DeleteGroup(s.groupSelect.Selected)
		s.dataAllGroups = new(pfdatabase.DBtype).ShowAllGroups()
		s.groupSelect.Options = new(pfdatabase.DBtype).ShowGroups()
		s.groupSelect.Selected = ""
		s.groupSelect.Refresh()
		s.groupname.SetText("")
		s.groupname.Refresh()
	}}
	GroupsInfo := "To use multiple mobile numbers, separate them with commas or Enter.\r\n"
	GroupsInfo += "Click Save Group to reuse in future."
	MessageInfo := "To insert firstname and/or lastname, use <<Fname>> and <<Lname>> in message."
	s.message = &widget.Entry{}
	s.message = widget.NewMultiLineEntry()
	s.message.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
	s.message.SetMinRowsVisible(8)
	s.message.Text = fyne.CurrentApp().Preferences().String("message")
	s.message.OnChanged = func (v string) {
		fyne.CurrentApp().Preferences().SetString("message",s.message.Text)
	}
	s.btnSubmit = &widget.Button{Text: "Click to send message", OnTapped: func() {
		s.HandleSendsms(s.phone.Text, s.groupname.Text, s.message.Text)
	}}
	s.form = &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "", Widget: NewBoldLabel(GroupsInfo)},
			{Text: "Groups", Widget: container.NewGridWithColumns(3, s.groupSelect, s.btnSaveGroup, s.btnDeleteGroup)},
			{Text: "Groupname", Widget: s.groupname},
			{Text: "Phone", Widget: s.phone},
			{Text: "", Widget: NewBoldLabel(MessageInfo)},
			{Text: "Message", Widget: s.message},
			{Text: "", Widget: s.btnSubmit},
		},
	}
	m := ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 12)
	s.logtext = &widget.Label{Text: m}
	s.form.Refresh()
	return container.NewScroll(
		container.NewVBox(
			s.form,
			s.logtext,
		))
}
func (s *theform) Getphonesforgroup(v string) string {
	var np string
	dag := new(pfdatabase.DBtype).ShowAllGroups()
	for i := 0; i < len(dag); i++ {
		if dag[i][0] == v {
			if np > "" {
				np = np + ", " + dag[i][1]
			} else {
				np = dag[i][1]
			}
		}
	}
	return np
}
func (s *theform) HandleSendsms(p, groupname, msg string) {
	// split phone into \r\n and ","
	ph := s.phone.Text
	ph = strings.Replace(ph, "\r", ",", -1)
	ph = strings.Replace(ph, "\n", ",", -1)
	ph = strings.Replace(ph, ",,", ",", -1)
	ph = strings.Replace(ph, ",,", ",", -1)
	ph2 := ""
	for i := 0; i < len(ph); i++ {
		if strings.Contains("0123456789+,", string(ph[i])) {
			ph2 = ph2 + string(ph[i])
		}
	}
	p1 := strings.Split(ph2, ",")
	countrycode := fyne.CurrentApp().Preferences().StringWithFallback("mobilecountry", "Sweden(+46)")
	for i := 0; i < len(p1); i++ {
		p1[i] = Fixphonenumber(p1[i], countrycode)
		p1[i] += "\t" + new(pfdatabase.DBtype).GetFname(p1[i]) + "\t" + new(pfdatabase.DBtype).GetLname(p1[i])
	}
	result := new(pfmobile.SMStype).SendMessage(p1, msg)
	if result != nil {
		log.Println("Sent messages ok")
		var sh [][]string
		for i := 0; i < len(result); i++ {
			// result = tstamp, phone, message
			sh = append(sh, []string{result[i][0], groupname, result[i][1], result[i][2]})
		}
		new(pfdatabase.DBtype).SaveHistory(sh)
		s.logtext.Text = ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 8)
		s.logtext.Refresh()
	}
}
func (s *theform) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Messages", Icon: theme.MailSendIcon(), Content: s.buildUI()}
}
