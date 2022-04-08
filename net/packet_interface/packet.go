package packet_interface

import (
	"bytes"
	"reflect"
)

type Packet interface {
	ToWriter(w *bytes.Buffer)
	FromReader(r *bytes.Reader) any
	AckRequired() bool
}

func CreatePacket(index uint32) Packet {
	return reflect.New(indexToPacket[index]).Interface().(Packet)
}

func BytesToPacket(data []byte) Packet {
	packet := CreatePacket(uint32(data[0]))
	bb := bytes.NewReader(data[1:])
	return packet.FromReader(bb).(Packet)
}

func PacketToBytes(packet Packet) []byte {
	b := new(bytes.Buffer)
	id := packetToIndex[reflect.TypeOf(packet)]
	err := b.WriteByte(byte(id))
	if err != nil {
		panic(err)
	}
	packet.ToWriter(b)
	return b.Bytes()
}
