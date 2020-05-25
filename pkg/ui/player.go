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

func (m *MPVController) log(msgs ...interface{}) {
	m.logger.Println(msgs...)

}

func (m *MPVController) Initialize() error {
	m.logger = GetLoggerInstance()

	m.log("Starting MPV process")
	process, err := StartMPV()
	if err != nil {
		panic(err)
	}
	m.process = process

	m.log("initializing clients")
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
	m.log("Playing track", item.RawTitle)

	err := m.client.Loadfile(item.URL, mpv.LoadFileModeReplace)
	if err != nil {
		return err
	}
	m.status.currentItem = item
	m.status.playState = "playing"

	return nil

}

func (m *MPVController) TogglePause() error {
	m.log("Pausing current track", m.status.currentItem)
	paused, err := m.client.Pause()
	if err != nil {
		return err
	}
	if paused {
		err := m.client.SetPause(false)
		m.status.playState = "playing"
		return err

	} else {
		err := m.client.SetPause(true)
		m.status.playState = "paused"
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
	if m.client == nil {
		m.log("Client not ready")
		return m.status
	}

	paused, err := m.client.Pause()
	if err != nil {
		m.log("Failed to get playing state!!")
		m.log("Player error: ", err.Error())
	}
	if paused {
		m.status.playState = "paused"

	} else {
		m.status.playState = "playing"
	}

	pos, err := m.GetPosition()
	if err != nil {
		m.log("Failed to get position!!")
		m.log("Player error: ", err.Error())
	} else {
		m.status.currentPosition = pos
	}

	l, err := m.client.Duration()
	if err != nil {
		m.log("Failed to get duration!!")
		m.log("Player error: ", err.Error())
	} else {

		m.status.currentDuration = l
	}

	return m.status

}

func (m *MPVController) GetPosition() (float64, error) {
	return m.client.Position()
}

func (m *MPVController) Seek(pos int) error {
	return m.client.Seek(pos, mpv.SeekModeAbsolute)
}
