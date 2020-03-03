package ui

import (
	"fmt"
	"github.com/blang/mpv"
	"github.com/jroimartin/gocui"
	"os/exec"
)

func StartMPV() *exec.Cmd {
    cmd := exec.Command("mpv", "--idle", "--input-ipc-server=/tmp/mpvsocket", "--no-video")
    cmd.Start()
    return cmd
}

func PlayTrack(g *gocui.Gui, v *gocui.View) error {

	selected_title := GetSelectedContent(g, v)

	item := resultMap[selected_title]

	statusv, _ := g.View("status_view")

	fmt.Fprintln(statusv, "Playing", item.URL)
	ipcc := mpv.NewIPCClient("/tmp/mpvsocket") // Lowlevel client
	c := mpv.NewClient(ipcc)                   // Highlevel client, can also use RPCClient

	err := c.Loadfile(item.URL, mpv.LoadFileModeReplace)
	if err != nil {
		fmt.Fprintln(statusv, err)
		return err
	}
	return nil
}

