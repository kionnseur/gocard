package ui

import (
	"gocard/data"
	"strconv"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type UICard interface {
	GetCard() data.Card
	Element
}

// Draw implements Element.

type UIMonsterCard struct {
	rect     sdl.FRect
	card     *data.MonsterCard
	quantity int
}

// GetCard implements UICard.
func (m *UIMonsterCard) GetCard() data.Card { return m.card }

func (m *UIMonsterCard) GetDescription() string { return m.card.GetDescription() }

func (m *UIMonsterCard) GetImage() string { return m.card.GetImage() }

func (m *UIMonsterCard) GetName() string { return m.card.GetName() }

func (m *UIMonsterCard) GetRect() *sdl.FRect { return &m.rect }

func (m *UIMonsterCard) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
	sdl.RenderFillRect(renderer, &m.rect)

	nameBox := TextBox{
		text:      m.card.GetName(),
		rect:      sdl.FRect{X: m.rect.X + 5, Y: m.rect.Y + 5, W: m.rect.W - 10, H: 20},
		color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
		font:      GetDefaultFont(13),
		textColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
	if m.quantity > 1 {
		qtyBox := TextBox{
			text:      strconv.Itoa(m.quantity),
			rect:      sdl.FRect{X: m.rect.X + 5, Y: m.rect.Y + 25, W: m.rect.W - 10, H: 20},
			color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
			font:      GetDefaultFont(20),
			textColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
		}
		qtyBox.Draw(renderer)
	}
	nameBox.Draw(renderer)
}

type UISpellTrapCard struct {
	rect     sdl.FRect
	card     *data.SpellTrapCard
	quantity int
}

// GetCard implements UICard.
func (m *UISpellTrapCard) GetCard() data.Card { return m.card }

func (m *UISpellTrapCard) GetDescription() string { return m.card.GetDescription() }

func (m *UISpellTrapCard) GetImage() string { return m.card.GetImage() }

func (m *UISpellTrapCard) GetName() string { return m.card.GetName() }

func (m *UISpellTrapCard) GetRect() *sdl.FRect { return &m.rect }

func (m *UISpellTrapCard) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
	sdl.RenderFillRect(renderer, &m.rect)

	nameBox := TextBox{
		text:      m.card.GetName(),
		rect:      sdl.FRect{X: m.rect.X + 5, Y: m.rect.Y + 5, W: m.rect.W - 10, H: 20},
		color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
		font:      GetDefaultFont(20),
		textColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
	if m.quantity > 1 {
		qtyBox := TextBox{
			text:      strconv.Itoa(m.quantity),
			rect:      sdl.FRect{X: m.rect.X + 5, Y: m.rect.Y + 25, W: m.rect.W - 10, H: 20},
			color:     sdl.Color{R: 150, G: 150, B: 150, A: 150},
			font:      GetDefaultFont(20),
			textColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
		}
		qtyBox.Draw(renderer)
	}
	nameBox.Draw(renderer)
}

func CreateUICard(card data.Card, rect sdl.FRect, quantity int) UICard {
	if c, ok := card.(*data.MonsterCard); ok {
		return &UIMonsterCard{
			rect:     rect,
			card:     c,
			quantity: quantity,
		}
	}
	if c, ok := card.(*data.SpellTrapCard); ok {
		return &UISpellTrapCard{
			rect:     rect,
			card:     c,
			quantity: quantity,
		}
	}
	return nil
}
