package ui
import (
	"fmt"
	"github.com/jroimartin/gocui"
	t "github.com/treethought/amnisiac/pkg/types"
)

// func subredditView(g *gocui.Gui)
func populateSearchResults(g *gocui.Gui, results []t.Item, query string) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 0, 5, maxX-5, maxY-5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.SetCursor(1, 6)
		v.Wrap = false
		v.Editable = false
		v.Autoscroll = true
		v.EditNewLine()
		v.EditNewLine()

		fmt.Fprintln(v, "Results:", len(results))
		x, y := 1, 7

		for _, i := range results[:10] {
			v, err := g.SetView(i.ID, x, y, maxX-6, y+2)
			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
			}
			fmt.Fprintln(v, i.RawTitle)
			y += 3
		}
	}
    if _, err := g.SetCurrentView(name); err != nil {
        return err
    }

    views = append(views, name)
	curView = len(views) - 1
	idxView += 1

	return nil
}

