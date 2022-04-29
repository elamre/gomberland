package menu_state_revised

import (
	"github.com/elamre/debug_gui/pkg/common"
	_ "github.com/elamre/debug_gui/pkg/composition"
	"github.com/elamre/debug_gui/pkg/gui"
	"github.com/elamre/gomberman/game/helpers"
	"github.com/elamre/tentsuyu"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type mainMenu struct {
	mainContainer   *gui.Container
	buttonContainer *gui.Container
	curState        string
}

func newMainMenu(canvas *gui.Container) *mainMenu {
	m := &mainMenu{curState: "we"}
	m.mainContainer = gui.NewContainer(0, 0, 0.9, 0.9, gui.ANCHOR_HORIZONTAL_CENTER|gui.ANCHOR_VERTICAL_CENTER)
	m.mainContainer.SetParent(canvas)

	m.buttonContainer = gui.NewContainer(0, 0, 0.2, 1, gui.ANCHOR_HORIZONTAL_CENTER|gui.ALIGN_UP|gui.STACK_VERITCAL)
	m.buttonContainer.SetParent(m.mainContainer)
	m.curState = helpers.MainMenuState
	m.buttonContainer.AddElement(gui.NewElement(gui.NewDebugButton("Singleplayer", func(button *gui.DebugButton, stateChanger common.StateChanger, gameState tentsuyu.GameState) {
		m.curState = helpers.SinglePlayerMenuState
	}), 0, 40, gui.EXTEND_WIDTH|gui.ANCHOR_LEFT))
	m.buttonContainer.AddElement(gui.NewElement(gui.NewDebugButton("Multiplayer", func(button *gui.DebugButton, stateChanger common.StateChanger, gameState tentsuyu.GameState) {
		m.curState = helpers.MultiplayerMenuState
	}), 0, 80, gui.EXTEND_WIDTH|gui.ANCHOR_HORIZONTAL_CENTER))
	m.buttonContainer.AddElement(gui.NewElement(gui.NewDebugButton("Options", func(button *gui.DebugButton, stateChanger common.StateChanger, gameState tentsuyu.GameState) {
		m.curState = helpers.OptionsMenuState
	}), 0, 80, gui.EXTEND_WIDTH|gui.ANCHOR_HORIZONTAL_CENTER))
	m.buttonContainer.AddElement(gui.NewElement(gui.NewDebugButton("Exit", func(button *gui.DebugButton, stateChanger common.StateChanger, gameState tentsuyu.GameState) {
		m.curState = helpers.QuitState
	}), 0, 120, gui.EXTEND_WIDTH|gui.ANCHOR_LEFT))
	return m
}

func (m *mainMenu) GetState() string {
	return m.curState
}

func (m *mainMenu) Init() any {
	return nil
}

func (m *mainMenu) Update() {

}

func (m *mainMenu) Draw(game *tentsuyu.Game) {
	ebitenutil.DebugPrintAt(game.Screen, "TTEESSTT", 10, 10)
}

func (m *mainMenu) DeInit() any {
	return nil
}
