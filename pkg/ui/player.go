package ui

import (
	"fmt"
	"strings"

	"gitlab.com/tslocum/cview"

	"github.com/gdamore/tcell/v2"
)

type ProgressBar struct {
	Widget
	view *cview.TextView
}

func NewProgressBar(app *UI) (w *ProgressBar) {
	w = &ProgressBar{}
	w.app = app
	w.Name = "Search"

	w.view = cview.NewTextView()
	w.view.SetBorder(true)

	w.view.SetTitle("Player")
	w.view.SetBackgroundColor(tcell.ColorDefault)
	go w.buildProgressBar()

	return

}

func (w *ProgressBar) View() cview.Primitive {
	return w.view
}
func (w *ProgressBar) buildProgressBar() {
	for true {

		_, _, width, _ := w.view.GetRect()
		status := w.app.Player.GetStatus()
		if status.CurrentItem == nil {
			continue

		}
		position := status.CurrentPosition
		duration := status.CurrentDuration

		if int(duration) == 0 {
			continue
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

		info := fmt.Sprintf("%s ||  %s", w.app.State.selectedSource, status.CurrentItem.RawTitle)

		content := fmt.Sprintf("%s\n%s", barContent, info)

		w.view.Clear()
		w.view.SetText(content)

		w.app.app.QueueUpdateDraw(func() {})
	}

}

func (w *ProgressBar) Render(g *cview.Grid) error {
	w.app.Logger.Println("Rendering progress bar")

	return nil
}

func secondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%d:%d", minutes, seconds)
	return str
}
