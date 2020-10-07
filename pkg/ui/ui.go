package ui

import (
	"log"

	"gitlab.com/tslocum/cview"

	"github.com/gdamore/tcell/v2"
	"github.com/treethought/amnisiac/pkg/logger"
	"github.com/treethought/amnisiac/pkg/player"
	"github.com/treethought/amnisiac/pkg/types"
)

// UI wraps the cview application which handles rendering and events
type UI struct {
	// Underlying Gui object to render views
	app *cview.Application

	Widgets []WidgetRenderer
	State   uiState
	Player  player.PlayerController
	Logger  *log.Logger
	grid    *cview.Grid
}

// guiState stores internal state of resources
type uiState struct {
	curView        int
	selectedSource string
	resultItems    []*types.Item
	sources        []string
	message        string
}

func NewApp() *UI {
	tapp := cview.NewApplication()
	app := &UI{
		app: tapp,
	}
	app.Player = player.NewMPVController()
	go app.Player.Initialize()

	app.State.message = "greetings"
	app.Logger = logger.GetLoggerInstance()

	return app
}
func newPrimitive(text string) cview.Primitive {
	return cview.NewTextView().
		SetTextAlign(cview.AlignCenter).
		SetText(text)
}

func (ui *UI) render() {
	for _, w := range ui.Widgets {
		w.Render(ui.grid)
	}

}

func (ui *UI) initWidgets() {
	sources := NewSourceList(ui)
	ui.Widgets = append(ui.Widgets, sources)

	search := NewSearchBox(ui)
	ui.Widgets = append(ui.Widgets, search)

	status := NewStatus(ui)
	ui.Widgets = append(ui.Widgets, status)

	results := NewItemList(ui)
	ui.Widgets = append(ui.Widgets, results)

	player := NewProgressBar(ui)
	ui.Widgets = append(ui.Widgets, player)

	ui.grid = cview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true)

	ui.grid.
		AddItem(sources.view, 1, 3, 3, 2, 0, 0, false).
		AddItem(search.view, 0, 0, 1, 3, 0, 0, false).
		AddItem(player.view, 3, 0, 1, 3, 0, 0, false).
		AddItem(results.view, 1, 0, 2, 3, 0, 0, false).
		AddItem(status.view, 0, 3, 1, 2, 0, 0, false)

	ui.render()

	ui.app.SetRoot(ui.grid, true).SetFocus(ui.grid)

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			if ui.State.curView == len(ui.Widgets)-1 {
				ui.State.curView = -1
			}
			ui.State.curView += 1
			widget := ui.Widgets[ui.State.curView]

			ui.app.SetFocus(widget.View())
			ui.render()

			return nil

		}

		return event
	})
	ui.Logger.Print("Initialized widgets")

}

func Start() {

	ui := NewApp()
	ui.initWidgets()

	err := ui.app.Run()
	if err != nil {
		panic(err)
	}
}
