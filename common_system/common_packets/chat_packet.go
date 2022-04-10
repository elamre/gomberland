package common_packets

import (
	"github.com/elamre/go_helpers/pkg/misc"

	"bytes"
)

type ChatPacket struct {
	From    uint32
	Message string
}

func (c ChatPacket) ToWriter(w *bytes.Buffer) {
	misc.CheckErrorRetVal(w.Write([]byte(c.Message)))

}

func (c ChatPacket) FromReader(r *bytes.Reader) any {
	message := make([]byte, r.Size())
	misc.CheckErrorRetVal(r.Read(message))
	c.Message = string(message)
	return c
}

func (c ChatPacket) AckRequired() bool {
	return true
}
