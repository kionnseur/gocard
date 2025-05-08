package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

type AppState int

const (
	StateStartMenu AppState = iota
	StateDeckBuilder
	StateDuel
	StateQuit
)

type Hud struct {
	Element
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
