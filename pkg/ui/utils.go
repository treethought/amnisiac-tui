package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (ui *UI) GetSelectedContent(v *gocui.View) string {
	_, cy := v.Cursor()

	lines := v.ViewBufferLines()
	selectedLine := lines[cy]

	return selectedLine
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeLog writes the message to the log UI view
func (ui *UI) writeLog(a ...interface{}) error {
	ui.log(a...)
	return nil

}

func (ui *UI) writeStatus(a ...interface{}) error {

	v, err := ui.g.View("status_view")
	if err != nil {
		return err
	}
	fmt.Fprintln(v, a...)
	ui.updateUI()

	return nil

}
