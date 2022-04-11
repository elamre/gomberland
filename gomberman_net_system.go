package main

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net_systems"
	"github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/game_system"
	"github.com/elamre/gomberman/net_systems/lobby_system"
	"github.com/elamre/gomberman/net_systems/ping_system"
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
