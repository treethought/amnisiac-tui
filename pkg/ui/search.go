package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
)

func doSearch(g *gocui.Gui, v *gocui.View) (err error) {

    sub_v, err := g.View("sub_list")
    if err != nil {
        return err
    }

    selectedSub := GetSelectedContent(g, sub_v)


    sv, err := g.View("status_view")
    if err != nil {
        return err
    }
    fmt.Fprintln(sv, "Fetching items from", selectedSub)


	items, err := r.FetchItemsFromReddit(selectedSub)
	if err != nil {
		return err
	}

	return populateSearchResults(g, items)

}




func statusView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	name := "status_view"
	v, err := g.SetView(name, 0, 0, maxX-30, 2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Wrap = true
		v.Editable = true
		v.Frame = true

	}

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	views = append(views, name)
	curView = len(views) - 1
	idxView += 1

	return nil
}

func simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
}
