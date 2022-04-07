package packet

import (
	"bytes"
)

type ChatPacket struct {
	Message string
}

func (c ChatPacket) ToWriter(w *bytes.Buffer) {
	_, err := w.Write([]byte(c.Message))
	if err != nil {
		panic(err)
	}
}

func (c ChatPacket) FromReader(r *bytes.Reader) any {
	message := make([]byte, r.Size())
	r.Read(message)
	c.Message = string(message)
	return c
}

func (c ChatPacket) AckRequired() bool {
	return true
}
