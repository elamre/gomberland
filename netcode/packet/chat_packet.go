package packet

import (
	"bytes"
	"log"
)

type ChatPacket struct {
	Message string
}

func (c *ChatPacket) ToWriter(w *bytes.Buffer) {
	_, err := w.Write([]byte(c.Message))
	log.Printf("err: %v", err)

}

func (c *ChatPacket) FromReader(r *bytes.Reader) {
	message := make([]byte, r.Size())
	r.Read(message)
	c.Message = string(message)
}
