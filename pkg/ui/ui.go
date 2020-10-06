package ui

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/treethought/amnisiac/pkg/player"
	"github.com/treethought/amnisiac/pkg/types"
)

// UI wraps the tview application which handles rendering and events
type UI struct {
	// Underlying Gui object to render views
	app *tview.Application

	Widgets []WidgetRenderer
	State   uiState
	Player  player.PlayerController
	Logger  *log.Logger
	grid    *tview.Grid
}

// guiState stores internal state of resources
type uiState struct {
	curView        int
	selectedSource string
	resultItems    []*types.Item
}

func NewApp() *UI {
	tapp := tview.NewApplication()
	app := &UI{
		app: tapp,
	}

	return app
}
func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
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

	menu := newPrimitive("Menu")

	ui.grid = tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true)

	ui.grid.AddItem(menu, 3, 0, 0, 1, 0, 0, false).
		AddItem(sources.view, 1, 3, 3, 2, 0, 0, false).
		AddItem(search.view, 0, 0, 1, 3, 0, 0, false).
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

}

func Start() {

	ui := NewApp()
	ui.initWidgets()

	err := ui.app.Run()
	if err != nil {
		panic(err)
	}
}
