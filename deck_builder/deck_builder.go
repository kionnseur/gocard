package deck_builder

import (
	"gocard/data"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

var (
	// deckBuilderData map[string]string
	scrollLevel float32 = 0.0
	gap         float32 = 20
	fontSize    float32 = 24
	_, font             = ttf.Init(), ttf.OpenFont("assets/fonts/arial.ttf", fontSize)
)

// Déclare le HUD de la listview à l'extérieur
var listViewHud = &ui.Hud{
	Rect:  sdl.FRect{X: gap, Y: 80, W: 200, H: 300},
	Color: sdl.Color{R: 100, G: 100, B: 100, A: 50},
}

func RenderDeckBuilder(renderer *sdl.Renderer, window *sdl.Window) ui.AppState {
	var buttons []ui.Button
	ui_deckList := getDeckListElements(data.GetDeckList())
	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)

		buttons = getDeckBuilderButtons()

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
			case sdl.EventMouseWheel:
				y := event.Wheel().Y
				scrollLevel -= float32(y) * gap
				if scrollLevel < 0 {
					scrollLevel = 0
				}
				maxScroll := float32(len(ui_deckList))*35 - listViewHud.Rect.H
				if scrollLevel > maxScroll {
					scrollLevel = maxScroll
				}
			}
		}

		sdl.SetRenderDrawColor(renderer, 255, 165, 0, 255)
		sdl.RenderClear(renderer)

		// Dessine le HUD de la listview
		listViewHud.Draw(renderer)

		// Dessine les éléments de la liste, scroll/clipping selon le HUD
		for _, e := range ui_deckList {
			if tb, ok := e.(*ui.TextBox); ok {
				rect := tb.Rect
				rect.Y -= scrollLevel
				if rect.Y+rect.H > listViewHud.Rect.Y && rect.Y < listViewHud.Rect.Y+listViewHud.Rect.H {
					tmp := *tb
					tmp.Rect = rect
					tmp.Draw(renderer)
				}
			}
		}
		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		sdl.RenderPresent(renderer)
		sdl.Delay(16)
	}
}

func getDeckBuilderButtons() []ui.Button {

	return []ui.Button{
		{
			Rect:      sdl.FRect{X: listViewHud.Rect.X, Y: listViewHud.Rect.Y - 30, W: 200, H: 30},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 100},
			Text:      "Nouveau Deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.StateStartMenu },
		},
		{
			Rect:      sdl.FRect{X: float32(data.ScreenWidth) - 50, Y: 0, W: 50, H: 50},
			Color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
			Text:      "Retour ⬅️",
			TextColor: sdl.Color{R: 255, G: 0, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.StateStartMenu },
		},
	}
}

func getDeckListElements(decksList []data.Deck) []ui.Element {
	elements := make([]ui.Element, len(decksList))
	for i, deck := range decksList {
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = &ui.TextBox{
			Rect:      sdl.FRect{X: listViewHud.Rect.X, Y: listViewHud.Rect.Y + float32(i*35), W: listViewHud.Rect.W, H: 30},
			Color:     sdl.Color{R: r, G: g, B: b, A: 255},
			Text:      deck.Name,
			TextColor: sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			Font:      font,
		}
	}
	return elements
}
