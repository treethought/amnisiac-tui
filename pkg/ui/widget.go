package ui

import "gitlab.com/tslocum/cview"

type Widget struct {
	app  *UI
	Name string
}

type WidgetRenderer interface {
	Render(grid *cview.Grid) error
	View() cview.Primitive
}
