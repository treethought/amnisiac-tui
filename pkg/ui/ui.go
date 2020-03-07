package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	r "github.com/treethought/amnisiac/pkg/reddit"
	t "github.com/treethought/amnisiac/pkg/types"
)

// Gui wraps the gocui Gui object which handles rendering and events
type UI struct {
	g      *gocui.Gui
	State  uiState
	Player PlayerController
}

// guiState stores internal state of resources
type uiState struct {
	ResultBuffer map[string]*t.Item
	views        []string
	curView      int
	idxView      int
}

func NewGui() (*UI, error) {
	fmt.Println("Creating app")

	initialState := uiState{
		ResultBuffer: map[string]*t.Item{},
	}

	mpvPlayer := NewMPVController()

	ui := &UI{
		State:  initialState,
		Player: mpvPlayer,
	}

	return ui, nil

}

func (ui *UI) Start() error {

	err := ui.Player.Initialize()
	if err != nil {
		return err
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	ui.g = g

	defer ui.Teardown()

	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan
	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	err = ui.initializeLayout()
	if err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return nil

}

func (ui *UI) Teardown() {
	ui.Player.Shutdown()
	ui.g.Close()

}

func StartApp() {

	gui, err := NewGui()
	if err != nil {
		panic(err)
	}
	gui.Start()

}

func (ui *UI) initializeLayout() error {
	if err := ui.statusView(ui.g); err != nil {
		log.Panicln(err)
	}

	var emptyResult []*t.Item

	if err := ui.populateSearchResults(emptyResult); err != nil {
		log.Panicln(err)

	}
	if err := logView(ui); err != nil {
		log.Panicln(err)
	}
	if err := subredditView(ui, ""); err != nil {
		log.Panicln(err)
	}

	if err := ui.initKeybindings(); err != nil {
		log.Panicln(err)
	}

	if _, err := ui.g.SetCurrentView("sub_list"); err != nil {
		log.Panicln(err)
	}
	return nil

}

func layout(g *gocui.Gui) error {

	return nil
}

func logView(ui *UI) error {
	maxX, maxY := ui.g.Size()
	name := "info_view"
	v, err := ui.g.SetView(name, 0, maxY-3, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Title = "Log"
		v.Highlight = true

	}

	ui.State.views = append(ui.State.views, name)
	ui.State.curView = len(ui.State.views) - 1
	ui.State.idxView += 1

	return nil

}

func subredditView(ui *UI, filter string) error {
	maxX, maxY := ui.g.Size()
	name := "sub_list"
	v, err := ui.g.SetView(name, maxX-25, 0, maxX-1, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Title = "Subreddits"
		v.Highlight = true

	}

	subs, err := r.SubRedditsFromWiki("Music", "musicsubreddits")
	if err != nil {
		fmt.Fprintln(v, "Failed to fetch subs", err)
	}
	for _, sub := range subs {
		fmt.Fprintln(v, sub)
	}

	ui.State.views = append(ui.State.views, name)
	ui.State.curView = len(ui.State.views) - 1
	ui.State.idxView += 1

	return nil

}

func (ui *UI) nextView(disableCurrent bool) error {
	next := ui.State.curView + 1
	if next > len(ui.State.views)-1 {
		next = 0
	}

	if _, err := ui.g.SetCurrentView(ui.State.views[next]); err != nil {
		return err
	}

	ui.State.curView = next
	return nil
}

func (ui *UI) writeLog(a ...interface{}) error {
	v, err := ui.g.View("log_view")
	if err != nil {
		return err
	}
	fmt.Fprintln(v, a)

	return nil

}

func (ui *UI) PlayTrack(gui *gocui.Gui, v *gocui.View) error {

	selectedLine := ui.GetSelectedContent(gui, v)

	item := ui.State.ResultBuffer[selectedLine]

	err := ui.Player.PlayTrack(item)
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) initKeybindings() error {
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			ui.Teardown()
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	if err := ui.g.SetKeybinding("sub_list", gocui.KeyEnter, gocui.ModNone, ui.doSearch); err != nil {
		return err
	}
	if err := ui.g.SetKeybinding("search_results", gocui.KeyEnter, gocui.ModNone, ui.PlayTrack); err != nil {
		return err
	}
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlSlash, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return ui.statusView(g)
		}); err != nil {
	}
	if err := ui.g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return ui.nextView(true)
		}); err != nil {
		return err
	}

	if err := ui.g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, cursorDown); err != nil {
	}
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, cursorUp); err != nil {
	}

	return nil
}
