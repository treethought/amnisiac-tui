package main

import (
	"os"

	"github.com/treethought/amnisiac/pkg/player"
	"github.com/treethought/amnisiac/pkg/ui"
)

func main() {
	cmd, err := player.StartMPV()
	if err != nil {
		panic(err)
	}
	defer cmd.Process.Signal(os.Kill)

	// ui.StartApp()
	ui.Start()
}
