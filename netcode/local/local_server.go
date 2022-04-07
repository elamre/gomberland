package local

import (
	. "github.com/elamre/gomberman/netcode/core"
	"time"
)

type LocalServer struct {
	packetIndex uint64
	n           *FakeNetwork
}

func NewLocalServer(n *FakeNetwork) LocalServer {
	return LocalServer{n: n}
}

func (l *LocalServer) Connect() any {
	return nil
}

func (l *LocalServer) Disconnect() any {
	return nil
}

func (l *LocalServer) Write(packet Packet) any {
	r := NewRawPacket(packet)
	r.PacketTime = time.Now().UnixMilli()
	r.PacketId = l.packetIndex
	r.PacketId++
	l.n.ServerWrite(r)
	return nil
}

func (l *LocalServer) SetPacketReceivedCallback(PacketReceived func(packet Packet)) {

}

func (l *LocalServer) GetPacket() Packet {
	pack := l.n.ServerRead()
	// Check index
	return pack.ContainingPacket
}

func (l *LocalServer) Close() any {
	return nil
}
