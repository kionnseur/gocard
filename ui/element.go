package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
)

// Element is the common interface for all UI elements.
type Element interface {
	GetRect() *sdl.FRect         // GetRect returns the element's rectangle.
	Draw(renderer *sdl.Renderer) // Draw renders the element.
}

// State represents the application's current state.
type State int

const (
	StateStartMenu   State = iota //  start menu.
	StateDeckMenu                 //  deck selection menu.
	StateDeckBuilder              //  deck building.
	StateDuel                     //  duel.
	StateQuit                     //  application should quit.
)

// AppState holds the current state and optional data.
type AppState struct {
	State
	Data map[string]string
}
