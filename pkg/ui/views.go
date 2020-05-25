package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

func (ui *UI) renderResultsView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := "search_results"

	v, err := g.SetView(name, 1, 5, maxX-30, maxY-7)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		ui.log("creating results view")
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

func (ui *UI) buildStatusMessage() (string, error) {
	// ui.log("getting player status")
	status := ui.Player.GetStatus()
	item := status.currentItem
	// ui.log(status)
	if item == nil {
		return "...", nil
	}

	var progressMsg string

	positionInt := int(status.currentPosition)
	durationInt := int(status.currentDuration)

	progressMsg = fmt.Sprintf("%d", positionInt) + "/" + fmt.Sprintf("%d", durationInt)

	statusMsg := progressMsg + " -- " + item.RawTitle
	return statusMsg, nil

}

func (ui *UI) buildProgressBar(v *gocui.View) {
	width, _ := v.Size()
	status := ui.Player.GetStatus()
	if status.currentItem == nil {
		return

	}
	position := status.currentPosition
	duration := status.currentDuration

	if int(duration) == 0 {
		ui.log("Duration 0, skipping bar")
		return
	}

	progressPercent := position / duration

	var barPosition float64

	barPosition = float64(width) * progressPercent

	barChars := strings.Repeat("+", int(barPosition))
	v.Clear()
	fmt.Fprintln(v, barChars)

}

func (ui *UI) renderProgressBar(g *gocui.Gui) error {
	maxX, maxY := ui.g.Size()
	name := "progress_bar"

	v, err := ui.g.SetView(name, 1, maxY-3, maxX-30, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.log("creating progress bar")
		v.Clear()
		v.Wrap = false
		v.Editable = false
		v.Frame = true

	}
	v.Clear()
	ui.buildProgressBar(v)

	return nil
}

func (ui *UI) renderStatusView(g *gocui.Gui) error {
	maxX, maxY := ui.g.Size()
	name := "status_view"

	v, err := ui.g.SetView(name, 1, maxY-6, maxX-30, maxY-4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.log("creating Status view")
		v.Clear()
		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Status"

		ui.State.views = append(ui.State.views, name)
		ui.State.curView = len(ui.State.views) - 1
		ui.State.idxView += 1

	}
	v.Clear()
	msg, err := ui.buildStatusMessage()
	if err != nil {
		fmt.Fprintln(v, err.Error())
	}
	fmt.Fprintln(v, msg)

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

		ui.log("creating log view")
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

		ui.log("creating subreddits view")
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
	ui.log("switching to next view")
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
