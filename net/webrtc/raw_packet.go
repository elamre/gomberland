package webrtc

import (
	"bytes"
	"encoding/binary"
	"github.com/elamre/gomberman/net/packet_interface"
	"reflect"
)

type RawPacket struct {
	Data             []byte
	PacketId         uint64
	PacketTime       int64
	ContainingPacket packet_interface.Packet
}

func NewRawPacket(containingPacket packet_interface.Packet) *RawPacket {
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

	newPacket := reflect.New(packet_interface.GetType(id)).Interface().(packet_interface.Packet)
	packet := newPacket.FromReader(b)
	pack.ContainingPacket = packet.(packet_interface.Packet)
	return pack
}

func (r RawPacket) GetPacket() packet_interface.Packet {
	return r.ContainingPacket
}

func (r RawPacket) GetBytes() []byte {
	b := new(bytes.Buffer)
	id := packet_interface.GetIndex(reflect.TypeOf(r.ContainingPacket))
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
