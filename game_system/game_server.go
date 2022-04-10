package game_system

import (
	"github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/net"
)

type GameServerSystem struct {
	players []*common_system.NetPlayer
	server  net.Server
}

func NewGameServerSystem(players []*common_system.NetPlayer, server net.Server) *GameServerSystem {
	return &GameServerSystem{players: players}
}

func (g *GameServerSystem) update() {
	
}
