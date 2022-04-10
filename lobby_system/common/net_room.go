package common

import (
	"github.com/elamre/gomberman/common_system"
	"sync"
)

type NetRoomSettings struct {
	HasPassword bool
}

type NetRoom struct {
	playerSync sync.Mutex
	RoomName   string
	Players    []*common_system.NetPlayer
	password   string
}
