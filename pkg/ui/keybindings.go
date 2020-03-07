package ui
import (
    "github.com/jroimartin/gocui"
)

func (ui *UI) initKeybindings() error {
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
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
	if err := ui.g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, ui.TogglePause); err != nil {
		return err
	}
	if err := ui.g.SetKeybinding("status_view", gocui.KeyEnter, gocui.ModNone, ui.PlayTrack); err != nil {
		return err
	}
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlSlash, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
            g.SetCurrentView("status_view")
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
