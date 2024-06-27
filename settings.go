package main

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
	themeSelect *widget.Select

	overwriteFiles     *widget.RadioGroup
	notificationRadio  *widget.RadioGroup

	componentSlider     *widget.Slider
	componentLabel      *widget.Label

	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newSettings(a fyne.App, w fyne.Window,  as *AppSettings) *settings {
	return &settings{app: a, window: w,  appSettings: as}
}

func (s *settings) onThemeChanged(selected string) {
	s.app.Preferences().SetString("Theme", checkTheme(selected, s.app))
}

func (s *settings) onOverwriteFilesChanged(selected string) {
//	s.client.OverwriteExisting = selected == "On"
	s.app.Preferences().SetString("OverwriteFiles", selected)
}

func (s *settings) onNotificationsChanged(selected string) {
	s.app.Preferences().SetString("Notifications", selected)
}

func (s *settings) onComponentsChange(value float64) {
	s.componentLabel.SetText(string('0' + byte(value)))
}



func (s *settings) buildUI() *container.Scroll {
	s.themeSelect = &widget.Select{Options: themes, OnChanged: s.onThemeChanged, Selected: s.appSettings.Theme}

	s.overwriteFiles = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onOverwriteFilesChanged}
	s.overwriteFiles.SetSelected(s.app.Preferences().StringWithFallback("OverwriteFiles", "Off"))

	s.notificationRadio = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onNotificationsChanged}
	s.notificationRadio.SetSelected(s.app.Preferences().StringWithFallback("Notifications", onOffOptions[1]))

	s.componentSlider, s.componentLabel = &widget.Slider{Min: 2.0, Max: 6.0, Step: 1, OnChanged: s.onComponentsChange}, &widget.Label{}
	s.componentSlider.SetValue(s.app.Preferences().FloatWithFallback("ComponentLength", 2))

	interfaceContainer := container.NewGridWithColumns(2,
		newBoldLabel("Application Theme"), s.themeSelect,
	)

	dataContainer := container.NewGridWithColumns(2,
		newBoldLabel("Overwrite Files"), s.overwriteFiles,
		newBoldLabel("Notifications"), s.notificationRadio,
	)

	
	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "User Interface", Content: interfaceContainer},
		&widget.Card{Title: "Data Handling", Content: dataContainer},
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
