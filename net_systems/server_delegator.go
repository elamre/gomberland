package net_systems

import (
	"github.com/elamre/go_helpers/pkg/misc"
	"github.com/elamre/go_helpers/pkg/slice_helpers"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	. "github.com/elamre/gomberman/net_systems/common_system"
	"log"
	"reflect"
	"sync"
)

type ServerDelegator struct {
	clientMutex sync.Mutex
	server      net.Server

	netPlayers []*ServerPlayer

	subSystems         map[string]ServerSubSystem
	subSystemsCallback map[reflect.Type][]PacketServerCallback
}

func NewServerDelegator(server net.Server) *ServerDelegator {
	s := &ServerDelegator{
		subSystems:         make(map[string]ServerSubSystem),
		server:             server,
		netPlayers:         make([]*ServerPlayer, 0),
		subSystemsCallback: make(map[reflect.Type][]PacketServerCallback),
	}
	server.AddConnectionCallback(s.clientConnect)
	server.AddDisconnectionCallback(s.clientDisconnect)
	return s
}

func (s *ServerDelegator) clientConnect(client net.ServerClient) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	sp := &ServerPlayer{
		NetPlayer: &NetPlayer{},
		Client:    client,
	}
	s.netPlayers = append(s.netPlayers, sp)
	log.Printf("Client connected! %+v [%+v]", sp, s.netPlayers)
}

func (s *ServerDelegator) clientDisconnect(client net.ServerClient) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	result := slice_helpers.RemoveFromListEquals[*ServerPlayer](s.netPlayers, func(s *ServerPlayer) bool {
		return s.Client == client
	})
	if result == nil {
		log.Println("Could not find the client")
	} else {
		s.netPlayers = result
	}
}

func (s *ServerDelegator) Update() {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	for loops := 0; loops < 10; loops++ {
		noPackets := true
		//log.Printf("Clients: %+v", s.netPlayers)

		for _, p := range s.netPlayers {
			pack := misc.CheckErrorRetVal[packet_interface.Packet](p.ReadPacket())
			if pack != nil {
				log.Printf("Got packet: %+v", pack)
				noPackets = false
				for _, cb := range s.subSystemsCallback[reflect.TypeOf(pack)] {
					cb(p, s, pack)
				}
			}
		}
		if noPackets {
			break
		}
	}

	for _, sub := range s.subSystems {
		sub.Update()
	}
}

func (s *ServerDelegator) RegisterPacketCallback(cb PacketServerCallback, packet packet_interface.Packet) {
	t := reflect.TypeOf(packet)
	if _, ok := s.subSystemsCallback[t]; !ok {
		s.subSystemsCallback[t] = make([]PacketServerCallback, 0)
	}
	s.subSystemsCallback[t] = append(s.subSystemsCallback[t], cb)
}

func (s *ServerDelegator) RemovePacketCallback(cb PacketServerCallback, packetType packet_interface.Packet) {
}
func (s *ServerDelegator) RemoveSubSystem(name string) {
	delete(s.subSystems, name)
}

func (s *ServerDelegator) GetSubsystem(name string) interface{} {
	return s.subSystems[name]
}

func (s *ServerDelegator) RegisterSubSystem(name string, system ServerSubSystem) {
	s.subSystems[name] = system
	system.RegisterCallbacks(s)
}
