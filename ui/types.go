package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

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

type Element interface {
	GetRect() *sdl.FRect
	Draw(renderer *sdl.Renderer)
}

type Hud struct {
	Rect  sdl.FRect
	Color sdl.Color
}

type TextBox struct {
	Rect      sdl.FRect
	Color     sdl.Color
	Text      string
	TextColor sdl.Color
	Font      *ttf.Font
}

type Button struct {
	Rect      sdl.FRect
	Color     sdl.Color
	Text      string
	TextColor sdl.Color
	Font      *ttf.Font
	OnClick   func() AppState
}
