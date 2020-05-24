package ui

import (
	"fmt"

	r "github.com/treethought/amnisiac/pkg/reddit"
	"github.com/treethought/amnisiac/pkg/types"
	t "github.com/treethought/amnisiac/pkg/types"
)

func (ui *UI) searchAndDisplayResults(subreddits ...string) error {
	ui.writeLog("Fetching items from", subreddits)

	var items []*types.Item
	for _, s := range subreddits {
		subItems, err := r.FetchItemsFromReddit(s)
		if err != nil {
			return err
		}
		for _, s := range subItems {
			items = append(items, s)
		}

		ui.writeLog("got items, populating")
		ui.populateSearchResults(items)

	}
	return nil
}

func (ui *UI) populateSubredditListing() error {
	v, err := ui.g.View("sub_list")
	if err != nil {
		return err
	}
	fmt.Fprintln(v, "Loading....")
	ui.updateUI()

	subs, err := r.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		ui.writeLog("Failed to fetch subs", err)
	}
	ui.writeLog("Subreddits retrieved")
	ui.updateUI()
	v.Clear()
	for _, sub := range subs {
		fmt.Fprintln(v, sub)
	}
	ui.updateUI()
	return err

}

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
