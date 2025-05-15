package start_menu

import (
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
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
					if x > btn.Rect.X && x < btn.Rect.X+btn.Rect.W &&
						y > btn.Rect.Y && y < btn.Rect.Y+btn.Rect.H {
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

func getStartMenuButtons() []ui.Button {
	ttf.Init()

	font := ttf.OpenFont("assets/fonts/arial.ttf", 24)

	return []ui.Button{
		{
			Rect:      sdl.FRect{X: 140.0, Y: 80.0, W: 200.0, H: 50.0},
			Color:     sdl.Color{R: 150, G: 150, B: 255, A: 255},
			Text:      "Deck Builder",
			TextColor: sdl.Color{R: 105, G: 105, B: 0, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.AppState{State: ui.StateDeckMenu} },
		},
		{
			Rect:      sdl.FRect{X: 140, Y: 180, W: 200, H: 50},
			Color:     sdl.Color{R: 0, G: 250, B: 0, A: 255},
			Text:      "Duel",
			TextColor: sdl.Color{R: 255, G: 5, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.AppState{State: ui.StateDuel} },
		},
	}
}
