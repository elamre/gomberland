package lobby_system

import (
	"fmt"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	common_system2 "github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/common_system/common_packets"
	"github.com/elamre/gomberman/net_systems/lobby_system/common"
	packets2 "github.com/elamre/gomberman/net_systems/lobby_system/lobby_system_packets"
	"log"
	"strings"
	"time"
)

type LobbyServerSystem struct {
	userManage *common_system2.UserManagement

	playersInLobby map[uint32]*common_system2.ServerPlayer
	playersInRoom  map[uint32]*common.NetRoom

	nameToRoom      map[string]*common.NetRoom
	rooms           []*common.NetRoom
	roomUpdateTimer time.Time

	OnRoomStart func(room *common.NetRoom)
}

func NewLobbyServerSystem(management *common_system2.UserManagement) *LobbyServerSystem {
	s := &LobbyServerSystem{
		userManage: management,
		rooms:      make([]*common.NetRoom, 0),
		nameToRoom: make(map[string]*common.NetRoom),
	}

	s.roomUpdateTimer = time.Now()
	return s
}

func (s *LobbyServerSystem) ClientDisconnect(c common_system2.ServerPlayer) {
	if disconnectedPlayer, ok := s.playersInLobby[c.NetPlayer.Id]; ok {
		delete(s.playersInLobby, disconnectedPlayer.NetPlayer.Id)
		// Write something to all other players
		for _, op := range s.playersInLobby {
			op.WritePacket(common_packets.ChatPacket{Message: fmt.Sprintf("[%s] Disconnected", c.NetPlayer.Name)})
		}
		if insideRoom, ok := s.playersInRoom[disconnectedPlayer.NetPlayer.Id]; ok {
			if insideRoom.Owner == disconnectedPlayer.NetPlayer.Id {
				// Have to delete the room
				for _, op := range s.playersInLobby {
					op.WritePacket(packets2.RoomPacket{Action: packets2.RoomDeleteAction})
				}
			} else {
				for _, op := range s.playersInLobby {
					op.WritePacket(packets2.RoomPacket{Action: packets2.RoomLeaveAction, UserId: disconnectedPlayer.NetPlayer.Id})
				}
			}
			delete(s.playersInLobby, disconnectedPlayer.NetPlayer.Id)
		}
	}
}

func (s *LobbyServerSystem) roomPacketCallback(c *common_system2.ServerPlayer, d common_system2.ServerRegulator, pack packet_interface.Packet) {
	t := pack.(packets2.RoomPacket)
	log.Printf("Server room: %s", t.String())

	if t.Action == packets2.RoomCreateAction {
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
		c.NetPlayer.RoomId = name
		newRoom := common.NetRoom{
			RoomName: name,
			Owner:    c.NetPlayer.Id,
			Players:  []*common_system2.NetPlayer{c.NetPlayer},
		}
		s.nameToRoom[name] = &newRoom
		s.rooms = append(s.rooms, &newRoom)
		c.WritePacket(packets2.RoomPacket{
			Action: packets2.RoomCreateSuccessAction,
			Name:   name,
		})
	} else {
		room := s.nameToRoom[c.NetPlayer.RoomId]
		switch t.Action {
		case packets2.RoomReadyAction:
			// Check for illegal action
			c.NetPlayer.Ready = !c.NetPlayer.Ready
			log.Printf("%+v", s.nameToRoom[c.NetPlayer.RoomId].Players)
			for _, p := range room.Players {
				s.userManage.IdToClient[p.Id].WritePacket(packets2.RoomPacket{
					UserId: room.Owner,
					Action: packets2.RoomStartAction,
					Name:   room.RoomName,
				})
			}
		case packets2.RoomStartAction:
			if room.IsReady() {
				for _, p := range room.Players {
					s.userManage.IdToClient[p.Id].WritePacket(packets2.RoomPacket{
						UserId: room.Owner,
						Action: packets2.RoomStartAction,
						Name:   room.RoomName,
					})
				}
				if s.OnRoomStart != nil {
					s.OnRoomStart(room)
				}
			} else {
				c.WritePacket(common_packets.ChatPacket{Message: "Can't start while not everybody is ready"})
			}
		case packets2.RoomLeaveAction:
			if room.Owner == c.NetPlayer.Id {
				delete(s.nameToRoom, room.RoomName)
				// Delete the room
			}
		}
	}
}

func (s *LobbyServerSystem) RegisterCallbacks(r common_system2.ServerRegulator) {
	r.RegisterPacketCallback(s.roomPacketCallback, packets2.RoomPacket{})
}

func (s *LobbyServerSystem) Update() {
	if time.Since(s.roomUpdateTimer) >= 500*time.Millisecond {
		for _, p := range s.playersInLobby {
			pack := packets2.RoomUpdatePacket{Rooms: s.rooms}
			p.WritePacket(pack)
		}
		s.roomUpdateTimer = time.Now()
	}
}

type ClientWorld struct {
	client net.Client
}
