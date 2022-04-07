package packet

import (
	"bytes"
	"encoding/binary"
)

type AckPacket struct {
	AckId uint32
}

func (c AckPacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, c.AckId); err != nil {
		panic(err)
	}

}
func (c AckPacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &c.AckId); err != nil {
		panic(err)
	}
	return c
}

func (c AckPacket) AckRequired() bool {
	return false
}
