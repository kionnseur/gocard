package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
)

type Hud struct {
	rect  sdl.FRect
	color sdl.Color
}

func NewHud(rect sdl.FRect, color sdl.Color) *Hud {
	return &Hud{
		rect:  rect,
		color: color,
	}
}

func (h *Hud) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, h.color.R, h.color.G, h.color.B, h.color.A)
	sdl.RenderFillRect(renderer, &h.rect)
}

func (e *Hud) GetRect() *sdl.FRect {
	return &e.rect
}
