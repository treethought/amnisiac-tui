package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/treethought/amnisiac/pkg/types"
	"gitlab.com/tslocum/cview"
)

type ItemList struct {
	Widget
	view  *cview.List
	items []types.Item
}

func NewItemList(app *UI) *ItemList {
	w := &ItemList{}
	w.app = app
	w.Name = "Items"

	w.view = cview.NewList()

	w.view.SetTitle("Items")
	w.view.SetInputCapture(w.HandleInput)
	w.view.SetBackgroundColor(tcell.ColorDefault)

	return w
}

func (w *ItemList) View() cview.Primitive {
	return w.view

}

func (w *ItemList) Render(g *cview.Grid) (err error) {
	w.view.Clear()

	for _, item := range w.app.State.resultItems {
		li := cview.NewListItem(item.RawTitle)
		li.SetSecondaryText(item.Domain)
		li.SetReference(item)
		w.view.AddItem(li)
	}
	return

}

func (w *ItemList) HandleInput(event *tcell.EventKey) *tcell.EventKey {

	_ = w.view.GetCurrentItem()
	key := event.Key()
	switch key {
	case tcell.KeyEnter:
		cur := w.view.GetCurrentItem()
		ref := cur.GetReference()
		item, ok := ref.(*types.Item)
		if !ok {
			log.Fatal("selected item ref was not an item")
		}

		w.app.Logger.Printf("Selected item %s\n", item.RawTitle)

		w.app.State.message = fmt.Sprintf("Playing %s", item.RawTitle)
		go func() {
			err := w.app.Player.PlayTrack(item)
			if err != nil {
				w.app.State.message = "Failed to play track"
			}
			w.app.app.QueueUpdateDraw(func() {})

		}()
		return nil

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
