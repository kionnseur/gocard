package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

type Element interface {
	Draw(renderer *sdl.Renderer)
}

func (h *Hud) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, h.Color.R, h.Color.G, h.Color.B, h.Color.A)
	sdl.RenderFillRect(renderer, &h.Rect)

}

func drawTextBoxLike(renderer *sdl.Renderer, rect sdl.FRect, color sdl.Color, text string, textColor sdl.Color, font *ttf.Font) {
	// Dessine le rectangle
	sdl.SetRenderDrawColor(renderer, color.R, color.G, color.B, color.A)
	sdl.RenderFillRect(renderer, &rect)
	sdl.SetRenderDrawColor(renderer, textColor.R, textColor.G, textColor.B, textColor.A)
	// Initialise TTF si ce n'est pas déjà fait
	if ttf.WasInit() == 0 {
		if ttf.Init() {
			return
		}
	}

	// Crée la surface du texte
	surface := ttf.RenderTextBlended(font, text, uint64(len(text)), textColor)
	if surface == nil {
		return
	}

	// Crée la texture à partir de la surface
	texture := sdl.CreateTextureFromSurface(renderer, surface)
	if texture == nil {
		return
	}
	defer sdl.DestroyTexture(texture)

	// Centre le texte dans le rectangle
	textW := surface.W
	textH := surface.H
	dstRect := sdl.FRect{
		X: rect.X + (rect.W-float32(textW))/2,
		Y: rect.Y + (rect.H-float32(textH))/2,
		W: float32(textW),
		H: float32(textH),
	}
	sdl.RenderTexture(renderer, texture, nil, &dstRect)
}

func (t *TextBox) Draw(renderer *sdl.Renderer) {
	drawTextBoxLike(renderer, t.Rect, t.Color, t.Text, t.TextColor, t.Font)
}

func (b *Button) Draw(renderer *sdl.Renderer) {
	drawTextBoxLike(renderer, b.Rect, b.Color, b.Text, b.TextColor, b.Font)
}
