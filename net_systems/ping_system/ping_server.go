package ping_system

import (
	"github.com/elamre/gomberman/net/packet_interface"
	common_system2 "github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/ping_system/ping_packets"
	"time"
)

type PingServerSystem struct {
	pastBroadcast time.Time
}

func NewPingServerSystem() *PingServerSystem {
	return &PingServerSystem{}
}

func (p *PingServerSystem) pingCallback(c *common_system2.ServerPlayer, d common_system2.ServerRegulator, pack packet_interface.Packet) {
	c.WritePacket(pack)
}

func (p *PingServerSystem) RegisterCallbacks(r common_system2.ServerRegulator) {
	r.RegisterPacketCallback(p.pingCallback, ping_packets.PingPacket{})
}

func (p *PingServerSystem) Update() {
	if p.pastBroadcast.Sub(time.Now()) > 10*time.Second {
		// Todo broadcsat the ping to all
	}
}
