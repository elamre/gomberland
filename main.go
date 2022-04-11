package main

import (
	webrtc2 "github.com/elamre/gomberman/net/webrtc"
	common_packets2 "github.com/elamre/gomberman/net_systems/common_system/common_packets"
	"github.com/elamre/gomberman/net_systems/lobby_system"
	"github.com/elamre/gomberman/net_systems/lobby_system/lobby_system_packets"
	"github.com/elamre/gomberman/net_systems/ping_system"
	"github.com/elamre/gomberman/net_systems/ping_system/ping_packets"
	"log"
	"time"
)

const port = 50001

func main() {
	common_packets2.Register()
	lobby_system_packets.Register()
	ping_packets.Register()

	go func() {
		server := webrtc2.NewWebrtcHost("127.0.0.1", port)
		/*			userManagement := common_system.NewUserManagement(server)
					lobby := lobby_system2.NewLobbyServerSystem(userManagement)

					//server := webrtc2.NewWebrtcHost("192.168.178.43", port)
					//server := webrtc2.NewWebrtcHost("78.47.36.203", port)
					serverDelegator := net_systems.NewServerDelegator(server)
					serverDelegator.RegisterSubSystem("users", userManagement)
					serverDelegator.RegisterSubSystem("serverlobby", lobby)
					serverDelegator.RegisterSubSystem("ping", ping_system.NewPingServerSystem())
					serverDelegator.RegisterSubSystem("game", game_system.NewGameServerSystem(nil, server, game_system.GameServerSystemOptions{TicksPerSecond: 1}))*/
		serverDelegator := CreateServerSystem(server)
		server.Start()

		for {
			serverDelegator.Update()
		}
	}()

	//client := webrtc2.NewWebrtcClient("192.168.178.43", port)
	//client := webrtc2.NewWebrtcClient("78.47.36.203", port)
	client := webrtc2.NewWebrtcClient("127.0.0.1", port)
	clientDelegator := CreateClientSystem(client)
	client.Connect()
	for !client.IsConnected() {
	}
	lobby := clientDelegator.GetSubsystem(clientLobbyTag).(*lobby_system.LobbyClientSystem)
	lobby.RegisterPlayer("Elmar")
	lobby.OnRegisteredAction = func() {
		log.Printf("We are registered")
	}
	start := time.Now()
	go func() {
		for {
			clientDelegator.Update()
			if time.Since(start) > time.Second*2 {
				ping := clientDelegator.GetSubsystem(clientPingTag).(*ping_system.PingClientSystem)
				log.Printf("Ping: %dus", ping.GetPing().Microseconds())
				start = time.Now()
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
