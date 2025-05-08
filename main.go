package main

import (
	"gocard/duel"
	"gocard/menu"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

func main() {
	if !sdl.Init(sdl.InitVideo) {
		panic(sdl.GetError())
	}
	defer sdl.Quit()

	var window *sdl.Window
	var renderer *sdl.Renderer
	if !sdl.CreateWindowAndRenderer("GoCard", 1280, 720, sdl.WindowResizable, &window, &renderer) {
		panic(sdl.GetError())
	}
	defer sdl.DestroyRenderer(renderer)
	defer sdl.DestroyWindow(window)

	state := ui.StateStartMenu

	for state != ui.StateQuit {
		switch state {
		case ui.StateStartMenu:
			state = ui.AppState(menu.RenderStartMenu(renderer))
		case ui.StateDeckBuilder:
			state = ui.AppState(menu.RenderDeckBuilder(renderer))
		case ui.StateDuel:
			state = ui.AppState(duel.RenderDuel(renderer))
		default:
			state = ui.StateQuit
		}
	}
}
