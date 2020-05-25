package ui

import (
	"log"
	"os"
	"os/exec"

	"github.com/blang/mpv"
	t "github.com/treethought/amnisiac/pkg/types"
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
	GetPosition() (float64, error)
	Seek(int) error
	QueueTrack(item *t.Item) error
	GetStatus() PlayerStatus
}

type PlayerStatus struct {
	currentItem          *t.Item
	currentPosition      float64
	currentDuration      float64
	currentQueuePosition int
	currentQueue         map[int]t.Item
	playState            string
}

type MPVController struct {
	client  *mpv.Client
	queue   map[int]t.Item
	process *exec.Cmd
	status  PlayerStatus
	logger  *log.Logger
}

// NewMPVController creates a new instance of an MPV Client satisfying the PlayerController interface
func NewMPVController() *MPVController {
	m := MPVController{
		queue:  map[int]t.Item{},
		status: PlayerStatus{},
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
	m.status.currentItem = item
	m.status.playState = "playing"

	return nil

}

func (m *MPVController) TogglePause() error {
	m.logger.Println("Pausing current track", m.status.currentItem)
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

}

func (m *MPVController) QueueTrack(item *t.Item) error {
	err := m.client.Loadfile(item.URL, mpv.LoadFileModeAppendPlay)
	if err != nil {
		return err
	}

	return nil
}

func (m *MPVController) GetStatus() PlayerStatus {
	m.logger.Println("building status")

	return m.status

}

func (m *MPVController) GetPosition() (float64, error) {
	return m.client.Position()
}

func (m *MPVController) Seek(pos int) error {
	return m.client.Seek(pos, mpv.SeekModeAbsolute)
}
