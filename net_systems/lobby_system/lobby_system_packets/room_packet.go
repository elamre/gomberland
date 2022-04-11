package lobby_system_packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/elamre/gomberman/net_systems/common_system"
	"github.com/elamre/gomberman/net_systems/lobby_system/common"
)

const (
	_                       = iota
	RoomJoinAction          = iota
	RoomJoinSuccessAction   = iota
	RoomJoinFailedAction    = iota
	RoomLeaveAction         = iota
	RoomReadyAction         = iota
	RoomCreateAction        = iota
	RoomDeleteAction        = iota
	RoomCreateSuccessAction = iota
	RoomCreateFailedAction  = iota
	RoomUpdateAction        = iota
	RoomStartAction         = iota
)

func (r RoomAction) String() string {
	switch r {
	case RoomJoinAction:
		return "RoomJoinAction"
	case RoomJoinSuccessAction:
		return "RoomJoinSuccessAction"
	case RoomJoinFailedAction:
		return "RoomJoinFailedAction"
	case RoomLeaveAction:
		return "RoomLeaveAction"
	case RoomReadyAction:
		return "RoomReadyAction"
	case RoomCreateAction:
		return "RoomCreateAction"
	case RoomDeleteAction:
		return "RoomDeleteAction"
	case RoomCreateSuccessAction:
		return "RoomCreateSuccessAction"
	case RoomCreateFailedAction:
		return "RoomCreateFailedAction"
	case RoomUpdateAction:
		return "RoomUpdateAction"
	case RoomStartAction:
		return "RoomStartAction"
	}
	return "unknown"
}

type RoomAction uint32

type RoomPacket struct {
	UserId   uint32
	Action   RoomAction
	Password string
	Name     string
}

func (r RoomPacket) String() string {
	return fmt.Sprintf("Name: %s Pass: %s Action: %s Owner: %d", r.Name, r.Password, r.Action.String(), r.UserId)
}

func (c RoomPacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &c.Action); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &c.UserId); err != nil {
		panic(err)
	}

	stringLen := int32(0)

	if err := binary.Read(r, binary.LittleEndian, &stringLen); err != nil {
		panic(err)
	}
	name := make([]byte, stringLen)
	if _, err := r.Read(name); err != nil {
		panic(err)
	}
	c.Name = string(name)

	if err := binary.Read(r, binary.LittleEndian, &stringLen); err != nil {
		panic(err)
	}
	password := make([]byte, stringLen)
	if stringLen > 0 {
		if _, err := r.Read(password); err != nil {
			panic(err)
		}
	}
	c.Password = string(password)

	return c
}

func (c RoomPacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, c.Action); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, c.UserId); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, int32(len(c.Name))); err != nil {
		panic(err)
	}
	w.WriteString(c.Name)

	if err := binary.Write(w, binary.LittleEndian, int32(len(c.Password))); err != nil {
		panic(err)
	}
	if len(c.Password) > 0 {
		w.WriteString(c.Password)
	}
}

func (c RoomPacket) AckRequired() bool {
	return true
}

type RoomUpdatePacket struct {
	Rooms []*common.NetRoom
}

func (c RoomUpdatePacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, int32(len(c.Rooms))); err != nil {
		panic(err)
	}
	if len(c.Rooms) > 0 {
		for i := 0; i < len(c.Rooms); i++ {
			room := c.Rooms[i]
			if err := binary.Write(w, binary.LittleEndian, int32(len(room.RoomName))); err != nil {
				panic(err)
			}
			w.WriteString(room.RoomName)
			if err := binary.Write(w, binary.LittleEndian, int32(len(room.Players))); err != nil {
				panic(err)
			}
			for b := 0; b < len(room.Players); b++ {
				if err := binary.Write(w, binary.LittleEndian, int32(len(room.Players[b].Name))); err != nil {
					panic(err)
				}
				w.WriteString(room.Players[b].Name)
				if err := binary.Write(w, binary.LittleEndian, room.Players[b].Id); err != nil {
					panic(err)
				}
				if room.Players[b].Ready {
					w.WriteByte(1)
				} else {
					w.WriteByte(0)
				}

			}
		}
	}
}

func (c RoomUpdatePacket) FromReader(r *bytes.Reader) any {
	rooms := int32(0)
	stringLen := int32(0)
	if err := binary.Read(r, binary.LittleEndian, &rooms); err != nil {
		panic(err)
	}
	c.Rooms = make([]*common.NetRoom, rooms)
	for i := 0; i < len(c.Rooms); i++ {
		c.Rooms[i] = &common.NetRoom{}
		if err := binary.Read(r, binary.LittleEndian, &stringLen); err != nil {
			panic(err)
		}
		roomName := make([]byte, stringLen)
		if _, err := r.Read(roomName); err != nil {
			panic(err)
		}
		c.Rooms[i].RoomName = string(roomName)
		players := int32(0)
		if err := binary.Read(r, binary.LittleEndian, &players); err != nil {
			panic(err)
		}
		c.Rooms[i].Players = make([]*common_system.NetPlayer, players)
		for b := 0; b < len(c.Rooms[i].Players); b++ {
			c.Rooms[i].Players[b] = &common_system.NetPlayer{}
			if err := binary.Read(r, binary.LittleEndian, &stringLen); err != nil {
				panic(err)
			}
			playerName := make([]byte, stringLen)
			if _, err := r.Read(playerName); err != nil {
				panic(err)
			}
			c.Rooms[i].Players[b].Name = string(playerName)
			playerId := uint32(0)
			if err := binary.Read(r, binary.LittleEndian, &playerId); err != nil {
				panic(err)
			}
			c.Rooms[i].Players[b].Id = playerId
			ready, err := r.ReadByte()
			if err != nil {
				panic(err)
			}
			c.Rooms[i].Players[b].Ready = ready == 1
		}
	}
	return c
}

func (c RoomUpdatePacket) AckRequired() bool {
	return false
}
