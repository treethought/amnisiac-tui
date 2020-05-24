package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
	"github.com/treethought/amnisiac/pkg/types"
)

func (ui *UI) doSearch(g *gocui.Gui, v *gocui.View) (err error) {
	ui.writeLog("Performing search")

	sub_v, err := g.View("sub_list")
	if err != nil {
		return err
	}

	selectedSub := ui.GetSelectedContent(sub_v)

	go ui.searchAndDisplayResults(selectedSub)
	return nil

}

func (ui *UI) searchAndDisplayResults(subreddits ...string) error {
	ui.writeLog("Fetching items from", subreddits)
	// sv, err := g.View("status_view")
	// if err != nil {
	// 	return err
	// }
	// fmt.Fprintln(sv, "Fetching items from", subreddit)

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

func (ui *UI) populateSubredditListing() error {
	v, err := ui.g.View("sub_list")
	if err != nil {
		return err
	}
	fmt.Fprintln(v, "Loading....")
	ui.updateUI()

	subs, err := r.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		ui.writeLog(v, "Failed to fetch subs", err)
	}
	ui.writeLog(v, "Subreddits retrieved", err)
	ui.updateUI()
	v.Clear()
	for _, sub := range subs {
		fmt.Fprintln(v, sub)
	}
	ui.updateUI()
	return err

}
