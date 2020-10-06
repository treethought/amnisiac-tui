package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/treethought/amnisiac/pkg/types"
)

type Item struct {
	Widget
}

type ItemList struct {
	Widget
	view  *tview.List
	items []types.Item
}

func NewItemList(app *UI) *ItemList {
	w := &ItemList{}
	w.app = app
	w.Name = "Items"

	w.view = tview.NewList()

	w.view.SetTitle("Items")
	w.view.AddItem("Testing", "test", '-', nil)
	w.view.SetInputCapture(w.HandleInput)

	return w
}

func (w *ItemList) View() tview.Primitive {
	return w.view

}

func (w *ItemList) Render(g *tview.Grid) (err error) {
	w.view.Clear()

	for i, item := range w.app.State.resultItems {
		w.view.AddItem(item.RawTitle, item.Domain, rune(i), nil)
	}
	return

}

func (w *ItemList) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	_ = w.view.GetCurrentItem()
	key := event.Key()
	switch key {
	case tcell.KeyEnter:
		cur := w.view.GetCurrentItem()
		item := w.app.State.resultItems[cur]

		err := w.app.Player.PlayTrack(item)
		if err != nil {
			panic(err)
		}
		return event

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
