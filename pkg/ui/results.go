package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	t "github.com/treethought/amnisiac/pkg/types"
)

// func subredditView(g *gocui.Gui)
func populateSearchResults(g *gocui.Gui, results []t.Item) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 0, 5, maxX-5, maxY-5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.SetCursor(1, 6)
		// v.SetCursor(0, 0)
		v.Highlight = true
		// v.BgColor = gocui.ColorMagenta
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorYellow
		v.Wrap = false
		v.Editable = false
		v.Autoscroll = true
		v.Frame = true
		v.Title = "Results"

	}

	for _, item := range results {
		// v, err := g.SetView(item.ID, x, y, maxX-6, y+2)
		// if err != nil {
		// 	if err != gocui.ErrUnknownView {
		// 		return err
		// 	}
		// }
		fmt.Fprintln(v, item.RawTitle)
		fmt.Fprintln(v, "")
		// y += 3
	}

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	views = append(views, name)
	curView = len(views) - 1
	idxView += 1

	return nil
}
