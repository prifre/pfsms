package ui

import (
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	appearance "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"

	"fyne.io/fyne/v2/layout"
)

type theabout struct {
	btnAppearance		*widget.Button
	dodebug				*widget.RadioGroup
	window  		    fyne.Window
}

func NewAbout(w fyne.Window) *theabout {
	return &theabout{ window: w}
}
func Getmemoryinfo() string {
// Alloc uint64
// Alloc is bytes of allocated heap objects.
// "Allocated" heap objects include all reachable objects, as well as unreachable objects that the garbage collector has not yet freed.
// Specifically, Alloc increases as heap objects are allocated and decreases as the heap is swept and unreachable objects are freed.
// Sweeping occurs incrementally between GC cycles, so these two processes occur simultaneously, and as a result Alloc tends to change smoothly (in contrast with the sawtooth that is typical of stop-the-world garbage collectors).
// TotalAlloc uint64
// TotalAlloc is cumulative bytes allocated for heap objects.
// TotalAlloc increases as heap objects are allocated, but unlike Alloc and HeapAlloc, it does not decrease when objects are freed.
// Sys uint64
// Sys is the total bytes of memory obtained from the OS.
// Sys is the sum of the XSys fields below. Sys measures the virtual address space reserved by the Go runtime for the heap, stacks, and other internal data structures. It's likely that not all of the virtual address space is backed by physical memory at any given moment, though in general it all was at some point.
// NumGC uint32
// NumGC is the number of completed GC cycles.
	var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
		var r string
        r += fmt.Sprintf("Memory Usage = %v MB", (m.Alloc / 1024 / 1024))
        // r += fmt.Sprintf("\r\nTotalAlloc = %v kB", (m.TotalAlloc  / 1024))
        r += fmt.Sprintf("\r\nApplication Memory = %v MB", (m.Sys  / 1024 / 1024))
        // r += fmt.Sprintf("\r\nNumGC = %v\n", m.NumGC)
		return r
}
func (a *theabout) abouttext() *fyne.Container {
	var m string
	bin := fyne.CurrentApp().UniqueID()
	var dtg string
	fi, err := os.Stat(bin)
	if err == nil {
		dtg = "on " + fi.ModTime().Format("2006-01-02.15:04:05")
	}
	m +=fyne.CurrentApp().UniqueID()
	m +="\r\n"
	m +=fyne.CurrentApp().Metadata().Version
	m +="\r\n"
	m +="by Peter Freund\r\nprifre@prifre.com\r\n\r\n"
	m += fmt.Sprintf("Compiled with go %s %s\r\n", runtime.Version(), dtg)
	m += "\r\n"+Getmemoryinfo()
	return container.NewVBox(NewBoldLabel(m))
}
func (a *theabout) aboutdatabase() *fyne.Container {
	var m string
	m += fmt.Sprintf("Number of customers: %d\r\n",len(new(pfdatabase.DBtype).ShowCustomers()))
	m += fmt.Sprintf("Number of groups: %d\r\n",len(new(pfdatabase.DBtype).ShowGroups()))
	m += fmt.Sprintf("Number of records in groups: %d\r\n",len(new(pfdatabase.DBtype).ShowAllGroups()))
	m += fmt.Sprintf("History records (# of sent sms): %d\r\n",len(new(pfdatabase.DBtype).ShowHistory()))

	return container.NewVBox(		
		NewBoldLabel(m),
	)
}
func (a *theabout) buildUI() *fyne.Container {
	a.btnAppearance = &widget.Button{Text: "Change appearance!", OnTapped: func() {
			dialog.NewCustom("Fix the looks for the application!","Close", 
			appearance.NewSettings().LoadAppearanceScreen(a.window),
			a.window).Show()
		},
	}
	a.dodebug = &widget.RadioGroup{Options: []string{"Yes","No"},
		Horizontal: true, 
		Required: true, 
		OnChanged: func(v string) {
			oldsetting:=fyne.CurrentApp().Preferences().Bool("debug")
			if oldsetting==(v=="Yes") || !oldsetting==(v=="No") {
				return
			}
			newsetting:=!oldsetting
			m:=""
			m +="Application needs to restart after Debug setting has been changed.\r\n"
			m +="Please note that now all files pfsms uses that are now located in:\r\n"
			m +=GetHomeDir()+ ".\r\n"
			fyne.CurrentApp().Preferences().SetBool("debug",newsetting)
			Setupfiles()
			m +="will change and instead be located in: \r\n"
			m += GetHomeDir()+".\r\n"
			fyne.CurrentApp().Preferences().SetBool("debug",oldsetting)
			Setupfiles()
			dialog.NewConfirm("Really change Debug setting?", m, func(confirmed bool) {
				if !confirmed {
					if v=="Yes" {
						a.dodebug.SetSelected("No")
					} else {
						a.dodebug.SetSelected("Yes")
					}
					return
				} else {
					fyne.CurrentApp().Preferences().SetBool("debug",newsetting)
					Setupfiles()
					fyne.CurrentApp().Quit()
				}}, a.window).Show()
		},
	}
	if fyne.CurrentApp().Preferences().Bool("debug") {
		a.dodebug.SetSelected("Yes")
	} else  {
		a.dodebug.SetSelected("No")
	}
	return container.NewVBox(
		&widget.Card{Title: "App Info", Content: container.NewHBox(		
			layout.NewSpacer(),a.abouttext(),layout.NewSpacer(),
		)},
		&widget.Card{Title: "Database", Content: container.NewHBox(
			layout.NewSpacer(),a.aboutdatabase(),layout.NewSpacer())},
		&widget.Card{Title: "Interface", Content: 
			container.NewHBox(
				layout.NewSpacer(),
				container.NewVBox(
					container.NewGridWithColumns(2,NewBoldLabel("Debug?"), a.dodebug),
					layout.NewSpacer(),
					a.btnAppearance),
				layout.NewSpacer(),
			),
		},
		layout.NewSpacer(),
	)
}
func (s *theabout) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About pfsms", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}
