package common_system

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
)

type PacketCallback = func(c net.ServerClient, d ServerRegulator, pack packet_interface.Packet)

type SubSystem interface {
	RegisterCallbacks(r ServerRegulator)
	Update()
}

type ServerRegulator interface {
	RegisterSubSystem(name string, system SubSystem)
	RegisterPacketCallback(cb PacketCallback, packet packet_interface.Packet)
	RemovePacketCallback(cb PacketCallback, packetType packet_interface.Packet)
	RemoveSubSystem(name string)
	GetSubsystem(name string) interface{}
}
