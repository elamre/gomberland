package main

import (
	"github.com/elamre/go_helpers/pkg/misc"
	"github.com/elamre/gomberman/assets"
	"github.com/elamre/gomberman/game/menu_state"
	"github.com/elamre/gomberman/game/room_state"
	"github.com/elamre/tentsuyu"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var assetsManager *tentsuyu.AssetsManager

type GomberlandGame struct {
	game      *tentsuyu.Game
	curState  tentsuyu.GameState
	nextState tentsuyu.GameState

	menuState *menu_state.MenuState
	roomState *room_state.RoomState

	curStateString, nextStateString string
}

func NewGomberlandGame() *GomberlandGame {
	g := GomberlandGame{}
	g.game = misc.CheckErrorRetVal[*tentsuyu.Game](tentsuyu.NewGame(640, 480))
	g.game.LoadAssetsManager(assets.GetManager)
	g.game.SetGameStateLoop(func() error {
		g.checkState()
		return nil
	})
	g.menuState = menu_state.NewMenuState()
	g.roomState = room_state.NewRoomState()
	g.nextState = g.menuState
	return &g
}

func (g *GomberlandGame) checkState() {
	switch g.curState {
	case g.menuState:
		switch g.curState.Msg() {
		case menu_state.InMenuState:
			// Do nothing
		case menu_state.JoinRoomState:
			log.Println(" Changing to room ")
			g.nextState = g.roomState
		case menu_state.BackState:
			// TODO close the game
		default:
			//Do nothing
		}
	case g.roomState:
		switch g.curState.Msg() {
		case room_state.InRoomState:
		// Do nothing
		case room_state.LaunchGameState:
			log.Printf("Launching with: %+v", g.roomState.GetSettings())
		}
	}
	if g.curState != g.nextState {
		if g.curState != nil {
			g.curState.SetMsg("reset")
		}
		g.game.SetGameState(g.nextState)
		g.curState = g.nextState
	}
}

func main() {
	defer assets.CleanUp()
	gomberland := NewGomberlandGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.RunGame(gomberland.game)
}
