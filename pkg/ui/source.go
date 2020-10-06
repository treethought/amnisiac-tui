package ui

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/treethought/amnisiac/pkg/reddit"
)

type SourceItem struct {
	name    string
	caption string
}

type SourceList struct {
	Widget
	items []string
	view  *tview.List
}

func NewSourceList(app *UI) *SourceList {
	w := &SourceList{}
	w.app = app
	w.Name = "Sources"

	subs, err := reddit.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		log.Fatal(err)
	}
	w.items = subs

	w.view = tview.NewList()

	w.view.SetTitle("Sources")
	w.view.SetInputCapture(w.HandleInput)
	w.view.ShowSecondaryText(false)

	return w
}

func (w *SourceList) fetchItems() error {

	idx := w.view.GetCurrentItem()
	s, _ := w.view.GetItemText(idx)
	items, err := reddit.FetchItemsFromReddit(s)
	if err != nil {
		return err
	}
	w.app.State.resultItems = items
	w.app.render()
	return nil

}

func (w *SourceList) View() tview.Primitive {
	return w.view
}

func (w *SourceList) Render(grid *tview.Grid) (err error) {
	w.view.Clear()
	for i, sub := range w.items {
		w.view.AddItem(sub, "", rune(i), nil)

	}

	return

}

func (w *SourceList) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	idx := w.view.GetCurrentItem()
	key := event.Key()
	switch key {
	case tcell.KeyEnter:
		err := w.fetchItems()
		if err != nil {
			panic(err)
		}

		s, _ := w.view.GetItemText(idx)
		w.app.State.selectedSource = s
		w.app.render()

	case tcell.KeyRune:
		switch event.Rune() {
		case 'g': // Home.
			w.view.SetCurrentItem(0)
		case 'G': // End.
			w.view.SetCurrentItem(-1)
		case 'j': // Down.
			cur := w.view.GetCurrentItem()
			w.view.SetCurrentItem(cur + 1)
		case 'k': // Up.
			cur := w.view.GetCurrentItem()
			w.view.SetCurrentItem(cur - 1)
		}

		return nil
	}

	return event
}
