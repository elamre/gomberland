package main

import (
	"fmt"
	"github.com/elamre/go_helpers/pkg/misc"
	"github.com/elamre/gomberman/assets"
	"github.com/elamre/gomberman/game/helpers"
	"github.com/elamre/gomberman/game/intro_state"
	"github.com/elamre/gomberman/game/menu_state_revised"
	"github.com/elamre/tentsuyu"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var assetsManager *tentsuyu.AssetsManager

type GomberlandGame struct {
	game            *tentsuyu.Game
	menuState       *menu_state_revised.MenuState
	introState      *intro_state.GameIntroState
	close           bool
	previousMessage tentsuyu.GameStateMsg

	/*	curStateString, nextStateString string
	 */
}

func NewGomberlandGame() *GomberlandGame {
	g := GomberlandGame{}
	g.game = misc.CheckErrorRetVal[*tentsuyu.Game](tentsuyu.NewGame(640, 480))
	g.game.LoadAssetsManager(assets.GetManager)
	g.game.SetGameStateLoop(func() error {
		if err := g.checkState(); err != nil {
			return err
		}
		if g.close {
			return fmt.Errorf("normal termination")
		}
		return nil
	})
	g.menuState = menu_state_revised.NewMenuState()
	g.introState = intro_state.NewGameIntroState()
	g.setState(g.menuState)
	return &g
}

func (g *GomberlandGame) setState(gs tentsuyu.GameState) {
	g.game.SetGameState(gs)
	g.previousMessage = g.game.GetGameState().Msg()
}

func (g *GomberlandGame) checkState() error {
	if g.game.GetGameState().Msg() != g.previousMessage {
		log.Printf("%v <- %v", g.game.GetGameState().Msg(), g.previousMessage)
		switch g.previousMessage {
		case helpers.IntroState:
			g.setState(g.menuState)
		// should really only go to the menu
		case helpers.MainMenuState:
			switch g.game.GetGameState().Msg() {
			case helpers.QuitState:
				g.close = true
				return nil
			}
		}
		if g.game.GetGameState().Msg() != g.previousMessage {
			return fmt.Errorf("we are still in a different state, not implemented? %s != %s", g.game.GetGameState().Msg(), g.previousMessage)
		}
	}
	return nil

	/*	switch g.curState {
		case g.menuState:
			switch g.curState.Msg() {
			case menu_state_revised.InMenuState:
				// Do nothing
			case menu_state_revised.JoinRoomState:
				log.Println(" Changing to room ")
				g.nextState = g.roomState
			case menu_state_revised.BackState:
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
		}*/
}

func main() {
	defer assets.CleanUp()
	gomberland := NewGomberlandGame()
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(gomberland.game); err != nil {
		panic(err)
	}
}
