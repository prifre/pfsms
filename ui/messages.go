package ui

import (
	"fmt"
	"log"
	"strings"

	"pfsms/general"
	"pfsms/pfdatabase"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	window         fyne.Window
	inQueue        *widget.Label
	db             *pfdatabase.DBtype
}

func NewMessages(w fyne.Window) *theform {
	return &theform{
		window: w,
		db:     new(pfdatabase.DBtype),
	}
}

func (s *theform) buildMessages() *container.Scroll {
	var v string

	s.groupname = widget.NewEntry()
	s.groupname.SetText(fyne.CurrentApp().Preferences().StringWithFallback("groupname", ""))
	s.groupname.OnChanged = func(val string) {
		fyne.CurrentApp().Preferences().SetString("groupname", val)
	}

	s.phone = widget.NewMultiLineEntry()
	s.phone.Wrapping = fyne.TextWrapWord
	s.phone.SetMinRowsVisible(6)
	s.phone.SetText(fyne.CurrentApp().Preferences().String("messagephone"))
	s.phone.OnChanged = func(val string) {
		fyne.CurrentApp().Preferences().SetString("messagephone", val)
	}

	s.dataAllGroups = s.db.ShowAllGroups()
	if s.groupname.Text != "" {
		v = s.groupname.Text
	} else if len(s.dataAllGroups) > 0 {
		v = s.dataAllGroups[0][0]
	}

	s.phone.SetText(s.Getphonesforgroup(v))

	s.groupSelect = &widget.Select{
		Options:  s.db.ShowGroups(),
		Selected: v,
		OnChanged: func(val string) {
			s.phone.SetText(s.Getphonesforgroup(val))
			fyne.CurrentApp().Preferences().SetString("messagephone", s.phone.Text)
			s.groupname.SetText(val)
			fyne.CurrentApp().Preferences().SetString("groupname", val)
		},
	}

	s.btnSaveGroup = widget.NewButton("Save Group", func() {
		if len(strings.TrimSpace(s.phone.Text)) < 5 {
			return
		}

		z := s.phone.Text
		z = strings.ReplaceAll(z, "\r", ", ")
		z = strings.ReplaceAll(z, "\n", ", ")

		// Filtrera godkända tecken
		var sb strings.Builder
		for _, r := range z {
			if strings.ContainsRune("0123456789+, ", r) {
				sb.WriteRune(r)
			}
		}

		// Säker rensning av komman och mellanslag i start/slut utan risk för index-krasch
		cleaned := strings.Trim(sb.String(), " ,")
		s.phone.SetText(cleaned)

		if strings.Contains(s.phone.Text, ",") && len(strings.TrimSpace(s.groupname.Text)) > 1 {
			s.db.SaveGroup(s.groupname.Text, s.phone.Text)
		}

		s.RefreshGroupData()
		s.groupSelect.SetSelected(s.groupname.Text)
		s.phone.SetText(s.Getphonesforgroup(s.groupSelect.Selected))
	})

	s.btnDeleteGroup = widget.NewButton("Delete Group", func() {
		selected := s.groupSelect.Selected
		if selected == "" {
			return
		}

		s.db.DeleteGroup(selected)
		s.RefreshGroupData()
		s.groupSelect.ClearSelected()
		s.groupname.SetText("")
		s.phone.SetText("")
	})

	s.message = widget.NewMultiLineEntry()
	s.message.Wrapping = fyne.TextWrapWord
	s.message.SetMinRowsVisible(8)

	msgText := fyne.CurrentApp().Preferences().String("message")
	msgText = strings.ReplaceAll(msgText, `\r`, "\r")
	msgText = strings.ReplaceAll(msgText, `\n`, "\n")
	s.message.SetText(msgText)

	s.message.OnChanged = func(val string) {
		fyne.CurrentApp().Preferences().SetString("message", val)
	}

	s.btnSubmit = widget.NewButton("Add message to sending queue", func() {
		err := s.AddToSMSQueue(s.phone.Text, s.groupname.Text, s.message.Text)
		if err != nil {
			log.Println("AddToSMSQueue failed: ", err.Error())
		}
		s.updateQueueCount()
	})

	s.inQueue = general.NewBoldLabel(fmt.Sprintf("%d", s.db.CountinQueue()))

	return container.NewScroll(
		container.NewVBox(
			s.SetForm(),
		),
	)
}

func (s *theform) SetForm() *widget.Form {
	GroupsInfo := "To use multiple mobile numbers, separate them with commas or Enter.\n" +
		"Click Save Group to reuse in future."
	MessageInfo := "To insert firstname and/or lastname, use <<Fname>> and <<Lname>> in message."

	s.form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: general.NewBoldLabel(GroupsInfo)},
			{Text: "Groups", Widget: container.NewGridWithColumns(3, s.groupSelect, s.btnSaveGroup, s.btnDeleteGroup)},
			{Text: "Groupname", Widget: s.groupname},
			{Text: "Phone", Widget: s.phone},
			{Text: "", Widget: general.NewBoldLabel(MessageInfo)},
			{Text: "Message", Widget: s.message},
			{Text: "", Widget: s.btnSubmit},
			{Text: "Phones in Queue: ", Widget: s.inQueue},
		},
	}
	return s.form
}

func (s *theform) Getphonesforgroup(v string) string {
	if v == "" {
		return ""
	}
	var phones []string
	dag := s.db.ShowAllGroups()
	for _, row := range dag {
		if len(row) >= 2 && row[0] == v {
			phones = append(phones, row[1])
		}
	}
	return strings.Join(phones, ", ")
}

func (s *theform) AddToSMSQueue(phone, groupname, msg string) error {
	ph := phone
	ph = strings.ReplaceAll(ph, "\r", ",")
	ph = strings.ReplaceAll(ph, "\n", ",")

	var sb strings.Builder
	for _, r := range ph {
		if strings.ContainsRune("0123456789+,", r) {
			sb.WriteRune(r)
		}
	}

	rawPhones := strings.Split(sb.String(), ",")
	var p2List []string

	for _, p := range rawPhones {
		trimmed := strings.TrimSpace(p)
		if len(trimmed) < 5 {
			continue
		}
		fixedNum := general.Fixphonenumber(trimmed)
		fname := s.db.GetFname(trimmed)
		lname := s.db.GetLname(trimmed)

		// Formatera post för köhantering
		item := fmt.Sprintf("%s\t%s\t%s", fixedNum, fname, lname)
		p2List = append(p2List, item)
	}

	if len(p2List) == 0 {
		return nil
	}

	return s.db.AddToQueue(p2List, groupname, msg)
}

func (s *theform) RefreshGroupData() {
	s.dataAllGroups = s.db.ShowAllGroups()
	s.groupSelect.Options = s.db.ShowGroups()
	s.groupSelect.Refresh()
}

func (s *theform) updateQueueCount() {
	if s.inQueue != nil {
		s.inQueue.SetText(fmt.Sprintf("%d", s.db.CountinQueue()))
	}
}

func (s *theform) TabItem() *container.TabItem {
	return &container.TabItem{
		Text:    "Messages",
		Icon:    theme.MailSendIcon(),
		Content: s.buildMessages(),
	}
}

func (s *theform) RefreshData() {
	s.RefreshGroupData()
	s.updateQueueCount()
}
