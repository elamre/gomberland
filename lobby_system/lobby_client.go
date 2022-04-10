package lobby_system

import (
	"github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/common_system/common_packets"
	packets2 "github.com/elamre/gomberman/lobby_system/lobby_system_packets"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"log"
)

type LobbyClientSystem struct {
	curPlayer *common_system.NetPlayer
	client    net.Client
}

func NewLobbyClientSystem(client net.Client) *LobbyClientSystem {
	c := &LobbyClientSystem{client: client, curPlayer: &common_system.NetPlayer{}}
	return c
}

func (c *LobbyClientSystem) SendPacket(pack packet_interface.Packet) any {
	return c.client.WritePacket(pack)
}

func (c *LobbyClientSystem) RegisterPlayer(name string) {
	if !c.curPlayer.HasRegistered {
		c.client.WritePacket(common_packets.ConnectionPacket{
			Action:  0,
			Message: name,
		})
		c.curPlayer.HasRegistered = true
	}
}

func (c *LobbyClientSystem) Update() {
	if !c.client.IsConnected() {
		return
	}

	pack, err := c.client.ReadPacket()
	if err != nil {
		panic(err)
	}
	switch t := pack.(type) {
	case common_packets.ConnectionPacket:
		if t.Action == common_packets.ConnectionRefusedAction {
			c.curPlayer.HasRegistered = false
			log.Printf("Unable to register: %s", t.Message)
		} else if t.Action == common_packets.ConnectionAcceptedAction {
			c.curPlayer.Id = t.UserId
			log.Println("accepted")
		} else {
			log.Printf("Player with name: %s registered (%d)", t.Message, t.UserId)
		}
	case packets2.RoomUpdatePacket:
		for _, room := range t.Rooms {
			log.Printf("Room: %s", room.RoomName)
			players := ""
			for _, p := range room.Players {
				players += p.Name + ", "
			}
			log.Println(players)
		}
	default:
		log.Printf("type; %T, contains: %+v", t, t)
	}
}
