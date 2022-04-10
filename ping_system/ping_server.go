package ping_system

import (
	"github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"github.com/elamre/gomberman/ping_system/ping_packets"
	"time"
)

type PingServerSystem struct {
	pastBroadcast time.Time
}

func NewPingServerSystem() *PingServerSystem {
	return &PingServerSystem{}
}

func (p *PingServerSystem) pingCallback(c net.ServerClient, d common_system.ServerRegulator, pack packet_interface.Packet) {
	c.WritePacket(pack)
}

func (p *PingServerSystem) RegisterCallbacks(r common_system.ServerRegulator) {
	r.RegisterPacketCallback(p.pingCallback, ping_packets.PingPacket{})
}

func (p *PingServerSystem) Update() {
	if p.pastBroadcast.Sub(time.Now()) > 10*time.Second {
		// Todo broadcsat the ping to all
	}
}
