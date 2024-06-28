package ui

type thetable struct {
	themeSelect *widget.Table
}

func newTable() *thetable{
	var data = [][]string{{"top left", "top right"},{"bottom left", "bottom right"}}
	list := widget.NewTable(
	func() (int, int) {
		return len(data), len(data[0])
	},
	func() fyne.CanvasObject {
		return widget.NewLabel("wide content")
	},
	func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(data[i.Row][i.Col])
	})
return 
