package main

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"github.com/elamre/gomberman/net/webrtc"
	"github.com/elamre/gomberman/net_systems/common_system/common_packets"
	"github.com/elamre/gomberman/net_systems/lobby_system/lobby_system_packets"
	"github.com/elamre/gomberman/net_systems/ping_system/ping_packets"
	"github.com/elamre/logger/pkg/logger"
	"sync"
)

const localIp = "127.0.0.1"
const hetznerIp = "78.47.36.203"

const port = 50001
const serverIp = localIp

const NetWebRtcInterface = "webrtc"
const NetLocalInterface = "local"
const NetTcpInterface = "tcp"

const netInterface = NetWebRtcInterface

var once sync.Once
var settingsLogger = logger.GetSettings().GetLogger("settings")

func init() {
	once.Do(func() {
		// Just to make sure, in came init gets called accidentally
		common_packets.Register()
		lobby_system_packets.Register()
		ping_packets.Register()
		settingsLogger.LogInfof("Packages registered: %+v", packet_interface.GetRegisteredPackets())
	})
}

func GetNetObjects() (net.Server, net.Client) {
	settingsLogger.LogInfof("creating interfaces at: %s:%d", serverIp, port)
	var serverObject net.Server
	var clientObject net.Client
	switch netInterface {
	case NetWebRtcInterface:
		serverObject = webrtc.NewWebrtcHost(serverIp, port)
		clientObject = webrtc.NewWebrtcClient(serverIp, port)
		settingsLogger.LogInfo("Created webrtc interface")
	case NetLocalInterface:
		settingsLogger.LogError("Local interface not supported")
	case NetTcpInterface:
		settingsLogger.LogError("Tcp interface not supported")
	}
	return serverObject, clientObject
}
