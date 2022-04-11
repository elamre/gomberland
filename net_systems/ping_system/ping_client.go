package ping_system

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/ping_system/ping_packets"
	"time"
)

type PingClientSystem struct {
	client        net.Client
	LastPing      int64
	lastBroadcast time.Time
	Interval      time.Duration
}

func NewPingClientSystem(client net.Client) *PingClientSystem {
	return &PingClientSystem{client: client, lastBroadcast: time.Now(), Interval: time.Second}
}

func (p *PingClientSystem) pingCallback(c net.Client, d common_system.ClientRegulator, pack packet_interface.Packet) {
	ping := pack.(ping_packets.PingPacket)
	p.LastPing = ping.GetPing()
}

func (p *PingClientSystem) RegisterCallbacks(r common_system.ClientRegulator) {
	r.RegisterPacketCallback(p.pingCallback, ping_packets.PingPacket{})
}

func (p *PingClientSystem) GetPing() time.Duration {
	return time.Duration(p.LastPing) * time.Millisecond
}

func (p *PingClientSystem) Update() {
	if time.Since(p.lastBroadcast) >= p.Interval {
		p.lastBroadcast = time.Now()
		p.client.WritePacket(ping_packets.PingPacket{CreationTime: time.Now()})
	}
}
