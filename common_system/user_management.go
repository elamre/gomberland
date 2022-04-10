package common_system

import (
	"github.com/elamre/go_helpers/pkg/slice_helpers"
	"github.com/elamre/gomberman/common_system/common_packets"
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
	log.Printf("OnConnection\n%+v\n%+v\n%+v\n%+v\n%+v", u.players, u.IdToClient, u.NameToClient, u.ClientToPlayer, u.PlayerToClient)
	netPlayer := &ServerPlayer{Client: client, NetPlayer: &NetPlayer{}}
	u.players = append(u.players, netPlayer)
	u.PlayerToClient[netPlayer] = client
	u.ClientToPlayer[client] = netPlayer
}

func (u *UserManagement) OnDisconnection(client net.ServerClient) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	tempList := slice_helpers.RemoveFromList[*ServerPlayer](u.ClientToPlayer[client], u.players)
	if tempList == nil {
		log.Println("client not found")
		return
	}
	u.players = tempList

	serverClient := u.ClientToPlayer[client]

	if serverClient.NetPlayer.HasRegistered {
		serverClient.NetPlayer.HasRegistered = false
		delete(u.NameToClient, serverClient.NetPlayer.Name)
		delete(u.IdToClient, serverClient.NetPlayer.Id)
	} else {
		log.Println("we were never registered")
	}
	delete(u.ClientToPlayer, client)
	delete(u.PlayerToClient, serverClient)
	log.Printf("Disconnected, current players: %+v", u.ClientToPlayer)
}

func (u *UserManagement) HandleConnectionPacket(connection net.ServerClient, pack common_packets.ConnectionPacket) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	log.Printf("HandleConnection\n%+v\n%+v\n%+v\n%+v\n%+v", u.players, u.IdToClient, u.NameToClient, u.ClientToPlayer, u.PlayerToClient)
	log.Printf("New request: %+v", pack)
	player := u.ClientToPlayer[connection]
	if player.NetPlayer.HasRegistered {
		connection.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "Already registered",
		})
		return
	}
	if len(pack.Message) == 0 {
		connection.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "No name given",
		})
		return
	}
	name := strings.TrimSpace(strings.ToLower(pack.Message))
	if _, ok := u.NameToClient[name]; ok {
		connection.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "Name already registered",
		})
		return
	}
	player.NetPlayer.HasRegistered = true
	player.NetPlayer.Name = name
	player.NetPlayer.Id = u.clientIdx

	u.IdToClient[u.clientIdx] = player
	u.NameToClient[name] = player

	pack.Action = common_packets.ConnectionAcceptedAction
	pack.Message = name
	pack.UserId = u.clientIdx
	connection.WritePacket(pack)
	u.clientIdx++
}
