package main

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net_systems"
	"github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/game_system"
	"github.com/elamre/gomberman/net_systems/lobby_system"
	"github.com/elamre/gomberman/net_systems/ping_system"
	"github.com/elamre/logger/pkg/logger"
	"time"
)

const (
	serverLobbyTag = "serverLobby"
	serverUsersTag = "serverUsers"
	serverPingTag  = "serverPing"
	serverGameTag  = "serverGame"
)

const (
	clientLobbyTag = "clientLobby"
	clientPingTag  = "clientPing"
)

func CreateServerSystem(server net.Server) *net_systems.ServerDelegator {
	userManagement := common_system.NewUserManagement(server)
	lobby := lobby_system.NewLobbyServerSystem(userManagement)
	serverDelegator := net_systems.NewServerDelegator(server)

	serverDelegator.RegisterSubSystem(serverUsersTag, userManagement)
	serverDelegator.RegisterSubSystem(serverLobbyTag, lobby)
	serverDelegator.RegisterSubSystem(serverPingTag, ping_system.NewPingServerSystem())
	serverDelegator.RegisterSubSystem(serverGameTag, game_system.NewGameServerSystem(nil, server, game_system.GameServerSystemOptions{TicksPerSecond: 1}))
	return serverDelegator
}

func CreateClientSystem(client net.Client) *net_systems.ClientDelegator {
	clientDelegator := net_systems.NewClientDelegator(client)
	clientLobby := lobby_system.NewLobbyClientSystem(client)
	clientDelegator.RegisterSubSystem(clientLobbyTag, clientLobby)
	clientDelegator.RegisterSubSystem(clientPingTag, ping_system.NewPingClientSystem(client))
	return clientDelegator
}

func NetTest() {
	mainLogger := logger.GetSettings().GetLogger("main")
	serverObject, clientObject := GetNetObjects()

	go func() {
		serverDelegator := CreateServerSystem(serverObject)
		serverObject.Start()

		for {
			serverDelegator.Update()
		}
	}()

	clientDelegator := CreateClientSystem(clientObject)
	clientObject.Connect()
	for !clientObject.IsConnected() {
	}
	lobby := clientDelegator.GetSubsystem(clientLobbyTag).(*lobby_system.LobbyClientSystem)
	lobby.RegisterPlayer("Elmar")
	state := "wait"
	start := time.Now()
	lobby.OnRegisteredAction = func(success bool, reason string, player *common_system.NetPlayer) {
		if success {
			state = "registered"
			start = time.Now()
			mainLogger.LogInfo("We are registered")
		} else {
			mainLogger.LogWarningf("Could not join: %s", reason)
		}
	}

	go func() {
		for {
			clientDelegator.Update()
			if time.Since(start) > 5*time.Second {
				start = time.Now()
				switch state {
				case "registered":
					state = "roomcreate"
					lobby.CreateRoom("naam", "")
				case "roomcreate":
					state = "roomready"
					lobby.SetReady()
				case "roomready":
					state = "nothing"
					lobby.StartRoom()
				}
			}
		}
	}()
}
