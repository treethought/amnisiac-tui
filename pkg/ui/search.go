package ui

import "gitlab.com/tslocum/cview"

type SearchBox struct {
	Widget
	view *cview.InputField
}

func NewSearchBox(app *UI) (w *SearchBox) {
	w = &SearchBox{}
	w.app = app
	w.Name = "Search"

	w.view = cview.NewInputField()
	w.view.SetTitle("Search please")
	w.view.SetText("Search")

	return

}

func (w *SearchBox) View() cview.Primitive {
	return w.view
}

func (w *SearchBox) Render(g *cview.Grid) error {
	w.view.SetTitle("Search please")
	return nil
}
