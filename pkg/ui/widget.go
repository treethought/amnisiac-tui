package ui

import (
	"github.com/rivo/tview"
)

type Widget struct {
	app  *UI
	Name string
}

type WidgetRenderer interface {
	Render(grid *tview.Grid) error
	View() tview.Primitive
}
