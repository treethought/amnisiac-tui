package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type SourceList struct {
	Widget
	view *tview.List
}

func NewSourceList(app *UI) *SourceList {
	w := &SourceList{}
	w.app = app
	w.Name = "Sources"

	w.view = tview.NewList()

	w.view.SetTitle("Sources")
	// .SetMouseCapture(w.H
	// w.view.SetInputCapture()

	return w
}

func (w *SourceList) View() tview.Primitive {
	return w.view
}

func (w *SourceList) Render(grid *tview.Grid) (err error) {

	w.view.Clear()
	w.view.AddItem("Subreddit 1", "", 'a', nil)
	w.view.AddItem("Reddit", "Musical Subreddits", 'r', nil)
	w.view.AddItem("Soundcloud", "The worlds music thing", 's', nil)
	w.view.AddItem("Bandcamp", "The other worlds music thing", 'b', nil)

	return

}

func (w *SourceList) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	key := event.Key()
	switch key {
	case tcell.KeyEnter:
		return nil
	}
	return event
}
