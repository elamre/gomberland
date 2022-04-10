package lobby_system

import (
	"github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/common_system/common_packets"
	"github.com/elamre/gomberman/lobby_system/common"
	packets2 "github.com/elamre/gomberman/lobby_system/lobby_system_packets"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"strings"
	"time"
)

type LobbyServerSystem struct {
	server     net.Server
	userManage *common_system.UserManagement

	nameToRoom      map[string]*common.NetRoom
	rooms           []*common.NetRoom
	roomUpdateTimer time.Time

	OnRoomStart func(room *common.NetRoom)
}

func NewLobbyServerSystem(server net.Server) *LobbyServerSystem {
	s := &LobbyServerSystem{
		server:     server,
		userManage: common_system.NewUserManagement(),
		rooms:      make([]*common.NetRoom, 0),
		nameToRoom: make(map[string]*common.NetRoom),
	}
	s.server.SetOnConnection(func(client net.ServerClient) {
		s.userManage.OnConnection(client)
	})
	s.server.SetOnDisconnection(func(client net.ServerClient) {
		serverClient := s.userManage.ClientToPlayer[client]
		s.userManage.OnDisconnection(client)
		if serverClient.NetPlayer.HasRegistered {
			// If we were registered, let the other people know
		}
	})

	s.server.Start()
	s.roomUpdateTimer = time.Now()
	return s
}

func (s *LobbyServerSystem) roomPacketCallback(c net.ServerClient, d common_system.ServerRegulator, pack packet_interface.Packet) {
	t := pack.(packets2.RoomPacket)
	room := s.nameToRoom[t.Name]
	switch t.Action {
	case packets2.RoomReadyAction:
		user := s.nameToRoom[t.Name].Players[t.UserId]
		s.nameToRoom[t.Name].Players[t.UserId].Ready = !user.Ready
	case packets2.RoomStartAction:
		if room.IsReady() {
			if s.OnRoomStart != nil {
				s.OnRoomStart(room)
			}
		} else {
			c.WritePacket(common_packets.ChatPacket{Message: "Can't start while not everybody is ready"})
		}
	case packets2.RoomCreateAction:
		name := strings.TrimSpace(strings.ToLower(t.Name))
		for _, r := range s.rooms {
			if strings.Compare(r.RoomName, name) == 0 {
				c.WritePacket(packets2.RoomPacket{
					Action: packets2.RoomCreateFailedAction,
					Name:   "Room already exists",
				})
				return
			}
		}
		s.userManage.ClientToPlayer[c].NetPlayer.RoomId = name
		newRoom := common.NetRoom{
			RoomName: name,
			Owner:    s.userManage.ClientToPlayer[c].NetPlayer.Id,
			Players:  []*common_system.NetPlayer{s.userManage.ClientToPlayer[c].NetPlayer},
		}
		s.rooms = append(s.rooms, &newRoom)
		c.WritePacket(packets2.RoomPacket{
			Action: packets2.RoomCreateSuccessAction,
			Name:   name,
		})
	case packets2.RoomLeaveAction:
		if room.Owner == s.userManage.ClientToPlayer[c].NetPlayer.Id {
			delete(s.nameToRoom, room.RoomName)
			// Delete the room
		}
	}
}

func (s *LobbyServerSystem) connectionCallback(c net.ServerClient, d common_system.ServerRegulator, pack packet_interface.Packet) {
	t := pack.(common_packets.ConnectionPacket)
	s.userManage.HandleConnectionPacket(c, t)
}

func (s *LobbyServerSystem) RegisterCallbacks(r common_system.ServerRegulator) {
	r.RegisterPacketCallback(s.roomPacketCallback, packets2.RoomPacket{})
	r.RegisterPacketCallback(s.connectionCallback, common_packets.ConnectionPacket{})
}

func (s *LobbyServerSystem) Update() {
	if time.Since(s.roomUpdateTimer) >= 500*time.Millisecond {
		s.roomUpdateTimer = time.Now()
		pack := packets2.RoomUpdatePacket{Rooms: s.rooms}
		s.userManage.BroadcastFilter(pack, func(player *common_system.ServerPlayer) bool {
			return true
		})
	}
}

type ClientWorld struct {
	client net.Client
}
