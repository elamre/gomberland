package local

import . "github.com/elamre/serverclient/netcode/core"

type LocalServer struct {
	n *FakeNetwork
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
	l.n.ServerWrite(packet)
	return nil
}

func (l *LocalServer) SetPacketReceivedCallback(PacketReceived func(packet Packet)) {

}

func (l *LocalServer) GetPacket() *Packet {
	/*if l.n.serverIncoming.Length() == 0 {
		return nil
	}*/
	dat := l.n.serverIncoming.Pop()
	pack := BytesToPacket(dat.Data)
	return &pack
}

func (l *LocalServer) Close() any {
	return nil
}
