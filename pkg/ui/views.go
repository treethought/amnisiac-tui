package ui

import "github.com/jroimartin/gocui"

func (ui *UI) renderResultsView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 1, 5, maxX-30, maxY-5)
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
	maxX, maxY := ui.g.Size()
	name := "status_view"

	v, err := ui.g.SetView(name, 1, maxY-3, maxX-30, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Status"

		ui.State.views = append(ui.State.views, name)
		ui.State.curView = len(ui.State.views) - 1
		ui.State.idxView += 1

	}

	return nil
}

func (ui *UI) renderLogView(g *gocui.Gui) error {
	maxX, _ := ui.g.Size()
	name := "log_view"
	v, err := g.SetView(name, 1, 0, maxX-30, 2)

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
	v, err := ui.g.SetView(name, maxX-25, 0, maxX-2, maxY-1)

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
