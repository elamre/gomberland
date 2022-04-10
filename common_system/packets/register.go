package packets

import "github.com/elamre/gomberman/net/packet_interface"

func Register() {
	packet_interface.RegisterPacket(ChatPacket{})
	packet_interface.RegisterPacket(ConnectionPacket{})
	packet_interface.RegisterPacket(PingPacket{})
}
