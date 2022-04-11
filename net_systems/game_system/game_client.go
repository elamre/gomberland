package game_system

import "github.com/elamre/gomberman/net"

type GameClientSystem struct {
	client net.Client
}

func NewGameClientSystem(client net.Client) *GameClientSystem {
	return &GameClientSystem{client: client}
}

func (g *GameClientSystem) Update() {

}
