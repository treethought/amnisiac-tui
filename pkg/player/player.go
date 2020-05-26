package player

import (
	t "github.com/treethought/amnisiac/pkg/types"
)

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
	CurrentItem          *t.Item
	CurrentPosition      float64
	CurrentDuration      float64
	CurrentQueuePosition int
	CurrentQueue         map[int]t.Item
	PlayState            string
}
