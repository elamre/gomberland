package packet

import "io"

type EntityStatePacket struct {
	packetId    int
	EntityId    int
	EntityState int
	EntityType  int
}

func (e *EntityStatePacket) ToWriter(w *io.Writer) {

}
func (e *EntityStatePacket) FromReader(r *io.Reader) {

}
func (e *EntityStatePacket) GetId() int {
	return 1
}
