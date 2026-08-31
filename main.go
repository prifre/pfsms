package main

import (
	_ "embed"

	"pfsms/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed resources/pfsms.png
var appIconBytes []byte

func main() {
	a := app.NewWithID("com.prifre.pfsms")

	// Sätt applikationens ikon från inbäddad bild
	if len(appIconBytes) > 0 {
		iconRes := fyne.NewStaticResource("pfsms.png", appIconBytes)
		a.SetIcon(iconRes)
	}

	w := a.NewWindow("pfsms")

	// 1. Sätt innehållet
	w.SetContent(ui.Create(w))

	// 2. Sätt storleken på fönstret
	w.Resize(fyne.NewSize(1024, 768))

	// Centrera fönstret på skärmen
	w.CenterOnScreen()

	// 3. Visa och starta applikationen
	w.ShowAndRun()
}
