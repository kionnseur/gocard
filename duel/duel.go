package duel

import (
	"fmt"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

func RenderDuel(renderer *sdl.Renderer) ui.AppState {
	elements := getDuelElements()
	buttons := make([]*ui.Button, 0, len(elements))
	for _, e := range elements {
		if btn, ok := e.(*ui.Button); ok {
			buttons = append(buttons, btn)
		}
	}

	duel := Duel{
		LeftPlayer: Player{
			Name:            "Red",
			LifePoints:      4000,
			InvocationPower: 9,
			IvocationNuber:  5,
			SpellTrapSet:    5,
			Deck:            40,
		},
		RightPlayer: Player{
			Name:            "Blue",
			LifePoints:      4000,
			InvocationPower: 9,
			IvocationNuber:  5,
			SpellTrapSet:    5,
			Deck:            40,
		},
		IsPaused: false,
		Timer:    100,
	}

	for {
		// Efface l'écran à chaque frame
		sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
		sdl.RenderClear(renderer)

		// Dessine tous les éléments à chaque frame
		for _, e := range elements {
			e.Draw(renderer)
		}

		// Gestion des événements
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				x, y := event.Button().X, event.Button().Y
				for _, btn := range buttons {
					if x > btn.Rect.X && x < btn.Rect.X+btn.Rect.W &&
						y > btn.Rect.Y && y < btn.Rect.Y+btn.Rect.H {
						return btn.OnClick()
					}
				}
			case sdl.EventType(sdl.KeycodeEscape):
				duel.pausedPressed()
			}
		}

		sdl.RenderPresent(renderer)
	}
}

func getDuelElements() []ui.Element {
	player := Player{
		Name:            "Red",
		LifePoints:      4000,
		InvocationPower: 9,
		IvocationNuber:  5,
		SpellTrapSet:    5,
		Deck:            40,
	}

	ttf.Init()

	font := ttf.OpenFont("assets/fonts/arial.ttf", 24)

	leftPlayerHud := getleftPlayerHud(font, player)
	rightPlayerHud := getRightPlayerHud(font, player)

	timer := getTimer(font)

	elements := []ui.Element{

		&ui.Button{
			Rect:      sdl.FRect{X: 140, Y: 280, W: 200, H: 50},
			Color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
			Text:      "Retour ⬅️",
			TextColor: sdl.Color{R: 255, G: 0, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.AppState{State: ui.StateStartMenu} },
		},
	}

	elements = append(elements, &timer)
	elements = append(elements, leftPlayerHud...)
	elements = append(elements, rightPlayerHud...)

	return elements
}

func getTimer(font *ttf.Font) ui.TextBox {
	countdown := 100

	return ui.TextBox{
		Rect:      sdl.FRect{X: 540, Y: 0, W: 200, H: 50},
		Color:     sdl.Color{R: 255, G: 255, B: 0, A: 255},
		Text:      "Timer: " + fmt.Sprintf("%d", countdown),
		TextColor: sdl.Color{R: 0, G: 0, B: 255, A: 255},
		Font:      font,
	}
}

func getleftPlayerHud(font *ttf.Font, player Player) []ui.Element {
	const (
		hudX      = 40
		hudY      = 0
		hudW      = 260
		hudH      = 100
		labelPadX = 10
		labelPadY = 10
		lpLabelW  = 50
		lpValueW  = 100
		lpH       = 40
	)

	return []ui.Element{
		&ui.Hud{
			Rect:  sdl.FRect{X: hudX, Y: hudY, W: hudW, H: hudH},
			Color: sdl.Color{R: 0, G: 120, B: 80, A: 255},
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX, Y: hudY + labelPadY, W: lpLabelW, H: lpH},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      "LP",
			TextColor: sdl.Color{R: 255, G: 215, B: 0, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + lpLabelW + 10, Y: hudY + labelPadY, W: lpValueW, H: lpH},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.LifePoints),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX, Y: hudY + lpH + labelPadY + 10, W: 100, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      player.Name,
			TextColor: sdl.Color{R: 200, G: 200, B: 200, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + 120, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.InvocationPower),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + 170, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.Deck),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + 220, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.SpellTrapSet),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
	}
}

func getRightPlayerHud(font *ttf.Font, player Player) []ui.Element {
	const (
		hudX      = 1020
		hudY      = 0
		hudW      = 260
		hudH      = 100
		labelPadX = 10
		labelPadY = 10
		lpLabelW  = 50
		lpValueW  = 100
		lpH       = 40
	)

	return []ui.Element{
		&ui.Hud{
			Rect:  sdl.FRect{X: hudX, Y: hudY, W: hudW, H: hudH},
			Color: sdl.Color{R: 0, G: 80, B: 120, A: 255},
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX, Y: hudY + labelPadY, W: lpLabelW, H: lpH},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      "LP",
			TextColor: sdl.Color{R: 255, G: 215, B: 0, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + lpLabelW + 10, Y: hudY + labelPadY, W: lpValueW, H: lpH},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.LifePoints),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX, Y: hudY + lpH + labelPadY + 10, W: 100, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      player.Name,
			TextColor: sdl.Color{R: 200, G: 200, B: 200, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + 120, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.InvocationPower),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + 170, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.Deck),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
		&ui.TextBox{
			Rect:      sdl.FRect{X: hudX + labelPadX + 220, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 0},
			Text:      fmt.Sprintf("%d", player.SpellTrapSet),
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
		},
	}
}

func (d *Duel) pausedPressed() {
	if d.IsPaused {
		d.IsPaused = false

	} else {
		d.IsPaused = true
	}
}
