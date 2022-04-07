package core

type ClientConnectionInterface interface {
	Connect() any
	Disconnect() any
	Write(packet Packet) any
	SetPacketReceivedCallback(PacketReceived func(packet Packet))
	GetPacket() *Packet
	Close() any
}

type ServerConnectionInterface interface {
	Start(port int)
	Stop()
	BroadcastPacket(packet Packet)
	SetClientConnectedCallback(Connected func(connectionInterface ClientConnectionInterface))
}
