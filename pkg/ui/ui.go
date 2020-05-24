package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	t "github.com/treethought/amnisiac/pkg/types"
)

// UI wraps the gocui Gui object which handles rendering and events
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

// Start initializes a player and builds the UI
// serves as the entrypoint for the application
// by starting the gocui event loop
func (ui *UI) Start() error {

	go ui.Player.Initialize()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan
	g.Cursor = true
	g.Mouse = true

	ui.g = g
	g.SetManager(ui)
	ui.writeLog("Manager set")

	if err := ui.initKeybindings(); err != nil {
		log.Panicln(err)
	}

	defer ui.Teardown()

	err = ui.initializeLayout()
	if err != nil {
		return err
	}
	ui.writeLog("Layout initialized")

	if _, err := ui.g.SetCurrentView("sub_list"); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return nil

}

func (ui *UI) initializeLayout() error {
	ui.writeLog("Initializing layout state")

	// render the base state
	ui.Layout(ui.g)

	var emptyResult []*t.Item

	if err := ui.populateSearchResults(emptyResult); err != nil {
		log.Panicln(err)
	}

	go ui.populateSubredditListing()

	ui.writeLog("UI initialized")
	return nil
}

// Layout updates the UI and keybindings on each event.
// This method allows UI to satisfy the gocui.Manager interface
// while wrapping the render updates with updating of app state
func (ui *UI) Layout(g *gocui.Gui) error {

	if err := ui.renderResultsView(g); err != nil {
		log.Panicln(err)
	}

	if err := ui.renderStatusView(g); err != nil {
		log.Panicln(err)
	}

	if err := ui.renderLogView(g); err != nil {
		log.Panicln(err)
	}
	if err := ui.renderSubredditView(g); err != nil {
		log.Panicln(err)
	}

	return nil

}

// updateUI forces a redraw of the views
// wrapper for *gocui.Gui.Update that must be called when
// changes are made that do not result from keybindings
func (ui *UI) updateUI() {
	ui.g.Update(ui.Layout)

}

func (ui *UI) renderResultsView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 0, 5, maxX-50, maxY-5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Clear()
		v.Highlight = true
		v.Wrap = false
		v.Editable = false
		v.Frame = true
		v.Title = "Results"

		ui.State.views = append(ui.State.views, name)
		ui.State.curView = len(ui.State.views) - 1
		ui.State.idxView += 1
	}
	return nil

}

func (ui *UI) renderStatusView(g *gocui.Gui) error {
	maxX, _ := g.Size()
	name := "status_view"

	v, err := g.SetView(name, 0, 0, maxX-30, 2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Wrap = true
		v.Editable = true
		v.Frame = true

		ui.State.views = append(ui.State.views, name)
		ui.State.curView = len(ui.State.views) - 1
		ui.State.idxView += 1

	}

	return nil
}

func (ui *UI) renderLogView(g *gocui.Gui) error {
	maxX, maxY := ui.g.Size()
	name := "log_view"
	v, err := ui.g.SetView(name, 0, maxY-3, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Title = "Log"
		v.Highlight = false
		v.Autoscroll = true

		ui.State.views = append(ui.State.views, name)
		ui.State.curView = len(ui.State.views) - 1
		ui.State.idxView += 1

	}

	return nil

}

func (ui *UI) renderSubredditView(g *gocui.Gui) error {
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

		ui.State.views = append(ui.State.views, name)
		ui.State.curView = len(ui.State.views) - 1
		ui.State.idxView += 1

	}

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

// writeLog writes the message to the log UI view
func (ui *UI) writeLog(a interface{}) error {
	v, err := ui.g.View("log_view")
	if err != nil {
		return err
	}
	fmt.Fprintln(v, a)
	ui.updateUI()

	return nil

}

func (ui *UI) SelectTrack(gui *gocui.Gui, v *gocui.View) error {

	selectedLine := ui.GetSelectedContent(v)

	item := ui.State.ResultBuffer[selectedLine]

	err := ui.Player.PlayTrack(item)
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) TogglePause(gui *gocui.Gui, v *gocui.View) error {
	err := ui.Player.TogglePause()
	return err
}
