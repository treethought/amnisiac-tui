package ui

import (
	"fmt"

	"gitlab.com/tslocum/cview"
)

type Status struct {
	Widget
	view *cview.TextView
}

func NewStatus(app *UI) *Status {
	w := &Status{}
	w.app = app
	w.Name = "Status"

	w.view = cview.NewTextView()

	w.view.SetTitle("Status")
	w.view.SetText("Greeting")

	return w
}

func (w *Status) View() cview.Primitive {
	return w.view
}

func (w *Status) Render(grid *cview.Grid) (err error) {
	w.view.Clear()

	statusMsg := w.app.State.message
	msg := fmt.Sprintf("source: %s | view: %d\n%s", w.app.State.selectedSource, w.app.State.curView, statusMsg)
	w.view.SetText(msg)
	return

}
