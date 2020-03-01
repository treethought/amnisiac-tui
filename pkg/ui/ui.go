package ui
import (
	"fmt"
	"github.com/jroimartin/gocui"
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

	g.SetManagerFunc(layout)

	if err := searchView(g); err != nil {
		log.Panicln(err)
	}

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {

	maxX, _ := g.Size()
	v, err := g.SetView("menu", maxX-25, 0, maxX-1, 3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "amnisiac")
		fmt.Fprintln(v, "^C: Exit")
	}

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

	if err := g.SetKeybinding("search_view", gocui.KeyEnter, gocui.ModNone, doSearch); err != nil {
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlSlash, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return searchView(g)
		}); err != nil {
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextView(g, true)
		}); err != nil {
		return err
	}
	return nil
}
