package packets

import (
	"bytes"
	"encoding/binary"
)

const (
	_                             = iota
	ConnectionRegisterAction      = iota
	ConnectionAcceptedAction      = iota
	ConnectionRefusedAction       = iota
	ConnectionNewRegisteredAction = iota
)

type ConnectionAction uint32

type ConnectionPacket struct {
	UserId  uint32
	Action  ConnectionAction
	Message string
}

func NewRegisterPacket(name string) ConnectionPacket {
	return ConnectionPacket{Action: ConnectionRegisterAction, Message: name}
}

func (c ConnectionPacket) ConnectionSuccessful() bool {
	return c.Action == ConnectionAcceptedAction
}

func (c ConnectionPacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, c.UserId); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, c.Action); err != nil {
		panic(err)
	}
	w.WriteString(c.Message)
}
func (c ConnectionPacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &c.UserId); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &c.Action); err != nil {
		panic(err)
	}
	sstring := make([]byte, r.Len())
	if _, err := r.Read(sstring); err != nil {
		panic(err)
	}
	c.Message = string(sstring)
	return c
}

func (c ConnectionPacket) AckRequired() bool {
	return true
}
