package intro_state

import (
	"github.com/elamre/gomberman/game/helpers"
	. "github.com/elamre/tentsuyu"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameIntroState struct {
	textOutOfArea bool
	textPos       float32
}

func NewGameIntroState() *GameIntroState {
	m := GameIntroState{}
	return &m
}

func (m *GameIntroState) Update(game *Game) error {

	return nil
}
func (m *GameIntroState) Draw(game *Game) error {
	m.textPos += 3
	if int(m.textPos) >= game.ScreenWidth() {
		m.textOutOfArea = true
	}
	ebitenutil.DebugPrintAt(game.Screen, "ELMAR", int(m.textPos), game.ScreenHeight()/2)
	return nil
}
func (m *GameIntroState) Msg() GameStateMsg {
	if m.textOutOfArea {
		return helpers.MainMenuState
	} else {
		return helpers.IntroState
	}
}

func (m *GameIntroState) SetMsg(state GameStateMsg) {

}
