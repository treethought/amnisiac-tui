package ui

import (
	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

type SearchBox struct {
	Widget
	view *cview.InputField
}

func NewSearchBox(app *UI) (w *SearchBox) {
	w = &SearchBox{}
	w.app = app
	w.Name = "Search"

	w.view = cview.NewInputField()
	w.view.SetBackgroundColor(tcell.ColorDefault)
	w.view.SetFieldBackgroundColor(tcell.ColorDefault)
	w.view.SetBorder(true)
	w.view.SetTitle("Search")
	w.view.SetText("Hit '/' to search")

	return

}

func (w *SearchBox) Selectable() bool { return true }

func (w *SearchBox) View() cview.Primitive {
	return w.view
}

func (w *SearchBox) Render(g *cview.Grid) error {
	return nil
}
