package menu_state_revised

import (
	"github.com/elamre/debug_gui/pkg/gui"
	. "github.com/elamre/tentsuyu"
)

type MenuState struct {
	baseContainer *gui.Container
	mainMenu      *mainMenu
}

func NewMenuState() *MenuState {
	m := MenuState{}
	m.baseContainer = gui.NewBaseContainer(0, 0, 640, 480)
	m.mainMenu = newMainMenu(m.baseContainer)
	return &m
}

func (m *MenuState) Update(game *Game) error {
	m.baseContainer.Update(game.Input, nil, m)
	return nil
}
func (m *MenuState) Draw(game *Game) error {
	m.baseContainer.Draw(game.Screen, game.DefaultCamera)
	m.mainMenu.Draw(game)

	return nil
}
func (m *MenuState) Msg() GameStateMsg { return GameStateMsg(m.mainMenu.GetState()) }
func (m *MenuState) SetMsg(state GameStateMsg) {

}
