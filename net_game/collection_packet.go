package net_game

import (
	"github.com/elamre/gomberman/net/packet_interface"
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

	PingPacketIndex = packet_interface.RegisterPacket(packet.PingPacket{})
	ChatPacketIndex = packet_interface.RegisterPacket(packet.ChatPacket{})
	EntityStatePacketIndex = packet_interface.RegisterPacket(packet.EntityStatePacket{})
	PhysicsPacketIndex = packet_interface.RegisterPacket(packet.PhysicsPacket{})
	PhysicsPacketIndex = packet_interface.RegisterPacket(packet.ConnectionPacket{})
	log.Printf("Registered packets: %+v\n", packet_interface.GetRegisteredPackets())
}
