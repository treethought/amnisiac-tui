package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	logger "github.com/treethought/amnisiac/pkg/logger"
	player "github.com/treethought/amnisiac/pkg/player"
	t "github.com/treethought/amnisiac/pkg/types"
)

// UI wraps the gocui Gui object which handles rendering and events
type UI struct {
	g      *gocui.Gui
	State  uiState
	Player player.PlayerController
	Logger *log.Logger
}

// guiState stores internal state of resources
type uiState struct {
	ResultBuffer map[string]*t.Item
	views        []string
	curView      int
	idxView      int
}

func NewUI() (*UI, error) {
	fmt.Println("Creating app")

	initialState := uiState{
		ResultBuffer: map[string]*t.Item{},
	}

	mpvPlayer := player.NewMPVController()

	ui := &UI{
		State:  initialState,
		Player: mpvPlayer,
	}
	ui.Logger = logger.GetLoggerInstance()

	return ui, nil

}

func (ui *UI) Teardown() {
	ui.Player.Shutdown()
	ui.g.Close()

}

func (ui *UI) log(msgs ...interface{}) {
	ui.Logger.Println(msgs...)
}

func StartApp() {

	ui, err := NewUI()
	if err != nil {
		panic(err)
	}
	ui.log("***************************")
	ui.Start()

}

func (ui *UI) pollPlayerStatus() {
	ui.log("Starting status polling")
	for {
		time.Sleep(1 * time.Second)
		ui.renderStatusView(ui.g)
		ui.updateUI()
	}
}

// Start initializes a player and builds the UI
// serves as the entrypoint for the application
// by starting the gocui event loop
func (ui *UI) Start() error {

	ui.log("initializing player")
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

	go ui.pollPlayerStatus()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		// log.Panicln(err)
		return err

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

// updateUI forces a redraw of the views
// wrapper for *gocui.Gui.Update that must be called when
// changes are made that do not result from keybindings
func (ui *UI) updateUI() {
	ui.g.Update(ui.Layout)
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

	if err := ui.renderProgressBar(g); err != nil {
		log.Panicln(err)
	}

	return nil

}
