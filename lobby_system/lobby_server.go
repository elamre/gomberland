package lobby_system

import (
	"github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/common_system/packets"
	"github.com/elamre/gomberman/lobby_system/common"
	packets2 "github.com/elamre/gomberman/lobby_system/packets"
	"github.com/elamre/gomberman/net"
	"log"
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

func (s *LobbyServerSystem) Update() {
	s.server.ClientIterator(func(c net.ServerClient) {
		pack, err := c.ReadPacket()
		if err != nil {
			log.Printf("Error reading: %v", err)
			return
		}
		if pack != nil {
			switch t := pack.(type) {
			case packets.ChatPacket:
			case packets2.RoomPacket:
				switch t.Action {
				case packets2.RoomStartAction:
					if s.OnRoomStart != nil {
						s.OnRoomStart(s.nameToRoom[t.Name])
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
					newRoom := common.NetRoom{
						RoomName: name,
						Players:  []*common_system.NetPlayer{s.userManage.ClientToPlayer[c].NetPlayer},
					}
					s.rooms = append(s.rooms, &newRoom)
					c.WritePacket(packets2.RoomPacket{
						Action: packets2.RoomCreateSuccessAction,
						Name:   name,
					})
				}
			case packets.ConnectionPacket:
				s.userManage.HandleConnectionPacket(c, t)
			default:
				log.Printf("unhandled type: %T", t)
			}
		}
	})
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
