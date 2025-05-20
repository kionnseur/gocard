package ui

import (
	"gocard/data"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

type UICard interface {
	Draw(renderer *sdl.Renderer)
	GetRect() sdl.FRect
}

type UIMonsterCard struct {
	Rect sdl.FRect
	Card data.Card
	Font *ttf.Font
}

func (m *UIMonsterCard) GetRect() sdl.FRect {
	return m.Rect
}

func (m *UIMonsterCard) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
	sdl.RenderFillRect(renderer, &m.Rect)

	nameBox := TextBox{
		Text:      m.Card.GetName(),
		Rect:      sdl.FRect{X: m.Rect.X + 5, Y: m.Rect.Y + 5, W: m.Rect.W - 10, H: 20},
		Color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
		Font:      GetDefaultFont(20),
		TextColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
	nameBox.Draw(renderer)
}
