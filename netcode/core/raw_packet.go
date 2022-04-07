package core

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

type RawPacket struct {
	Data             []byte
	PacketId         uint64
	PacketTime       int64
	ContainingPacket Packet
}

func NewRawPacket(containingPacket Packet) *RawPacket {
	return &RawPacket{ContainingPacket: containingPacket}
}

func RawPacketFrom(data []byte) *RawPacket {
	pack := &RawPacket{Data: data}
	b := bytes.NewReader(data)
	id := uint32(0)
	if err := binary.Read(b, binary.LittleEndian, &id); err != nil {
		panic(err)
	}
	if err := binary.Read(b, binary.LittleEndian, &pack.PacketTime); err != nil {
		panic(err)
	}
	if err := binary.Read(b, binary.LittleEndian, &pack.PacketId); err != nil {
		panic(err)
	}

	newPacket := reflect.New(indexToPacket[id]).Interface().(Packet)
	packet := newPacket.FromReader(b)
	pack.ContainingPacket = packet.(Packet)
	return pack
}

func (r RawPacket) GetPacket() Packet {
	return r.ContainingPacket
}

func (r RawPacket) GetBytes() []byte {
	b := new(bytes.Buffer)
	id := packetToIndex[reflect.TypeOf(r.ContainingPacket)]
	if err := binary.Write(b, binary.LittleEndian, id); err != nil {
		panic(err)
	}
	if err := binary.Write(b, binary.LittleEndian, r.PacketTime); err != nil {
		panic(err)
	}
	if err := binary.Write(b, binary.LittleEndian, r.PacketId); err != nil {
		panic(err)
	}
	r.ContainingPacket.ToWriter(b)
	bb := b.Bytes()
	return bb
}
