package start_menu

import (
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

func RenderStartMenu(renderer *sdl.Renderer) ui.AppState {
	buttons := getStartMenuButtons()

	for {
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				x, y := event.Button().X, event.Button().Y
				for _, btn := range buttons {
					if x > btn.GetRect().X && x < btn.GetRect().X+btn.GetRect().W &&
						y > btn.GetRect().Y && y < btn.GetRect().Y+btn.GetRect().H {
						return btn.OnClick()
					}
				}
			}
		}
		sdl.SetRenderDrawColor(renderer, 200, 200, 200, 255)
		sdl.RenderClear(renderer)

		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		sdl.RenderPresent(renderer)
	}
}

func getStartMenuButtons() []*ui.Button {
	font := ui.GetDefaultFont(24)

	return []*ui.Button{
		ui.NewButton(
			"Deck Builder",
			sdl.FRect{X: 140.0, Y: 80.0, W: 200.0, H: 50.0},
			sdl.Color{R: 150, G: 150, B: 255, A: 255},
			sdl.Color{R: 105, G: 105, B: 0, A: 255},
			font,
			func() ui.AppState { return ui.AppState{State: ui.StateDeckMenu} },
		),
		ui.NewButton(
			"Duel",
			sdl.FRect{X: 140, Y: 180, W: 200, H: 50},
			sdl.Color{R: 0, G: 250, B: 0, A: 255},
			sdl.Color{R: 255, G: 5, B: 255, A: 255},
			font,
			func() ui.AppState { return ui.AppState{State: ui.StateDuel} },
		),
	}
}
