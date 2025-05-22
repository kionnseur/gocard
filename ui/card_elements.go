package ui

import (
	"gocard/data"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type UICard interface {
	GetRect() sdl.FRect
	GetCard() data.Card
	Element
}

// Draw implements Element.

type UIMonsterCard struct {
	Rect sdl.FRect
	Card *data.MonsterCard
}

// GetCard implements UICard.
func (m *UIMonsterCard) GetCard() data.Card { return m.Card }

func (m *UIMonsterCard) GetDescription() string { return m.Card.Description }

func (m *UIMonsterCard) GetImage() string { return m.Card.Image }

func (m *UIMonsterCard) GetName() string { return m.Card.Name }

func (m *UIMonsterCard) GetRect() sdl.FRect { return m.Rect }

func (m *UIMonsterCard) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
	sdl.RenderFillRect(renderer, &m.Rect)

	nameBox := TextBox{
		Text:      m.Card.GetName(),
		Rect:      sdl.FRect{X: m.Rect.X + 5, Y: m.Rect.Y + 5, W: m.Rect.W - 10, H: 20},
		Color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
		Font:      GetDefaultFont(13),
		TextColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
	nameBox.Draw(renderer)
}

type UISpellTrapCard struct {
	Rect sdl.FRect
	Card data.SpellTrapCard
}

// GetCard implements UICard.
func (m *UISpellTrapCard) GetCard() data.Card { return &m.Card }

func (m *UISpellTrapCard) GetDescription() string { return m.Card.Description }

func (m *UISpellTrapCard) GetImage() string { return m.Card.Image }

func (m *UISpellTrapCard) GetName() string { return m.Card.Name }

func (m *UISpellTrapCard) GetRect() sdl.FRect { return m.Rect }

func (m *UISpellTrapCard) Draw(renderer *sdl.Renderer) {
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

func CreateUICard(card data.Card, rect sdl.FRect) UICard {
	if c, ok := card.(*data.MonsterCard); ok {
		return &UIMonsterCard{
			Rect: rect,
			Card: c,
		}
	}
	if c, ok := card.(*data.SpellTrapCard); ok {
		return &UISpellTrapCard{
			Rect: rect,
			Card: *c,
		}
	}
	return nil
}
