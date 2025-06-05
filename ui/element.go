package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
)

type Element interface {
	GetRect() *sdl.FRect
	Draw(renderer *sdl.Renderer)
}

type State int

const (
	StateStartMenu State = iota
	StateDeckMenu
	StateDeckBuilder
	StateDuel
	StateQuit
)

type AppState struct {
	State
	Data map[string]string
}
