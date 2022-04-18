package menu_state

import (
	"github.com/elamre/gomberman/game"
	. "github.com/elamre/tentsuyu"
	"log"
)

const BackState = "menu_state_back"
const InMenuState = "menu_state"
const OptionsState = "menu_state_options"
const JoinRoomState = "menu_state_join"

const (
	mainMenuState         = "main_menu"
	optionsMenuState      = "options_menu"
	singlePlayerMenuState = "singlePlayer_menu"
	multiPlayerMenuState  = "multiPlayer_menu"
)

type MenuState struct {
	curState         GameStateMsg
	inMenu           string
	stateToObjectMap map[string]game.MenuUiInterface
	curMenu          game.MenuUiInterface
}

func NewMenuState() *MenuState {
	m := MenuState{curState: "reset", inMenu: mainMenuState}
	m.stateToObjectMap = map[string]game.MenuUiInterface{mainMenuState: newMainMenu(), multiPlayerMenuState: newlobbyState()}
	return &m
}

func (m *MenuState) Update(*Game) error {
	if m.curState == "reset" || m.curMenu == nil {
		log.Println("Is reset?")
		m.SetMsg("reset")
	}
	m.curMenu.Update()
	if m.curMenu != m.stateToObjectMap[m.curMenu.GetState()] {
		log.Printf("Switching to: %s", m.curMenu.GetState())
		m.curMenu = m.stateToObjectMap[m.curMenu.GetState()]

	}
	return nil
}
func (m *MenuState) Draw(game *Game) error {
	if m.curMenu != nil {
		m.curMenu.Draw(game)
	}
	return nil
}
func (m *MenuState) Msg() GameStateMsg { return m.curState }
func (m *MenuState) SetMsg(state GameStateMsg) {
	if state == "reset" {
		// Reset everything here
		state = InMenuState
		/*for _, v := range m.stateToObjectMap {
			v.DeInit()
			v.Init()
		}*/
		m.curMenu = m.stateToObjectMap[mainMenuState]
		log.Printf("Reset to: %s", m.curMenu.GetState())
	}
	m.curState = state
}
