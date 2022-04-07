package local

import (
	"github.com/elamre/gomberman/netcode/core"
	"github.com/elamre/queue/pkg/queue"
)

type FakeNetworkPacket struct {
	Data []byte
}

type FakeNetwork struct {
	serverIncoming *queue.Queue[*FakeNetworkPacket]
	serverOutgoing *queue.Queue[*FakeNetworkPacket]
}

func NewFakeNetwork() *FakeNetwork {
	return &FakeNetwork{serverIncoming: queue.New[*FakeNetworkPacket](), serverOutgoing: queue.New[*FakeNetworkPacket]()}
}

func (f *FakeNetwork) ServerRead() *core.RawPacket {
	dat := f.serverIncoming.Pop()
	return core.RawPacketFrom(dat.Data)
}

func (f *FakeNetwork) ServerWrite(packet *core.RawPacket) {
	f.serverOutgoing.Append(&FakeNetworkPacket{packet.GetBytes()})
}

func (f *FakeNetwork) ClientRead() *core.RawPacket {
	dat := f.serverOutgoing.Pop()
	return core.RawPacketFrom(dat.Data)
}

func (f *FakeNetwork) ClientWrite(packet *core.RawPacket) {
	f.serverIncoming.Append(&FakeNetworkPacket{packet.GetBytes()})
}
