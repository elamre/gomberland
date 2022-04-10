package ping_packets

import (
	"github.com/elamre/gomberman/net/packet_interface"
)

func Register() {
	packet_interface.RegisterPacket(PingPacket{})
}
