package ui

import (
	"github.com/blang/mpv"
	t "github.com/treethought/amnisiac/pkg/types"
	"os"
	"os/exec"
)

// StartMPV starts mpv in idle mode and specifies the ipc socket
func StartMPV() (*exec.Cmd, error) {
	cmd := exec.Command("mpv", "--idle", "--input-ipc-server=/tmp/mpvsocket", "--no-video")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

type PlayerController interface {
	Initialize() error
	Shutdown() error
	PlayTrack(item *t.Item) error
	TogglePause() error
	GetPosition() (int32, error)
	Seek(int32) error
}

type MPVController struct {
	client  *mpv.Client
	queue   map[int]t.Item
	process *exec.Cmd
	logger  *log.Logger
}

// NewMPVController creates a new instance of an MPV Client satisfying the PlayerController interface
func NewMPVController() *MPVController {
	m := MPVController{
		queue: map[int]t.Item{},
	}
	return &m

}

func (m *MPVController) Initialize() error {
	m.logger = GetLoggerInstance()
	process, err := StartMPV()
	if err != nil {
		return err
	}
	m.process = process
	ipcc := mpv.NewIPCClient("/tmp/mpvsocket") // Lowlevel client
	c := mpv.NewClient(ipcc)                   // Highlevel client, can also use RPCClient
	m.client = c

	return nil

}

func (m *MPVController) Shutdown() error {

	err := m.process.Process.Signal(os.Kill)
	if err != nil {
		return err
	}
	return nil

}

func (m *MPVController) PlayTrack(item *t.Item) error {
	m.logger.Println("Playing track", item.RawTitle)

	err := m.client.Loadfile(item.URL, mpv.LoadFileModeReplace)
	if err != nil {
		return err
	}
	return nil

}

func (m *MPVController) TogglePause() error {
    paused, err := m.client.Pause()
    if err != nil {
        return err
    }
    if paused {
        err := m.client.SetPause(false)
        return err
        
    } else {
        err := m.client.SetPause(true)
        return err
    }

	return err
}
	m.logger.Println("Pausing current track", m.status.currentItem)

func (m *MPVController) GetPosition() (int32, error) {
	return 0, nil
}

func (m *MPVController) Seek(int32) error {
	return nil
}
