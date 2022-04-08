package webrtc

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"github.com/elamre/gomberman/net/webrtc/internal/webrtc_server"
	"sync"
)

type WebrtcHost struct {
	clientsMutex        sync.Mutex
	server              *webrtc_server.Server
	connectedClientsMap map[*webrtc_server.Connection]net.ServerClient
	connectedClients    []net.ServerClient
	onConnectionCb      func(client net.ServerClient)
	onDisconnectCb      func(client net.ServerClient)
}

func NewWebrtcHost(publicIp string, port int) *WebrtcHost {
	w := WebrtcHost{
		connectedClientsMap: make(map[*webrtc_server.Connection]net.ServerClient),
		connectedClients:    make([]net.ServerClient, 0),
	}
	w.server = webrtc_server.New(webrtc_server.Options{
		MaxConnections: 10,
		HttpPort:       port,
		PublicIP:       publicIp,
		ICEServerURLs:  []string{"stun:127.0.0.1:3478"},
	})
	w.server.OnConnection = w.onConnect
	w.server.OnDisconnect = w.onDisconnect
	return &w
}

func (w *WebrtcHost) ClientIterator(iterator func(c net.ServerClient)) {
	for _, c := range w.connectedClients {
		iterator(c)
	}
}

func (w *WebrtcHost) Start() any {
	w.server.Start()
	return nil
}

func (w *WebrtcHost) Close() any {
	// todo
	return nil
}

func (w *WebrtcHost) onDisconnect(connection *webrtc_server.Connection) {
	w.clientsMutex.Lock()
	defer w.clientsMutex.Unlock()
	if w.onDisconnectCb != nil {
		w.onDisconnectCb(w.connectedClientsMap[connection])
	}
	delete(w.connectedClientsMap, connection)
}

func (w *WebrtcHost) onConnect(connection *webrtc_server.Connection) {
	w.clientsMutex.Lock()
	defer w.clientsMutex.Unlock()

	wrappedHostclient := NewWebrtcHostClient(connection)
	w.connectedClients = append(w.connectedClients, wrappedHostclient)
	w.connectedClientsMap[connection] = wrappedHostclient

	if w.onConnectionCb != nil {
		w.onConnectionCb(wrappedHostclient)

	}
}

func (w *WebrtcHost) SetOnConnection(callb func(client net.ServerClient)) {
	w.onConnectionCb = callb

}

func (w *WebrtcHost) SetOnDisconnection(callb func(client net.ServerClient)) {
	w.onDisconnectCb = callb
}

func (w *WebrtcHost) BroadcastPacket(packet packet_interface.Packet) {

}
