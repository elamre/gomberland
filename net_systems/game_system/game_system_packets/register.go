package game_system_packets

import (
	"github.com/elamre/gomberman/net/packet_interface"
)

func Register() {
	packet_interface.RegisterPacket(EntityStatePacket{})
	packet_interface.RegisterPacket(PhysicsPacket{})
}
