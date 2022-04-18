package game

import "github.com/elamre/tentsuyu"

type MenuUiInterface interface {
	Init() any
	DeInit() any
	Update()
	Draw(game *tentsuyu.Game)
	GetState() string
}
