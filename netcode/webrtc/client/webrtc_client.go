package client

import (
	"github.com/elamre/serverclient/netcode/core"
	"github.com/pion/webrtc/v3"
)

type WebrtcClient struct {
	IpAddress string

	packets        chan []byte
	peerConnection *webrtc.PeerConnection
	dataChannel    *webrtc.DataChannel
}

func (c *WebrtcClient) Disconnect() any {
	//TODO implement me
	panic("implement me")
}

func (c *WebrtcClient) Write(packet core.Packet) any {

	//c.dataChannel.Send(data)
	//TODO implement me
	panic("implement me")
}

func (c *WebrtcClient) SetPacketReceivedCallback(PacketReceived func(packet core.Packet)) {
	//TODO implement me
	panic("implement me")
}

func (c *WebrtcClient) Write(data []byte) error {
	return
}

func (c *WebrtcClient) Received(data []byte) {
	//TODO implement me
	panic("implement me")
}

func NewClient(ipAddress string) core.ConnectionInterface {
	c := &WebrtcClient{IpAddress: ipAddress}
	return c
}

func (c *WebrtcClient) onMessage(msg webrtc.DataChannelMessage) {

}

func (c *WebrtcClient) Connect() any {
	c.dataChannel.OnMessage(c.onMessage)
	return nil
}
