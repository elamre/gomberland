package server

import (
	"github.com/elamre/serverclient/netcode/core"
	"github.com/pion/webrtc/v3"
)

type ClientConnection struct {
	peerConnection *webrtc.PeerConnection
	dataChannel    *webrtc.DataChannel
}

func (c *ClientConnection) WritePacket(packet core.Packet) {

}
