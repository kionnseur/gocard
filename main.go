package main

import (
	"gocard/data"
	"gocard/deck_builder"
	"gocard/duel"
	"gocard/start_menu"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

var (
	screenWidth  int32 = 1280
	screenHeight int32 = 720
)

func main() {

	var window *sdl.Window
	var renderer *sdl.Renderer
	if !sdl.CreateWindowAndRenderer("GoCard", screenWidth, screenHeight, sdl.WindowResizable, &window, &renderer) {
		panic(sdl.GetError())
	}
	defer sdl.DestroyRenderer(renderer)
	defer sdl.DestroyWindow(window)

	state := ui.AppState{State: ui.StateStartMenu}

	for state.State != ui.StateQuit {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		switch state.State {
		case ui.StateStartMenu:
			state = ui.AppState(start_menu.RenderStartMenu(renderer))
		case ui.StateDeckBuilder:
			state = ui.AppState(deck_builder.RenderDeckBuilder(renderer, window))
		case ui.StateDuel:
			state = ui.AppState(duel.RenderDuel(renderer))
		default:
			state = ui.AppState{State: ui.StateQuit}
		}
	}
}
