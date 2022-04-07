package webrtc

import (
	"fmt"
	"github.com/elamre/gomberman/net/webrtc_client"
	"github.com/elamre/gomberman/netcode/core"
	"time"
)

type WebrtcClient struct {
	client    *webrtc_client.Client
	port      int
	ip        string
	packetIdx uint64
}

func NewWebrtcClient(ip string, port int) WebrtcClient {
	c := WebrtcClient{ip: ip, port: port, packetIdx: 1}
	c.client = webrtc_client.New(webrtc_client.Options{
		IPAddress:     fmt.Sprintf("%s:%d", ip, port),
		ICEServerURLs: []string{"stun:127.0.0.1:3478"},
	})
	return c
}

func (w *WebrtcClient) Connect() any {
	w.client.Start()
	startTime := time.Now()
	for time.Since(startTime) < time.Second*5 || w.client.GetLastError() != nil {
		if w.client.HasConnectedOnce() {
			return nil
		}
		if err := w.client.GetLastError(); err != nil {
			return err
		}
	}

	return fmt.Errorf("timed out connecting to %s:%d", w.ip, w.port)
}
func (w *WebrtcClient) Disconnect() any {
	return w.client.GetLastError()
}
func (w *WebrtcClient) Write(packet core.Packet) any {
	if !w.client.IsConnected() {
		return fmt.Errorf("not connected")
	}
	rawPacket := core.NewRawPacket(packet)
	rawPacket.PacketId = w.packetIdx
	w.packetIdx++
	rawPacket.PacketTime = time.Now().UnixMilli()
	if err := w.client.Send(rawPacket.GetBytes()); err != nil {
		return nil
	}
	return w.client.GetLastError()
}
func (w *WebrtcClient) SetPacketReceivedCallback(PacketReceived func(packet core.Packet)) {

}
func (w *WebrtcClient) GetPacket() *core.Packet {
	if data, success := w.client.Read(); success {
		rawPacket := core.RawPacketFrom(data)
		return &rawPacket.ContainingPacket
	}
	return nil
}
func (w *WebrtcClient) Close() any {
	return nil
}
