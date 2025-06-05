package ui

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
)

type UIScrollableStackView struct {
	rect        sdl.FRect
	color       sdl.Color
	scrollY     float32
	scrollSpeed float32
	elements    []Element
	OnScroll    func(event *sdl.Event)
}

type UIScrollableGridView struct {
	rect        sdl.FRect
	color       sdl.Color
	scrollY     float32
	scrollSpeed float32
	elements    []Element
	OnScroll    func(event *sdl.Event)
	gridConfig  *GridConfig
}

type GridConfig struct {
	cardWidth  float32
	cardHeight float32
	cardGap    float32
}

func NewGridConfig(cardWidth, cardHeight, cardGap float32) *GridConfig {
	return &GridConfig{
		cardWidth:  cardWidth,
		cardHeight: cardHeight,
		cardGap:    cardGap,
	}
}

// /////////////
// UIScrollableStackView
// /////////////
func NewUIScrollableStackView(rect sdl.FRect, color sdl.Color, scrollSpeed float32) *UIScrollableStackView {
	return &UIScrollableStackView{
		rect:        rect,
		color:       color,
		scrollY:     0,
		scrollSpeed: scrollSpeed,
		elements:    make([]Element, 0),
		OnScroll:    nil,
	}
}

func (e *UIScrollableStackView) GetRect() *sdl.FRect {
	return &e.rect
}

func (e *UIScrollableStackView) SetRect(rect sdl.FRect) {
	e.rect = rect
}

func (e *UIScrollableStackView) GetElements() []Element {
	return e.elements
}
func (e *UIScrollableStackView) SetElements(elements []Element) {
	e.elements = elements
}
func (e *UIScrollableStackView) GetScrollY() float32 {
	return e.scrollY
}
func (e *UIScrollableStackView) SetScrollY(scrollY float32) {
	e.scrollY = scrollY
}

func (e *UIScrollableStackView) GetScrollSpeed() float32 {
	return e.scrollSpeed
}

func (e *UIScrollableStackView) Draw(renderer *sdl.Renderer) {
	sdl.SetRenderDrawColor(renderer, e.color.R, e.color.G, e.color.B, e.color.A)
	sdl.RenderFillRect(renderer, &e.rect)
	for _, elem := range e.elements {
		rect := elem.GetRect()
		tmpRect := *rect
		tmpRect.Y -= e.scrollY
		// Test de visibilité
		if tmpRect.Y+tmpRect.H > e.rect.Y && tmpRect.Y < e.rect.Y+e.rect.H {
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
func NewUIScrollableGridView(renderer *sdl.Renderer, rect sdl.FRect, color sdl.Color, scrollSpeed float32, gridConfig GridConfig) *UIScrollableGridView {
	return &UIScrollableGridView{
		rect:        rect,
		color:       color,
		scrollY:     0,
		scrollSpeed: scrollSpeed,
		elements:    make([]Element, 0),
		gridConfig:  &gridConfig,
	}
}

func (e *UIScrollableGridView) GetRect() *sdl.FRect {
	return &e.rect
}

func (e *UIScrollableGridView) SetRect(rect sdl.FRect) {
	e.rect = rect
}

func (e *UIScrollableGridView) SetScrollY(scrollY float32) {
	e.scrollY = scrollY
}

func (e *UIScrollableGridView) GetElements() []Element {
	return e.elements
}

func (e *UIScrollableGridView) SetElements(elements []Element) {
	e.elements = elements
}

func (e *UIScrollableGridView) GetScrollY() float32 {
	return e.scrollY
}

func (e *UIScrollableGridView) GetScrollSpeed() float32 {
	return e.scrollSpeed
}

func (e *UIScrollableGridView) Draw(renderer *sdl.Renderer) {
	e.SetElementsPosition(&e.rect)

	sdl.SetRenderDrawColor(renderer, e.color.R, e.color.G, e.color.B, e.color.A)
	sdl.RenderFillRect(renderer, &e.rect)

	for _, elem := range e.elements {
		rect := elem.GetRect()
		rect.Y -= e.scrollY

		if rect.Y+rect.H > e.rect.Y && rect.Y < e.rect.Y+e.rect.H {
			elem.Draw(renderer)
		}
	}
}

func (e *UIScrollableGridView) SetElementsPosition(parent *sdl.FRect) {
	cfg := e.gridConfig
	maxColCards := int((e.rect.W + cfg.cardGap) / (cfg.cardWidth + cfg.cardGap))
	startX := e.rect.X + (e.rect.W-float32(maxColCards)*(cfg.cardWidth+cfg.cardGap)+cfg.cardGap)/2
	y := e.rect.Y + cfg.cardGap
	for i, elem := range e.elements {
		x := startX + float32(i%maxColCards)*(cfg.cardWidth+cfg.cardGap)
		if i%maxColCards == 0 && i != 0 {
			y += cfg.cardHeight + cfg.cardGap
		}
		rect := elem.GetRect()
		rect.X = x
		rect.Y = y
		rect.W = cfg.cardWidth
		rect.H = cfg.cardHeight
	}
}
