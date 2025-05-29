package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

var fontCache = map[float32]*ttf.Font{}

func drawTextBoxLike(renderer *sdl.Renderer, rect sdl.FRect, color sdl.Color, text string, textColor sdl.Color, font *ttf.Font) {
	// Dessine le rectangle
	sdl.SetRenderDrawColor(renderer, color.R, color.G, color.B, color.A)
	sdl.RenderFillRect(renderer, &rect)
	sdl.SetRenderDrawColor(renderer, textColor.R, textColor.G, textColor.B, textColor.A)
	// Initialise TTF si ce n'est pas déjà fait

	// Crée la surface du texte
	surface := ttf.RenderTextBlendedWrapped(font, text, uint64(len(text)), textColor, int32(rect.W))
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

func GetDefaultFont(size float32) *ttf.Font {
	// stock une font par taille passé en parametre
	if ttf.WasInit() == 0 {
		ttf.Init()
	}
	size = GetDefaultFontSize(size)
	if f, ok := fontCache[size]; ok && f != nil {
		return f
	}
	font := ttf.OpenFont("assets/fonts/D2CodingLigatureNerdFontMono-Regular.ttf", size)
	if font != nil {
		fontCache[size] = font
	}
	return font
}

func GetDefaultFontSize(size float32) float32 {
	// pour que 1 et 0 soit une valeur par defaut
	if size < 2 {
		return 18
	}
	return size
}

// /////////////
// HUD
// /////////////
func (h *Hud) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, h.Color.R, h.Color.G, h.Color.B, h.Color.A)
	sdl.RenderFillRect(renderer, &h.Rect)
}
func (e *Hud) GetRect() *sdl.FRect {
	return &e.Rect
}

// /////////////
// TextBox
// /////////////
func (t *TextBox) Draw(renderer *sdl.Renderer) {
	drawTextBoxLike(renderer, t.Rect, t.Color, t.Text, t.TextColor, t.Font)
}

func (e *TextBox) GetRect() *sdl.FRect {
	return &e.Rect
}

// /////////////
// Button
// /////////////
func (b *Button) Draw(renderer *sdl.Renderer) {
	drawTextBoxLike(renderer, b.Rect, b.Color, b.Text, b.TextColor, b.Font)
}

func (e *Button) GetRect() *sdl.FRect {
	return &e.Rect
}

// /////////////
// UIScrollableStackView
// /////////////
func (e *UIScrollableStackView) GetRect() *sdl.FRect {
	return &e.Rect
}

func (e *UIScrollableStackView) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, e.Color.R, e.Color.G, e.Color.B, e.Color.A)
	sdl.RenderFillRect(renderer, &e.Rect)
	for _, elem := range e.Elements {
		rect := elem.GetRect()
		tmpRect := *rect
		tmpRect.Y -= e.ScrollY
		// Test de visibilité
		if tmpRect.Y+tmpRect.H > e.Rect.Y && tmpRect.Y < e.Rect.Y+e.Rect.H {
			original := *rect
			*rect = tmpRect
			elem.Draw(renderer)
			*rect = original
		}
	}
}

// /////////////
// UIScrollableGridView
// /////////////

func (e *UIScrollableGridView) GetRect() *sdl.FRect {
	return &e.Rect
}
func (e *UIScrollableGridView) Draw(renderer *sdl.Renderer) {
	e.SetElementsPosition(&e.Rect)

	sdl.SetRenderDrawColor(renderer, e.Color.R, e.Color.G, e.Color.B, e.Color.A)
	sdl.RenderFillRect(renderer, &e.Rect)

	for _, elem := range e.Elements {
		rect := elem.GetRect()
		rect.Y -= e.ScrollY

		if rect.Y+rect.H > e.Rect.Y && rect.Y < e.Rect.Y+e.Rect.H {
			elem.Draw(renderer)
		}
	}
}

func (e *UIScrollableGridView) SetElementsPosition(parent *sdl.FRect) {
	cfg := e.GridConfig
	maxColCards := int((e.Rect.W + cfg.CardGap) / (cfg.CardWidth + cfg.CardGap))
	startX := e.Rect.X + (e.Rect.W-float32(maxColCards)*(cfg.CardWidth+cfg.CardGap)+cfg.CardGap)/2
	y := e.Rect.Y + cfg.CardGap
	for i, elem := range e.Elements {
		x := startX + float32(i%maxColCards)*(cfg.CardWidth+cfg.CardGap)
		if i%maxColCards == 0 && i != 0 {
			y += cfg.CardHeight + cfg.CardGap
		}
		rect := elem.GetRect()
		rect.X = x
		rect.Y = y
		rect.W = cfg.CardWidth
		rect.H = cfg.CardHeight
	}
}
