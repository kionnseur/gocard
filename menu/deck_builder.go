package menu

import (
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

func RenderDeckBuilder(renderer *sdl.Renderer) ui.AppState {
	buttons := getDeckBuilderButtons()

	for {
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.StateQuit
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
		sdl.SetRenderDrawColor(renderer, 255, 165, 0, 255)
		sdl.RenderClear(renderer)

		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		sdl.RenderPresent(renderer)
		sdl.Delay(16)
	}
}

func getDeckBuilderButtons() []ui.Button {
	ttf.Init()

	font := ttf.OpenFont("assets/fonts/arial.ttf", 24)

	return []ui.Button{
		{
			Rect:      sdl.FRect{X: 140.0, Y: 80.0, W: 200.0, H: 50.0},
			Color:     sdl.Color{R: 255, G: 255, B: 0, A: 255},
			Text:      "Nouveau Deck",
			TextColor: sdl.Color{R: 0, G: 0, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { print("New Deck"); return ui.StateStartMenu },
		},
		{
			Rect:      sdl.FRect{X: 140, Y: 180, W: 200, H: 50},
			Color:     sdl.Color{R: 255, G: 0, B: 0, A: 255},
			Text:      "Edit Deck",
			TextColor: sdl.Color{R: 0, G: 255, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { print("Edit Deck"); return ui.StateStartMenu },
		},
		{
			Rect:      sdl.FRect{X: 140, Y: 280, W: 200, H: 50},
			Color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
			Text:      "Retour ⬅️",
			TextColor: sdl.Color{R: 255, G: 0, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { print("Retour"); return ui.StateStartMenu },
		},
	}
}
