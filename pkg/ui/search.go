package ui

import "github.com/rivo/tview"

type SearchBox struct {
	Widget
	view *tview.InputField
}

func NewSearchBox(app *UI) (w *SearchBox) {
	w = &SearchBox{}
	w.app = app
	w.Name = "Search"

	w.view = tview.NewInputField()
	w.view.SetTitle("Search please")

	return

}

func (w *SearchBox) View() tview.Primitive {
	return w.view
}

func (w *SearchBox) Render(g *tview.Grid) error {
	w.view.SetTitle("Search please")
	return nil
}
