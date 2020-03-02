package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
	t "github.com/treethought/amnisiac/pkg/types"
	"log"
)

var (
	views   = []string{}
	curView = -1
	idxView = 0
)

func StartApp() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan
	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := statusView(g); err != nil {
		log.Panicln(err)
	}

	var emptyResult []t.Item

	if err := populateSearchResults(g, emptyResult); err != nil {
		log.Panicln(err)

	}
	if err := subredditView(g, ""); err != nil {
		log.Panicln(err)
	}

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if _, err := g.SetCurrentView("sub_list"); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {

	return nil
}

func subredditView(g *gocui.Gui, filter string) error {
	maxX, maxY := g.Size()
	name := "sub_list"
	v, err := g.SetView(name, maxX-25, 0, maxX-1, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Title = "Subreddits"
		// v.Autoscroll = true
		v.Highlight = true

	}

	subs, err := r.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		fmt.Fprintln(v, "Failed to fetch subs", err)
	}
	for _, sub := range subs {
		fmt.Fprintln(v, sub)
	}

	views = append(views, name)
	curView = len(views) - 1
	idxView += 1

	return nil

}

func nextView(g *gocui.Gui, disableCurrent bool) error {
	next := curView + 1
	if next > len(views)-1 {
		next = 0
	}

	if _, err := g.SetCurrentView(views[next]); err != nil {
		return err
	}

	curView = next
	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("sub_list", gocui.KeyEnter, gocui.ModNone, doSearch); err != nil {
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlSlash, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return statusView(g)
		}); err != nil {
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextView(g, true)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, cursorDown); err != nil {
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, cursorUp); err != nil {
	}

	return nil
}

