package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

type Status struct {
	Widget
	view *tview.TextView
}

func NewStatus(app *UI) *Status {
	w := &Status{}
	w.app = app
	w.Name = "Status"

	w.view = tview.NewTextView()

	w.view.SetTitle("Status")
	w.view.SetText("Greeting")

	return w
}

func (w *Status) View() tview.Primitive {
	return w.view
}

func (w *Status) Render(grid *tview.Grid) (err error) {
	w.view.Clear()

	msg := fmt.Sprintf("source: %s | view: %d", w.app.State.selectedSource, w.app.State.curView)
	w.view.SetText(msg)
	return

}
