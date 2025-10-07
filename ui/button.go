package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

type Button struct {
	text      string
	textColor sdl.Color
	font      *ttf.Font
	rect      sdl.FRect
	color     sdl.Color
	OnClick   func() *AppState
}

func (b *Button) Draw(renderer *sdl.Renderer) {
	drawTextBoxLike(renderer, b.rect, b.color, b.text, b.textColor, b.font)
}

func (e *Button) GetRect() *sdl.FRect {
	return &e.rect
}

func NewButton(text string, rect sdl.FRect, color, textColor sdl.Color, font *ttf.Font, onClick func() *AppState) *Button {

	return &Button{
		text:      text,
		rect:      rect,
		color:     color,
		textColor: textColor,
		font:      font,
		OnClick:   onClick,
	}
}
