package game_system

import (
	"github.com/elamre/gomberman/net"
	common_system2 "github.com/elamre/gomberman/net_systems/common_system"
	"log"
	"time"
)

type GameServerSystemOptions struct {
	TicksPerSecond int
}

type GameServerSystem struct {
	players    []*common_system2.NetPlayer
	server     net.Server
	options    GameServerSystemOptions
	nextUpdate time.Time

	tickAmount time.Duration
}

func NewGameServerSystem(players []*common_system2.NetPlayer, server net.Server, options GameServerSystemOptions) *GameServerSystem {
	g := &GameServerSystem{players: players, server: server, options: options}
	fraction := time.Second.Microseconds() / (time.Duration(options.TicksPerSecond) * time.Second).Milliseconds()
	g.tickAmount = time.Duration(fraction) * time.Millisecond
	return g
}

func (p *GameServerSystem) RegisterCallbacks(r common_system2.ServerRegulator) {
}

func (p *GameServerSystem) Update() {
	if time.Since(p.nextUpdate) >= p.tickAmount {
		p.nextUpdate = time.Now()
		log.Printf("update: %+v", p.nextUpdate)
	}
}
