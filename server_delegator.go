package main

import (
	. "github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"reflect"
)

type ServerDelegator struct {
	server             net.Server
	subSystems         map[string]SubSystem
	subSystemsCallback map[reflect.Type][]PacketCallback
}

func NewServerDelegator(server net.Server) *ServerDelegator {
	return &ServerDelegator{
		server:             server,
		subSystems:         make(map[string]SubSystem),
		subSystemsCallback: make(map[reflect.Type][]PacketCallback),
	}
}

func (s *ServerDelegator) RegisterSubSystem(name string, system SubSystem) {
	s.subSystems[name] = system
	system.RegisterCallbacks(s)
}

func (s *ServerDelegator) Update() {
	s.server.ClientIterator(func(c net.ServerClient) {
		pack, err := c.ReadPacket()
		if err != nil {
			panic(err)
		}
		if pack != nil {
			for _, cb := range s.subSystemsCallback[reflect.TypeOf(pack)] {
				cb(c, s, pack)
			}
		}
	})
	for _, sub := range s.subSystems {
		sub.Update()
	}
}

func (s *ServerDelegator) RegisterPacketCallback(cb PacketCallback, packet packet_interface.Packet) {
	t := reflect.TypeOf(packet)
	if _, ok := s.subSystemsCallback[t]; !ok {
		s.subSystemsCallback[t] = make([]PacketCallback, 0)
	}
	s.subSystemsCallback[t] = append(s.subSystemsCallback[t], cb)
}

func (s *ServerDelegator) RemovePacketCallback(cb PacketCallback, packetType packet_interface.Packet) {
}

func (s *ServerDelegator) RemoveSubSystem(name string) {
	delete(s.subSystems, name)
}

func (s *ServerDelegator) GetSubsystem(name string) interface{} {
	return s.subSystems[name]
}
