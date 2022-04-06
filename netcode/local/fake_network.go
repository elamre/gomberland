package local

import (
	"github.com/elamre/queue/pkg/queue"
	"github.com/elamre/serverclient/netcode/core"
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

func (f *FakeNetwork) ServerRead() core.Packet {
	dat := f.serverIncoming.Pop()
	return core.BytesToPacket(dat.Data)
}

func (f *FakeNetwork) ServerWrite(packet core.Packet) {
	f.serverOutgoing.Append(&FakeNetworkPacket{core.PacketToBytes(packet)})
}

func (f *FakeNetwork) ClientRead() core.Packet {
	dat := f.serverOutgoing.Pop()
	return core.BytesToPacket(dat.Data)
}

func (f *FakeNetwork) ClientWrite(packet core.Packet) {
	packetBytes := core.PacketToBytes(packet)
	f.serverIncoming.Append(&FakeNetworkPacket{Data: packetBytes})
}
