package core

import (
	"github.com/elamre/serverclient/netcode/packet"
)

var (
	PingPacketIndex        int
	ChatPacketIndex        int
	EntityStatePacketIndex int
	PhysicsPacketIndex     int
)

func RegisterPackets() {
	// Registration here, order is important

	PingPacketIndex = RegisterPacket(packet.PingPacket{})
	ChatPacketIndex = RegisterPacket(packet.ChatPacket{})
	EntityStatePacketIndex = RegisterPacket(packet.EntityStatePacket{})
	PhysicsPacketIndex = RegisterPacket(packet.PhysicsPacket{})

}
