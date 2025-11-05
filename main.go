package main

import (
	"gocard/data"
	"gocard/deck_builder"
	"gocard/duel"
	"gocard/start_menu"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

// Initializes the SDL window and renderer, then runs the main game loop
// handling different application states.
func main() {
	// Initialize SDL window and renderer
	var window *sdl.Window
	var renderer *sdl.Renderer
	if !sdl.CreateWindowAndRenderer("GoCard", data.ScreenWidth, data.ScreenHeight, sdl.WindowResizable, &window, &renderer) {
		panic(sdl.GetError())
	}
	defer sdl.DestroyRenderer(renderer)
	defer sdl.DestroyWindow(window)

	// Initial application state: start menu
	state := ui.AppState{State: ui.StateStartMenu}

	// Main game loop
	for state.State != ui.StateQuit {
		// Update window size
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		// Handle different application states
		switch state.State {
		case ui.StateStartMenu:
			state = *start_menu.RenderStartMenu(renderer)
		case ui.StateDeckMenu:
			state = *deck_builder.RenderDeckMenu(renderer, window, &state)
		case ui.StateDeckBuilder:
			state = *deck_builder.RenderDeckEditor(renderer, window, &state)
		case ui.StateDuel:
			state = *duel.RenderDuel(renderer)
		default:
			state = ui.AppState{State: ui.StateQuit}
		}
	}
}
