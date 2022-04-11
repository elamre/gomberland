package common_packets

import (
	"bytes"
	"encoding/binary"
	"github.com/elamre/go_helpers/pkg/misc"
)

const (
	_               = iota
	GameStateMenu   = iota
	GameStateLobby  = iota
	GameStateRoom   = iota
	GameStateInGame = iota
	GameStateScore  = iota
)

type GameState uint32

type GameStateChangePacket struct {
	GameState GameState
}

func (c GameStateChangePacket) ToWriter(w *bytes.Buffer) {
	misc.CheckError(binary.Write(w, binary.LittleEndian, c.GameState))

}
func (c GameStateChangePacket) FromReader(r *bytes.Reader) any {
	misc.CheckError(binary.Read(r, binary.LittleEndian, &c.GameState))
	return c
}

func (c GameStateChangePacket) AckRequired() bool {
	return true
}
