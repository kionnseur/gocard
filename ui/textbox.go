package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

type TextBox struct {
	rect      sdl.FRect
	color     sdl.Color
	text      string
	textColor sdl.Color
	font      *ttf.Font
}

func NewTextBox(text string, rect sdl.FRect, color, textColor sdl.Color, font *ttf.Font) *TextBox {
	return &TextBox{
		rect:      rect,
		color:     color,
		text:      text,
		textColor: textColor,
		font:      font,
	}
}

func (t *TextBox) Draw(renderer *sdl.Renderer) {
	drawTextBoxLike(renderer, t.rect, t.color, t.text, t.textColor, t.font)
}

func (e *TextBox) GetRect() *sdl.FRect {
	return &e.rect
}
