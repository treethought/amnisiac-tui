package ui

import (
	"fmt"

	t "github.com/treethought/amnisiac/pkg/types"
)

var (
	resultMap = map[string]t.Item{}
)

// populateSearchResults replaces the results buffer with the current search results
func (ui *UI) populateSearchResults(results []*t.Item) error {
	ui.writeLog("populating search results")
	maxX, maxY := ui.g.Size()
	name := "search_results"

	v, err := ui.g.SetView(name, 0, 5, maxX-50, maxY-5)
	if err != nil {
		return err
	}

	v.Clear()
	for _, item := range results {
		fmt.Fprintln(v, item.RawTitle)
		ui.State.ResultBuffer[item.RawTitle] = item
	}

	// ui.updateUI()
	if _, err := ui.g.SetCurrentView(name); err != nil {
		return err
	}

	return nil
}
