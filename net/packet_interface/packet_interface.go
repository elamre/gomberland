package packet_interface

import (
	"fmt"
	"reflect"
)

var packetToIndex = map[reflect.Type]uint32{}
var indexToPacket = map[uint32]reflect.Type{}
var packetIndex = uint32(0)

func GetType(t uint32) reflect.Type {
	return indexToPacket[t]
}

func GetIndex(t reflect.Type) uint32 {
	return packetToIndex[t]
}

func RegisterPacket(pack any) uint32 {
	tt := reflect.TypeOf(pack)
	if id, ok := packetToIndex[tt]; ok {
		panic(fmt.Sprintf("Packet %T already registered with id: %d", tt, id))
	}
	packetIndex++
	packetToIndex[tt] = packetIndex
	indexToPacket[packetIndex] = tt
	return packetIndex
}

func GetRegisteredPackets() map[uint32]reflect.Type {
	return indexToPacket
}
