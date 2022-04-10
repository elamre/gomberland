package main

import (
	. "github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/net"
	"github.com/elamre/gomberman/net/packet_interface"
	"reflect"
)

type ClientDelegator struct {
	client             net.Client
	subSystems         map[string]ClientSubSystem
	subSystemsCallback map[reflect.Type][]PacketClientCallback
}

func NewClientDelegator(client net.Client) *ClientDelegator {
	return &ClientDelegator{
		client:             client,
		subSystems:         make(map[string]ClientSubSystem),
		subSystemsCallback: make(map[reflect.Type][]PacketClientCallback),
	}
}

func (s *ClientDelegator) Update() {
	pack, err := s.client.ReadPacket()
	if err != nil {
		panic(err)
	}
	if pack != nil {
		for _, cb := range s.subSystemsCallback[reflect.TypeOf(pack)] {
			cb(s.client, s, pack)
		}
	}
	for _, sub := range s.subSystems {
		sub.Update()
	}
}

func (s *ClientDelegator) RegisterPacketCallback(cb PacketClientCallback, packet packet_interface.Packet) {
	t := reflect.TypeOf(packet)
	if _, ok := s.subSystemsCallback[t]; !ok {
		s.subSystemsCallback[t] = make([]PacketClientCallback, 0)
	}
	s.subSystemsCallback[t] = append(s.subSystemsCallback[t], cb)
}

func (s *ClientDelegator) RemovePacketCallback(cb PacketClientCallback, packetType packet_interface.Packet) {
}

func (s *ClientDelegator) RemoveSubSystem(name string) {
	delete(s.subSystems, name)
}

func (s *ClientDelegator) GetSubsystem(name string) interface{} {
	return s.subSystems[name]
}

func (s *ClientDelegator) RegisterSubSystem(name string, system ClientSubSystem) {
	s.subSystems[name] = system
	system.RegisterCallbacks(s)
}
