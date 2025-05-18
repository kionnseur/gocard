package data

import (
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type UICard interface {
	Draw(renderer *sdl.Renderer)
}

type UIMonsterCard struct {
	Rect sdl.FRect
	Card Card
}

func (m *UIMonsterCard) Draw(renderer *sdl.Renderer) {
	// Utilise m.CardRect pour dessiner à la bonne position
	// sdl.SetRenderDrawColor(renderer, 0, 0, 0, 255)
	// sdl.RenderRect(renderer, &m.Rect)
	sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
	sdl.RenderFillRect(renderer, &m.Rect)

	// Affiche le nom, description, etc. dans la carte
	nameBox := ui.TextBox{
		Text:      m.Card.GetName(),
		Rect:      sdl.FRect{X: m.Rect.X + 5, Y: m.Rect.Y + 5, W: m.Rect.W - 10, H: 20},
		Color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
		Font:      ui.GetDefaultFont(),
		TextColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
	nameBox.Draw(renderer)
}
