
package ui
import (
	"fmt"
	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
)


func doSearch(g *gocui.Gui, v *gocui.View) (err error) {
	v.Clear()

	g.Cursor = false
	v.Editable = false

	items, err := r.FetchItemsFromReddit()
	if err != nil {
		return err
	}

	return searchResults(g, items)

}

func searchView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	name := "search_view"
	v, err := g.SetView(name, 0, 0, maxX-30, 2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Wrap = true
		fmt.Fprintln(v, "Search:")
		v.Editable = true

		v.MoveCursor(9, 0, true)

	}

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

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
