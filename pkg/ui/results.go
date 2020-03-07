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
func (ui *UI) populateSearchResults(results []*t.Item) error {
	maxX, maxY := ui.g.Size()
	name := "search_results"

	v, err := ui.g.SetView(name, 0, 5, maxX-50, maxY-5)
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
		ui.State.ResultBuffer[item.RawTitle] = item
	}

	if _, err := ui.g.SetCurrentView(name); err != nil {
		return err
	}

	ui.State.views = append(ui.State.views, name)
	ui.State.curView = len(ui.State.views) - 1
	ui.State.idxView += 1

	return nil
}
