package local

import . "github.com/elamre/serverclient/netcode/core"

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
	l.n.ClientWrite(packet)
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

func (l *LocalClient) Close() any {
	return nil
}
