package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	t "github.com/treethought/amnisiac/pkg/types"
)

var (
	resultMap = map[string]t.Item{}
)

// func subredditView(g *gocui.Gui)
func populateSearchResults(g *gocui.Gui, results []t.Item) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 0, 5, maxX-50, maxY-5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Wrap = false
		v.Editable = false
		v.Frame = true
		v.Title = "Results"

	}

	v.Clear()
	for _, item := range results {
		fmt.Fprintln(v, item.RawTitle)
		resultMap[item.RawTitle] = item
	}

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	views = append(views, name)
	curView = len(views) - 1
	idxView += 1

	return nil
}
