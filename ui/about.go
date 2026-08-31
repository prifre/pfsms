package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"pfsms/general"
	"pfsms/pfdatabase"

	"fyne.io/fyne/v2"
	appearance "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type theabout struct {
	btnAppearance *widget.Button
	dodebug       *widget.RadioGroup
	dbLabel       *widget.Label
	memLabel      *widget.Label
	window        fyne.Window
	db            *pfdatabase.DBtype
	updatingUI    bool // Flagga för att förhindra oändliga loopar i event handlern
}

func NewAbout(w fyne.Window) *theabout {
	return &theabout{
		window: w,
		db:     new(pfdatabase.DBtype),
	}
}

func (a *theabout) buildAbout() *fyne.Container {
	a.btnAppearance = widget.NewButton("Change appearance!", func() {
		dialog.NewCustom(
			"Fix the looks for the application!",
			"Close",
			appearance.NewSettings().LoadAppearanceScreen(a.window),
			a.window,
		).Show()
	})

	a.dodebug = widget.NewRadioGroup([]string{"Yes", "No"}, nil)
	a.dodebug.Horizontal = true
	a.dodebug.Required = true

	// Sätt initialt värde
	a.updatingUI = true
	if fyne.CurrentApp().Preferences().Bool("debug") {
		a.dodebug.Selected = "Yes"
	} else {
		a.dodebug.Selected = "No"
	}
	a.updatingUI = false

	a.dodebug.OnChanged = func(v string) {
		// Förhindra körning om ändringen gjordes från kod
		if a.updatingUI {
			return
		}

		currentSetting := fyne.CurrentApp().Preferences().Bool("debug")
		wantDebug := (v == "Yes")

		if currentSetting == wantDebug {
			return
		}

		// Räkna ut sökvägarna för informationstexten
		currentDir := general.GetHomeDir()

		var newDir string
		if wantDebug {
			exePath, err := os.Executable()
			if err == nil {
				newDir = filepath.Join(filepath.Dir(exePath), "pfsmsdata")
			} else {
				newDir, _ = os.Getwd()
				newDir = filepath.Join(newDir, "pfsmsdata")
			}
		} else {
			home, _ := os.UserHomeDir()
			newDir = filepath.Join(home, "pfsmsdata")
		}

		m := "Application needs to restart after Debug setting has been changed.\n\n" +
			"Current data location:\n" + currentDir + "\n\n" +
			"New data location after restart:\n" + newDir

		dialog.NewConfirm("Really change Debug setting?", m, func(confirmed bool) {
			if !confirmed {
				// Återställ värdet i UI utan att utlösa dialogen igen
				a.updatingUI = true
				if currentSetting {
					a.dodebug.SetSelected("Yes")
				} else {
					a.dodebug.SetSelected("No")
				}
				a.updatingUI = false
				return
			}

			// Om bekräftat: spara och avsluta
			fyne.CurrentApp().Preferences().SetBool("debug", wantDebug)
			general.Setupfiles()
			fyne.CurrentApp().Quit()
		}, a.window).Show()
	}

	a.dbLabel = general.NewBoldLabel(a.getDatabaseText())
	a.memLabel = general.NewBoldLabel(a.getAboutText())

	return container.NewVBox(
		&widget.Card{
			Title: "App Info",
			Content: container.NewHBox(
				layout.NewSpacer(),
				container.NewVBox(a.memLabel),
				layout.NewSpacer(),
			),
		},
		&widget.Card{
			Title: "Database",
			Content: container.NewHBox(
				layout.NewSpacer(),
				container.NewVBox(a.dbLabel),
				layout.NewSpacer(),
			),
		},
		&widget.Card{
			Title: "Interface",
			Content: container.NewHBox(
				layout.NewSpacer(),
				container.NewVBox(
					container.NewGridWithColumns(2, general.NewBoldLabel("Debug?"), a.dodebug),
					layout.NewSpacer(),
					a.btnAppearance,
				),
				layout.NewSpacer(),
			),
		},
		layout.NewSpacer(),
	)
}

func (a *theabout) getAboutText() string {
	var dtg string
	exePath, err := os.Executable()
	if err == nil {
		fi, err := os.Stat(exePath)
		if err == nil {
			dtg = "on " + fi.ModTime().Format("2006-01-02 15:04:05")
		}
	}

	appID := fyne.CurrentApp().UniqueID()
	version := fyne.CurrentApp().Metadata().Version

	return fmt.Sprintf(
		"%s\n%s\nby Peter Freund\nprifre@prifre.com\n\nCompiled with go %s %s\n\n%s",
		appID, version, runtime.Version(), dtg, general.Getmemoryinfo(),
	)
}

func (a *theabout) getDatabaseText() string {
	return fmt.Sprintf(
		"Number of customers: %d\nNumber of groups: %d\nNumber of records in groups: %d\nHistory records (# of sent sms): %d\n",
		len(a.db.ShowCustomers()),
		len(a.db.ShowGroups()),
		len(a.db.ShowAllGroups()),
		len(a.db.ShowHistory()),
	)
}

func (a *theabout) RefreshData() {
	if a.dbLabel != nil {
		a.dbLabel.SetText(a.getDatabaseText())
	}
	if a.memLabel != nil {
		a.memLabel.SetText(a.getAboutText())
	}
}

func (a *theabout) TabItem() *container.TabItem {
	return &container.TabItem{
		Text:    "About pfsms",
		Icon:    theme.SettingsIcon(),
		Content: a.buildAbout(),
	}
}
