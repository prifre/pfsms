package ui

import (
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

const version = "v0.0.2"

type about struct {
}

func NewAbout() *about {
	return &about{}
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
func (a *about) buildUI() *fyne.Container {
	m:=Getmemoryinfo()
	return container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(		
			layout.NewSpacer(),
			container.NewVBox(		
				newBoldLabel("PFSMS"), 
				layout.NewSpacer(),
				newBoldLabel(version),
				layout.NewSpacer(),
				newBoldLabel(m),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}
func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About", Content: a.buildUI()}
}
