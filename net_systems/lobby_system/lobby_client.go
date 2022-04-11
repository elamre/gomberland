package lobby_system

import (
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	common_system2 "github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/common_system/common_packets"
	packets2 "github.com/elamre/gomberman/net_systems/lobby_system/lobby_system_packets"
	"log"
)

type LobbyClientSystem struct {
	OnRegisteredAction func()
	curPlayer          *common_system2.NetPlayer
	client             net.Client
}

func NewLobbyClientSystem(client net.Client) *LobbyClientSystem {
	c := &LobbyClientSystem{client: client, curPlayer: &common_system2.NetPlayer{}}
	return c
}

func (c *LobbyClientSystem) SendRoomPacket(pack packets2.RoomPacket) any {
	pack.UserId = c.curPlayer.Id
	return c.client.WritePacket(pack)

}

func (c *LobbyClientSystem) SendPacket(pack packet_interface.Packet) any {
	return c.client.WritePacket(pack)
}

func (c *LobbyClientSystem) RegisterPlayer(name string) {
	if !c.curPlayer.HasRegistered {
		c.client.WritePacket(common_packets.ConnectionPacket{
			Action:  common_packets.ConnectionRegisterAction,
			Message: name,
		})
		c.curPlayer.HasRegistered = true
	}
}

func (p *LobbyClientSystem) connectionCallback(c net.Client, d common_system2.ClientRegulator, pack packet_interface.Packet) {
	t := pack.(common_packets.ConnectionPacket)
	log.Printf("Client received: %s", t.String())
	if t.Action == common_packets.ConnectionRefusedAction {
		p.curPlayer.HasRegistered = false
		log.Printf("Unable to register: %s", t.Message)
	} else if t.Action == common_packets.ConnectionAcceptedAction {
		p.curPlayer.Id = t.UserId
		if p.OnRegisteredAction != nil {
			p.OnRegisteredAction()
		}
		log.Println("accepted")
	} else {
		log.Printf("Player with name: %s registered (%d)", t.Message, t.UserId)
	}
}
func (lp *LobbyClientSystem) roomupdateCallback(c net.Client, d common_system2.ClientRegulator, pack packet_interface.Packet) {
	/*	t := pack.(packets2.RoomUpdatePacket)
		for _, room := range t.Rooms {
			log.Printf("Room: %s", room.RoomName)
			players := ""
			for _, p := range room.Players {
				players += p.Name + ", "
			}
			log.Println(players)
		}*/
}

func (lp *LobbyClientSystem) roomPacketCallback(c net.Client, d common_system2.ClientRegulator, pack packet_interface.Packet) {
	log.Printf("Received: %+v", pack)
}

func (p *LobbyClientSystem) RegisterCallbacks(r common_system2.ClientRegulator) {
	r.RegisterPacketCallback(p.connectionCallback, common_packets.ConnectionPacket{})
	r.RegisterPacketCallback(p.roomupdateCallback, packets2.RoomUpdatePacket{})
	r.RegisterPacketCallback(p.roomPacketCallback, packets2.RoomPacket{})
}

func (c *LobbyClientSystem) Update() {
	if !c.client.IsConnected() {
		return
	}
	// We can draw here something
}
