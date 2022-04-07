package local

import (
	. "github.com/elamre/gomberman/netcode/core"
	"time"
)

type LocalClient struct {
	n *FakeNetwork
}

func NewLocalClient(n *FakeNetwork) LocalClient {
	return LocalClient{n: n}
}

func (l *LocalClient) Connect() any {
	return nil
}

func (l *LocalClient) Disconnect() any {
	return nil
}

func (l *LocalClient) Write(packet Packet) any {
	p := NewRawPacket(packet)
	p.PacketTime = time.Now().UnixMilli()
	l.n.ClientWrite(p)
	return nil
}

func (l *LocalClient) SetPacketReceivedCallback(PacketReceived func(packet Packet)) {

}

func (l *LocalClient) GetPacket() *Packet {
	if l.n.serverOutgoing.Length() == 0 {
		return nil
	}
	dat := l.n.serverOutgoing.Pop()
	pack := BytesToPacket(dat.Data)
	return &pack
}

func (l *LocalClient) WaitForPacket() Packet {
	r := l.n.ClientRead()
	return r.ContainingPacket
}

func (l *LocalClient) Close() any {
	return nil
}
