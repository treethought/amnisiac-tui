package main
import (
	"fmt"
	"log"
	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
	t "github.com/treethought/amnisiac/pkg/types"
)

func main() {
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
	if err := PlayerView(g); err != nil {
		log.Panicln(err)
	}
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

func doSearch(g *gocui.Gui, v *gocui.View) (err error) {
	v.Clear()

	g.Cursor = false
	v.Editable = false

	items, err := r.FetchItemsFromReddit()
	if err != nil {
		return err
	}

	return searchResults(g, items)

}

func searchView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	name := "search_view"
	v, err := g.SetView(name, 0, 0, maxX-30, 2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Wrap = true
		fmt.Fprintln(v, "Search:")
		v.Editable = true

		v.MoveCursor(9, 0, true)

	}

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	return nil
}

// func subredditView(g *gocui.Gui)
func searchResults(g *gocui.Gui, results []t.Item) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 0, 5, maxX-5, maxY-5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.SetCursor(1, 6)
		v.Wrap = false
		v.Editable = false
		v.Autoscroll = true
		v.EditNewLine()
		v.EditNewLine()

		fmt.Fprintln(v, "Results:", len(results))
		x, y := 1, 7

		for _, i := range results[:10] {
			v, err := g.SetView(i.ID, x, y, maxX-6, y+2)
			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
			}
			fmt.Fprintln(v, i.RawTitle)
			y += 3
		}
	}

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	return nil

}
	return nil

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
