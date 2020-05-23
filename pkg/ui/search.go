package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
)

func (ui *UI) doSearch(g *gocui.Gui, v *gocui.View) (err error) {

	sub_v, err := g.View("sub_list")
	if err != nil {
		return err
	}

	selectedSub := ui.GetSelectedContent(g, sub_v)

	sv, err := g.View("status_view")
	if err != nil {
		return err
	}
	fmt.Fprintln(sv, "Fetching items from", selectedSub)

	items, err := r.FetchItemsFromReddit(selectedSub)
	if err != nil {
		return err
	}

	return ui.populateSearchResults(items)

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

	subs, err := r.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		fmt.Fprintln(v, "Failed to fetch subs", err)
	}
	for _, sub := range subs {
		fmt.Fprintln(v, sub)
	}
	return err

}
