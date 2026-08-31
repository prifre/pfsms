package ui

import (
	"pfsms/general"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// Create syr ihop alla UI-komponenter och returnerar huvudflikarna
func Create(window fyne.Window) *container.AppTabs {
	fyne.CurrentApp().Settings().SetTheme(theme.DefaultTheme())
	general.Setupfiles()

	// Skapa vyerna en gång
	customersView := NewCustomers(window)
	messagesView := NewMessages(window)
	queueView := NewQueue(window)
	//	emailView := NewEmail(window)
	settingsView := NewSettings(window)
	aboutView := NewAbout(window)

	// Skapa flik-objektet med Fynes inbyggda konstruktor
	at := container.NewAppTabs(
		customersView.TabItem(),
		messagesView.TabItem(),
		queueView.TabItem(),
		//		emailView.TabItem(),
		settingsView.TabItem(),
		aboutView.TabItem(),
	)

	// Om du vill göra något när en flik väljs (t.ex. uppdatera data)
	at.OnSelected = func(t *container.TabItem) {
		switch t {
		case at.Items[0]:
			customersView.RefreshData() // Exempel: Ladda om kundlistan från databasen
		case at.Items[1]:
			messagesView.RefreshData()
		case at.Items[2]:
			queueView.RefreshTables()
		}
	}
	return at
}
