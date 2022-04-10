package game_system_packets

import (
	"bytes"
	"encoding/binary"
)

type EntityStatePacket struct {
	EntityId    int
	EntityState int
	EntityType  int
}

func (e EntityStatePacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, e.EntityId); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.EntityState); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.EntityType); err != nil {
		panic(err)
	}

}
func (e EntityStatePacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &e.EntityId); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.EntityState); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.EntityType); err != nil {
		panic(err)
	}
	return e
}

func (e EntityStatePacket) AckRequired() bool {
	return true
}
