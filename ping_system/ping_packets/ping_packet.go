package ping_packets

import (
	"bytes"
	"encoding/binary"
	"time"
)

type PingPacket struct {
	CreationTime  time.Time
	receptionTime time.Time
}

func (e PingPacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, e.CreationTime.UnixMilli()); err != nil {
		panic(err)
	}
}

func (e PingPacket) FromReader(r *bytes.Reader) any {
	e.receptionTime = time.Now()
	ttime := int64(0)
	if err := binary.Read(r, binary.LittleEndian, &ttime); err != nil {
		panic(err)
	}
	e.CreationTime = time.UnixMilli(ttime)
	return e
}

func (e PingPacket) AckRequired() bool {
	return false
}

func (e PingPacket) GetPing() int64 {
	return e.receptionTime.UnixMilli() - e.CreationTime.UnixMilli()
}
