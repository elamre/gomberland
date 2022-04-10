package local

import (
	"github.com/elamre/gomberman/net/webrtc"
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

func (f *FakeNetwork) ServerRead() *webrtc.RawPacket {
	dat := f.serverIncoming.Pop()
	return webrtc.RawPacketFrom(dat.Data)
}

func (f *FakeNetwork) ServerWrite(packet *webrtc.RawPacket) {
	f.serverOutgoing.Append(&FakeNetworkPacket{packet.GetBytes()})
}

func (f *FakeNetwork) ClientRead() *webrtc.RawPacket {
	dat := f.serverOutgoing.Pop()
	return webrtc.RawPacketFrom(dat.Data)
}

func (f *FakeNetwork) ClientWrite(packet *webrtc.RawPacket) {
	f.serverIncoming.Append(&FakeNetworkPacket{packet.GetBytes()})
}
