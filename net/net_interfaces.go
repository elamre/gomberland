package net

import "github.com/elamre/gomberman/net/packet_interface"

type ServerOptions struct {
	Port           int
	MaxConnections int
}

type ClientOptions struct {
	Target string
	Port   int
}

type Server interface {
	Start() any
	Close() any
	SetOnConnection(func(client ServerClient))
	SetOnDisconnection(func(client ServerClient))
	BroadcastPacket(packet packet_interface.Packet)
	ClientIterator(iterator func(c ServerClient))
}

type ServerClient interface {
	WritePacket(packet packet_interface.Packet) any
	GotPacket() bool
	ReadPacket() (packet packet_interface.Packet, err error)
	Close() any
}

type Client interface {
	Connect() any
	IsConnected() bool
	WritePacket(packet packet_interface.Packet) any
	GotPacket() bool
	ReadPacket() (packet *packet_interface.Packet, err error)
	Close() any
}
