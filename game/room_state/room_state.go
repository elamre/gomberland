package room_state

import (
	"github.com/elamre/gomberman/game"
	. "github.com/elamre/tentsuyu"
)

const BackState = "room_state_back"
const InRoomState = "room_state"
const LaunchGameState = "room_state_start"

type RoomState struct {
	curState GameStateMsg
}

func NewRoomState() *RoomState {
	m := RoomState{curState: InRoomState}
	return &m
}

func (m *RoomState) Update(*Game) error {
	return nil
}
func (m *RoomState) Draw(*Game) error  { return nil }
func (m *RoomState) Msg() GameStateMsg { return m.curState }

func (m *RoomState) SetMsg(state GameStateMsg) {
	if state == "reset" {
		// Reset everything here
		m.curState = InRoomState
	}
	m.curState = state
}

func (m *RoomState) GetSettings() game.Settings {
	return game.Settings{
		MapIdentifier: "IceMap",
		PlayerNames:   []string{"p1", "p2"},
		GameMode:      "co-op",
	}
}
