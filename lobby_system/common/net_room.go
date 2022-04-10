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
	Owner      uint32
	Players    []*common_system.NetPlayer
	password   string
}

func (n *NetRoom) IsReady() bool {
	for _, p := range n.Players {
		if !p.Ready {
			return false
		}
	}
	return true
}
