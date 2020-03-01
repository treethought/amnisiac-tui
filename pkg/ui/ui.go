package ui
import (
	"fmt"
	"log"
	"github.com/jroimartin/gocui"
    // github.com/treethought/amnisiac/pkg/ui
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

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	if err := searchView(g); err != nil {
		log.Panicln(err)
	}

	maxX, _ := g.Size()
	v, err := g.SetView("menu", maxX-25, 0, maxX-1, 3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "amnisiac")
		fmt.Fprintln(v, "^C: Exit")
	}
	// if err := PlayerView(g); err != nil {
	// 	log.Panicln(err)
	// }
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
	return nil
}

