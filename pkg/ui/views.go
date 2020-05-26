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
	status := ui.Player.GetStatus()
	item := status.CurrentItem

	statusMsg := ""

	// ui.log(status)
	if item != nil {
		statusMsg += item.RawTitle
		statusMsg += " (" + item.SourcePlatform + ")"
		statusMsg += " (" + item.SubReddit + ")"
	}

	return statusMsg, nil

}

func (ui *UI) buildProgressBar(v *gocui.View) {
	width, _ := v.Size()
	status := ui.Player.GetStatus()
	if status.CurrentItem == nil {
		return

	}
	position := status.CurrentPosition
	duration := status.CurrentDuration

	if int(duration) == 0 {
		return
	}

	progressPercent := position / duration

	posFormatted := secondsToMinutes(int(position))
	durationFormatted := secondsToMinutes(int(duration))

	barContent := posFormatted

	var barPosition float64

	// barspace occupies the view between the positon and duration
	//
	barMaxWidth := width - len(posFormatted) - len(durationFormatted)
	barPosition = float64(barMaxWidth) * progressPercent

	barChars := strings.Repeat("=", int(barPosition))
	emptyChars := strings.Repeat(" ", (int(barMaxWidth) - int(barPosition)))

	barContent += barChars + emptyChars + durationFormatted

	v.Clear()
	fmt.Fprintln(v, barContent)

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
