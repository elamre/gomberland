package packets

import (
	"bytes"
	"encoding/binary"
	"github.com/elamre/gomberman/common_system"
	"github.com/elamre/gomberman/lobby_system/common"
)

const (
	_                       = iota
	RoomJoinAction          = iota
	RoomJoinSuccessAction   = iota
	RoomJoinFailedAction    = iota
	RoomLeaveAction         = iota
	RoomCreateAction        = iota
	RoomCreateSuccessAction = iota
	RoomCreateFailedAction  = iota
	RoomUpdateAction        = iota
	RoomStartAction         = iota
)

type RoomAction uint32

type RoomPacket struct {
	UserId   uint32
	Action   RoomAction
	Password string
	Name     string
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
		}
	}
	return c
}

func (c RoomUpdatePacket) AckRequired() bool {
	return false
}
