package player

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/blang/mpv"
	logger "github.com/treethought/amnisiac/pkg/logger"
	t "github.com/treethought/amnisiac/pkg/types"
)

type MPVController struct {
	client  *mpv.Client
	queue   map[int]t.Item
	process *os.Process
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

func createSocketFile() error {
	fileName := "/tmp/mpvsocket"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		_, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// StartMPV starts mpv in idle mode and specifies the ipc socket
func StartMPV() (*exec.Cmd, error) {
	createSocketFile()
	cmd := exec.Command("mpv", "--idle=once", "--no-terminal", "--input-ipc-server=/tmp/mpvsocket", "--no-video", "--no-config")
	// cmd := exec.Command("mpv", "--idle", "--input-ipc-server=/tmp/mpvsocket")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("MPV started")
	return cmd, nil
}

func (m *MPVController) Initialize() error {
	m.logger = logger.GetLoggerInstance()
	m.log("initializing MPVController")
	// m.log("Starting MPV process")

	// cmd, err := m.StartMPV()

	// if err != nil {
	// 	return err
	// }

	// for {
	// 	if m.process == nil {
	// 		m.log("Waiting for MPV to start")
	// 	} else {
	// 		m.log("MPV ready")
	// 		break

	// 	}
	// }

	// m.log("Started mpv PID: ", m.process.Pid)

	ipcc := mpv.NewIPCClient("/tmp/mpvsocket") // Lowlevel client
	c := mpv.NewClient(ipcc)                   // Highlevel client, can also use RPCClient
	m.client = c

	return nil

}

func (m *MPVController) Shutdown() error {
	m.log("Shutting down MPV")
	// m.log("removing socket file")
	// os.Remove("/tmp/mpvsocket")
	// m.log("Killing process ", m.process.Pid)
	// err := m.process.Signal(os.Kill)
	// if err != nil {
	// 	m.log("Failed to kill mpv process", m.process.Pid)
	// 	panic(err)
	// }
	// m.log("Killed process", m.process.Pid)
	return nil

}

func (m *MPVController) PlayTrack(item *t.Item) error {
	m.log("Playing track", item.RawTitle)

	err := m.client.Loadfile(item.URL, mpv.LoadFileModeReplace)
	if err != nil {
		return err
	}
	m.status.CurrentItem = item
	m.status.PlayState = "playing"

	return nil

}

func (m *MPVController) TogglePause() error {
	m.log("Pausing current track", m.status.CurrentItem)
	paused, err := m.client.Pause()
	if err != nil {
		return err
	}
	if paused {
		err := m.client.SetPause(false)
		m.status.PlayState = "playing"
		return err

	} else {
		err := m.client.SetPause(true)
		m.status.PlayState = "paused"
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
		return m.status
	}

	paused, err := m.client.Pause()
	if err != nil {
		m.log("Failed to get playing state!!")
		m.log("Player error: ", err.Error())
	}
	if paused {
		m.status.PlayState = "paused"

	} else {
		m.status.PlayState = "playing"
	}

	pos, err := m.GetPosition()
	if err == nil {
		m.status.CurrentPosition = pos
	}

	l, err := m.client.Duration()
	if err == nil {
		m.status.CurrentDuration = l
	}

	return m.status

}

func (m *MPVController) GetPosition() (float64, error) {
	return m.client.Position()
}

func (m *MPVController) Seek(pos int) error {
	return m.client.Seek(pos, mpv.SeekModeAbsolute)
}
