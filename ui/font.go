package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

// caches fonts by size.
var fontCache = map[float32]*ttf.Font{}

// Draw a colored rectangle with centered text.
func drawTextBoxLike(renderer *sdl.Renderer, rect sdl.FRect, color sdl.Color, text string, textColor sdl.Color, font *ttf.Font) {
	// Draw background rectangle
	sdl.SetRenderDrawColor(renderer, color.R, color.G, color.B, color.A)
	sdl.RenderFillRect(renderer, &rect)
	sdl.SetRenderDrawColor(renderer, textColor.R, textColor.G, textColor.B, textColor.A)

	// Create text surface
	surface := ttf.RenderTextBlendedWrapped(font, text, uint64(len(text)), textColor, int32(rect.W))
	if surface == nil {
		return
	}
	defer sdl.DestroySurface(surface)

	// Create texture and render centered
	texture := sdl.CreateTextureFromSurface(renderer, surface)
	if texture == nil {
		return
	}
	defer sdl.DestroyTexture(texture)

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

// returns a cached font for the given size.
func GetDefaultFont(size float32) *ttf.Font {
	if ttf.WasInit() == 0 {
		ttf.Init()
	}
	size = getDefaultFontSize(size)
	if f, ok := fontCache[size]; ok && f != nil {
		return f
	}
	font := ttf.OpenFont("assets/fonts/D2CodingLigatureNerdFontMono-Regular.ttf", size)
	if font != nil {
		fontCache[size] = font
	}
	return font
}

// normalizes font size (defaults to 18 if < 2).
func getDefaultFontSize(size float32) float32 {
	if size < 2 {
		return 18
	}
	return size
}

// CleanupFontCache closes all cached fonts. Call on app exit.
func CleanupFontCache() {
	for _, font := range fontCache {
		if font != nil {
			ttf.CloseFont(font)
		}
	}
	fontCache = make(map[float32]*ttf.Font)
}
