package core

import (
	"github.com/elamre/gomberman/netcode/packet"
	"log"
)

var (
	PingPacketIndex        uint32
	ChatPacketIndex        uint32
	EntityStatePacketIndex uint32
	PhysicsPacketIndex     uint32
	ConnectionPacketIndex  uint32
)

func RegisterPackets() {
	// Registration here, order is important

	PingPacketIndex = RegisterPacket(packet.PingPacket{})
	ChatPacketIndex = RegisterPacket(packet.ChatPacket{})
	EntityStatePacketIndex = RegisterPacket(packet.EntityStatePacket{})
	PhysicsPacketIndex = RegisterPacket(packet.PhysicsPacket{})
	PhysicsPacketIndex = RegisterPacket(packet.ConnectionPacket{})
	log.Printf("Registered packets: %+v\n", GetRegisteredPackets())
}
