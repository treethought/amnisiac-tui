package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/treethought/amnisiac/pkg/reddit"
	"gitlab.com/tslocum/cview"
)

type SourceItem struct {
	name    string
	caption string
}

type SourceList struct {
	Widget
	items []string
	view  *cview.List
}

func NewSourceList(app *UI) *SourceList {
	w := &SourceList{}
	w.app = app
	w.Name = "Sources"

	go w.fetchSubs()
	w.view = cview.NewList()

	w.view.SetTitle("Sources")
	w.view.SetInputCapture(w.HandleInput)
	w.view.ShowSecondaryText(false)

	return w
}

func (w *SourceList) fetchSubs() {
	subs, err := reddit.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		log.Fatal(err)
	}
	w.app.State.sources = subs

}

func (w *SourceList) fetchItems() error {

	selected := w.view.GetCurrentItem()
	s := selected.GetMainText()

	items, err := reddit.FetchItemsFromReddit(s)
	if err != nil {
		return err
	}
	w.app.State.resultItems = items
	w.app.render()
	return nil

}

func (w *SourceList) View() cview.Primitive {
	return w.view
}

func (w *SourceList) Render(grid *cview.Grid) (err error) {
	w.view.Clear()
	for _, sub := range w.app.State.sources {
		item := cview.NewListItem(sub)
		w.view.AddItem(item)

	}

	return

}

func (w *SourceList) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	key := event.Key()
	switch key {
	case tcell.KeyEnter:
		err := w.fetchItems()
		if err != nil {
			panic(err)
		}

		selected := w.view.GetCurrentItem()
		s := selected.GetMainText()
		w.app.State.selectedSource = s
		w.app.render()

	case tcell.KeyRune:
		switch event.Rune() {
		case 'g': // Home.
			w.view.SetCurrentItem(0)
		case 'G': // End.
			w.view.SetCurrentItem(-1)
		case 'j': // Down.
			cur := w.view.GetCurrentItemIndex()
			w.view.SetCurrentItem(cur + 1)
		case 'k': // Up.
			cur := w.view.GetCurrentItemIndex()
			w.view.SetCurrentItem(cur - 1)
		}

		return nil
	}

	return event
}
