package core

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
)

var packetToIndex = map[reflect.Type]int{}
var indexToPacket = map[int]reflect.Type{}
var packetIndex = 0

type Packet interface {
	ToWriter(w *bytes.Buffer)
	FromReader(r *bytes.Reader)
}

func RegisterPacket(pack any) int {
	tt := reflect.TypeOf(pack)
	if id, ok := packetToIndex[tt]; ok {
		panic(fmt.Sprintf("Packet %T already registered with id: %d", tt, id))
	}
	packetIndex++
	packetToIndex[tt] = packetIndex
	indexToPacket[packetIndex] = tt
	return packetIndex
}

func CreatePacket(index int) Packet {
	return reflect.New(indexToPacket[index]).Interface().(Packet)
}

func BytesToPacket(data []byte) Packet {
	log.Printf("creating from bytes:  %02x", data)
	packet := CreatePacket(int(data[0]))
	bb := bytes.NewReader(data[1:])
	packet.FromReader(bb)
	return packet
}

func PacketToBytes(packet Packet) []byte {
	b := new(bytes.Buffer)
	id := packetToIndex[reflect.TypeOf(packet)]
	err := b.WriteByte(byte(id))
	log.Printf("%v", err)
	packet.ToWriter(b)
	return b.Bytes()
}
