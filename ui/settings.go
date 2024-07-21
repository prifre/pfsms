package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	themes       = []string{"Adaptive (requires restart)", "Light", "Dark"}
	onOffOptions = []string{"On", "Off"}
)

// AppSettings contains settings specific to the application
type AppSettings struct {
	// Theme holds the current theme
	Theme string
}

type settings struct {
	themeSelect 		*widget.Select

	useEmail	     	*widget.RadioGroup
	emailServer  		*widget.Entry
	emailSLabel      	*widget.Label
	emailUser  			*widget.Entry
	emailULabel      	*widget.Label
	emailPassword  		*widget.Entry
	emailPLabel      	*widget.Label
	emailFLabel      	*widget.Label
	emailFrequency      *widget.Slider

	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func NewSettings(a fyne.App, w fyne.Window,  as *AppSettings) *settings {
	return &settings{app: a, window: w,  appSettings: as}
}

func (s *settings) onThemeChanged(selected string) {
	s.app.Preferences().SetString("Theme", checkTheme(selected, s.app))
}

func (s *settings) onUseEmailChanged(selected string) {
//	s.client.OverwriteExisting = selected == "On"
	s.app.Preferences().SetString("UseEmail", selected)
	s.emailServer.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailUser.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailPassword.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailSLabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailULabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailPLabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
}

func (s *settings) buildUI() *container.Scroll {
	s.themeSelect = &widget.Select{Options: themes, OnChanged: s.onThemeChanged, Selected: s.appSettings.Theme}

	s.emailSLabel = &widget.Label{Text: "Email Server", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailServer = &widget.Entry{Text:""}
	s.emailULabel = &widget.Label{Text: "Email User", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailUser = &widget.Entry{Text:""}
	s.emailPLabel = &widget.Label{Text: "Email Password", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailPassword = &widget.Entry{Text:""}
	s.emailFLabel = &widget.Label{Text: "Email frequency (min)", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailFrequency = &widget.Slider{}
	s.emailFrequency=widget.NewSlider(0,60)

	s.useEmail = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onUseEmailChanged}
	s.useEmail.SetSelected(s.app.Preferences().StringWithFallback("UseEmail", "Off"))
	s.onUseEmailChanged(s.app.Preferences().StringWithFallback("UseEmail", "Off"))

	interfaceContainer := container.NewGridWithColumns(2,
		newBoldLabel("Application Theme"), s.themeSelect,
	)

	dataContainer := container.NewGridWithColumns(2,
		newBoldLabel("Use Email"), s.useEmail,
		s.emailSLabel, s.emailServer,
		s.emailULabel, s.emailUser,
		s.emailPLabel, s.emailPassword,
		s.emailFLabel,s.emailFrequency,
	)

	
	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "User Interface", Content: interfaceContainer},
		&widget.Card{Title: "Email Settings", Content: dataContainer},
	))
}

func (s *settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}

func checkTheme(themec string, a fyne.App) string {
	switch themec {
	case "Light":
		//lint:ignore SA1019 Not quite ready for removal on Linux.
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		//lint:ignore SA1019 Not quite ready for removal on Linux.
		a.Settings().SetTheme(theme.DarkTheme())
	}

	return themec
}

func newBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}
