package common_system

import (
	"fmt"
	"github.com/elamre/gomberman/net"
)

type ServerPlayer struct {
	NetPlayer *NetPlayer
	Client    net.ServerClient
}

func (s ServerPlayer) String() string {
	if s.NetPlayer == nil {
		return fmt.Sprintf("%+v nil", s.Client)
	}
	return fmt.Sprintf("%+v, %s", s.NetPlayer, s.NetPlayer.String())
}

type NetPlayer struct {
	HasRegistered bool
	Ready         bool
	Name          string
	Id            uint32
	RoomId        string
}

func (n NetPlayer) String() string {
	retString := fmt.Sprintf("%s[%d]", n.Name, n.Id)
	if !n.HasRegistered {
		retString += " [not registered]"
	}
	if n.Ready {
		retString += " [ready]"
	}
	if n.RoomId != "" {
		retString += fmt.Sprintf(" [room: %s]", n.RoomId)
	}
	return retString
}
