package common_system

import (
	"github.com/elamre/gomberman/common_system/packets"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"log"
	"strings"
	"sync"
)

type ServerPlayer struct {
	NetPlayer *NetPlayer
	Client    net.ServerClient
}

type UserManagement struct {
	playersAdjustMutex sync.Mutex
	players            []*ServerPlayer
	ClientToPlayer     map[net.ServerClient]*ServerPlayer
	PlayerToClient     map[*ServerPlayer]net.ServerClient
	NameToClient       map[string]*ServerPlayer
	IdToClient         map[uint32]*ServerPlayer
	clientIdx          uint32
}

func NewUserManagement() *UserManagement {
	return &UserManagement{
		players:        make([]*ServerPlayer, 0),
		NameToClient:   make(map[string]*ServerPlayer),
		IdToClient:     make(map[uint32]*ServerPlayer),
		ClientToPlayer: make(map[net.ServerClient]*ServerPlayer),
		PlayerToClient: make(map[*ServerPlayer]net.ServerClient),
	}
}

func (u *UserManagement) BroadcastFilter(pack packet_interface.Packet, filter func(player *ServerPlayer) bool) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	for _, p := range u.players {
		if filter(p) {
			p.Client.WritePacket(pack)
		}
	}
}

func (u *UserManagement) OnConnection(client net.ServerClient) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	netPlayer := &ServerPlayer{Client: client, NetPlayer: &NetPlayer{}}
	u.players = append(u.players, netPlayer)
	u.PlayerToClient[netPlayer] = client
	u.ClientToPlayer[client] = netPlayer
}

func (u *UserManagement) OnDisconnection(client net.ServerClient) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	// TODO check this out
	for i := 0; i < len(u.players); i++ {
		if u.players[i] == u.ClientToPlayer[client] {
			if i == len(u.players)-1 {
				u.players = u.players[:i-1]
			} else {
				u.players = append(u.players[:i], u.players[i+1:]...)

			}
			break
		}
	}
	serverClient := u.ClientToPlayer[client]

	if serverClient.NetPlayer.HasRegistered {
		delete(u.NameToClient, serverClient.NetPlayer.Name)
		delete(u.IdToClient, serverClient.NetPlayer.Id)
	}
	delete(u.ClientToPlayer, client)
	delete(u.PlayerToClient, serverClient)
}

func (u *UserManagement) HandleConnectionPacket(connection net.ServerClient, pack packets.ConnectionPacket) {
	log.Printf("New request: %+v", pack)
	player := u.ClientToPlayer[connection]
	if player.NetPlayer.HasRegistered {
		connection.WritePacket(packets.ConnectionPacket{
			UserId:  0,
			Action:  packets.ConnectionRefusedAction,
			Message: "Already registered",
		})
		return
	}
	if len(pack.Message) == 0 {
		connection.WritePacket(packets.ConnectionPacket{
			UserId:  0,
			Action:  packets.ConnectionRefusedAction,
			Message: "No name given",
		})
		return
	}
	name := strings.TrimSpace(strings.ToLower(pack.Message))
	if _, ok := u.NameToClient[name]; ok {
		connection.WritePacket(packets.ConnectionPacket{
			UserId:  0,
			Action:  packets.ConnectionRefusedAction,
			Message: "Name already registered",
		})
		return
	}
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	u.IdToClient[u.clientIdx] = player
	u.NameToClient[name] = player

	pack.Action = packets.ConnectionAcceptedAction
	pack.Message = name
	pack.UserId = u.clientIdx
	connection.WritePacket(pack)
	u.clientIdx++
}
