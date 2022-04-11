package common_system

import "github.com/elamre/gomberman/net"

type ServerPlayer struct {
	NetPlayer *NetPlayer
	Client    net.ServerClient
}

type NetPlayer struct {
	HasRegistered bool
	Ready         bool
	Name          string
	Id            uint32
	RoomId        string
}
