package common_system

import (
	"github.com/elamre/go_helpers/pkg/slice_helpers"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"github.com/elamre/gomberman/net_systems/common_system/common_packets"
	"log"
	"strings"
	"sync"
)

// These 2 functions are for ease

func (s *ServerPlayer) WritePacket(pack packet_interface.Packet) any {
	return s.Client.WritePacket(pack)
}

func (s *ServerPlayer) ReadPacket() (packet packet_interface.Packet, err error) {
	return s.Client.ReadPacket()
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

func NewUserManagement(server net.Server) *UserManagement {
	s := &UserManagement{
		players:        make([]*ServerPlayer, 0),
		NameToClient:   make(map[string]*ServerPlayer),
		IdToClient:     make(map[uint32]*ServerPlayer),
		ClientToPlayer: make(map[net.ServerClient]*ServerPlayer),
		PlayerToClient: make(map[*ServerPlayer]net.ServerClient),
	}
	server.AddConnectionCallback(s.OnConnection)
	server.AddDisconnectionCallback(s.OnDisconnection)
	return s
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
}

func (u *UserManagement) HandleConnectionPacket(c *ServerPlayer, d ServerRegulator, pack packet_interface.Packet) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()

	conPacket := pack.(common_packets.ConnectionPacket)

	log.Printf("connection packet: %+v", conPacket)

	if c.NetPlayer.HasRegistered {
		c.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "Already registered",
		})
		return
	}
	if len(conPacket.Message) == 0 {
		c.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "No name given",
		})
		return
	}
	name := strings.TrimSpace(strings.ToLower(conPacket.Message))
	if _, ok := u.NameToClient[name]; ok {
		c.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "Name already registered",
		})
		return
	}
	c.NetPlayer.HasRegistered = true
	c.NetPlayer.Name = name
	c.NetPlayer.Id = u.clientIdx

	u.IdToClient[u.clientIdx] = c
	u.NameToClient[name] = c

	conPacket.Action = common_packets.ConnectionAcceptedAction
	conPacket.Message = name
	conPacket.UserId = u.clientIdx
	c.WritePacket(conPacket)
	u.clientIdx++
}

func (s *UserManagement) RegisterCallbacks(r ServerRegulator) {
	r.RegisterPacketCallback(s.HandleConnectionPacket, common_packets.ConnectionPacket{})
}
func (s *UserManagement) Update() {

}
