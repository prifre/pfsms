package ui

import (
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	appearance "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/layout"
)

type theabout struct {
	dodebug				*widget.RadioGroup
	window  		    fyne.Window
	app        			fyne.App
}

func NewAbout(a fyne.App, w fyne.Window) *theabout {
	return &theabout{app: a, window: w}
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
	bin := "pfsms"
	var dtg string
	fi, err := os.Stat(bin)
	if err == nil {
		dtg = "on " + fi.ModTime().Format("2006-01-02.15:04:05")
	}
	m := fmt.Sprintf("program %q, compiled with go %s %s\n",bin, runtime.Version(), dtg)
	m += "\r\n"+Getmemoryinfo()
	return container.NewVBox(		
		NewBoldLabel("PFSMS"), 
		NewBoldLabel("by Peter Freund"), 
		NewBoldLabel("prifre@prifre.com"), 
		NewBoldLabel(m),
	)
}
func (a *theabout) buildUI() *fyne.Container {
	interfaceContainer := appearance.NewSettings().LoadAppearanceScreen(a.window)
	a.dodebug = &widget.RadioGroup{Options: []string{"Yes","No"}, Horizontal: true, Required: true, OnChanged: func(v string) {
		fyne.CurrentApp().Preferences().SetBool("debug",v=="Yes")
		Setupfiles()
		new(pfsettings).buildFilePart()
	}}
	if fyne.CurrentApp().Preferences().Bool("debug") {
		a.dodebug.SetSelected("Yes")
	} else  {
		a.dodebug.SetSelected("No")
	}
	return container.NewVBox(
		&widget.Card{Title: "App Info", Content: container.NewHBox(		
			layout.NewSpacer(),
			a.abouttext(),
			container.NewVBox(NewBoldLabel("Debug?"), a.dodebug),
			layout.NewSpacer(),
		)},
		layout.NewSpacer(),

		&widget.Card{Title: "User Interface", Content: interfaceContainer},
	)
}
func (s *theabout) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About pfsms", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}
