package ui

import "github.com/jroimartin/gocui"

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

func (ui *UI) doSearch(g *gocui.Gui, v *gocui.View) (err error) {
	ui.writeLog("Performing search")

	sub_v, err := g.View("sub_list")
	if err != nil {
		return err
	}

	selectedSub := ui.GetSelectedContent(sub_v)

	go ui.searchAndDisplayResults(selectedSub)
	return nil

}
