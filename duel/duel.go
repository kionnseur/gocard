package duel

import (
	"fmt"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

// Renders the duel screen.
func RenderDuel(renderer *sdl.Renderer) *ui.AppState {
	elements := getDuelElements()
	buttons := make([]*ui.Button, 0, len(elements))
	for _, e := range elements {
		if btn, ok := e.(*ui.Button); ok {
			buttons = append(buttons, btn)
		}
	}

	duel := Duel{
		LeftPlayer: Player{
			Name:             "Red",
			LifePoints:       4000,
			InvocationPower:  9,
			InvocationNumber: 5,
			SpellTrapSet:     5,
			Deck:             40,
		},
		RightPlayer: Player{
			Name:             "Blue",
			LifePoints:       4000,
			InvocationPower:  9,
			InvocationNumber: 5,
			SpellTrapSet:     5,
			Deck:             40,
		},
		IsPaused: false,
		Timer:    100,
	}

	for {
		// Clear screen each frame
		sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
		sdl.RenderClear(renderer)

		// Draw all elements each frame
		for _, e := range elements {
			e.Draw(renderer)
		}

		// Handle events
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return &ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				x, y := event.Button().X, event.Button().Y
				for _, btn := range buttons {
					if ui.HitTest(btn.GetRect(), int32(x), int32(y)) {
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

// Creates the duel screen elements.
func getDuelElements() []ui.Element {
	player := Player{
		Name:             "Red",
		LifePoints:       4000,
		InvocationPower:  9,
		InvocationNumber: 5,
		SpellTrapSet:     5,
		Deck:             40,
	}

	font := ui.GetDefaultFont(24)
	leftPlayerHud := getLeftPlayerHud(font, player)
	rightPlayerHud := getRightPlayerHud(font, player)

	timer := getTimer(font)

	elements := []ui.Element{
		ui.NewButton("Back ⬅️", sdl.FRect{X: 140, Y: 280, W: 200, H: 50}, sdl.Color{R: 0, G: 255, B: 0, A: 255}, sdl.Color{R: 255, G: 0, B: 255, A: 255}, font, func() *ui.AppState { return &ui.AppState{State: ui.StateStartMenu} }),
	}

	elements = append(elements, &timer)
	elements = append(elements, leftPlayerHud...)
	elements = append(elements, rightPlayerHud...)

	return elements
}

// Creates the timer text box.
func getTimer(font *ttf.Font) ui.TextBox {
	countdown := 100

	return *ui.NewTextBox("Timer: "+fmt.Sprintf("%d", countdown),
		sdl.FRect{X: 540, Y: 0, W: 200, H: 50},
		sdl.Color{R: 255, G: 255, B: 0, A: 255},
		sdl.Color{R: 0, G: 0, B: 255, A: 255},
		font,
	)
}

// Creates the left player's HUD.
func getLeftPlayerHud(font *ttf.Font, player Player) []ui.Element {
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
		ui.NewHud(sdl.FRect{X: hudX, Y: hudY, W: hudW, H: hudH}, sdl.Color{R: 0, G: 120, B: 80, A: 255}),
		ui.NewTextBox(
			fmt.Sprintf("LP: %d", player.LifePoints),
			sdl.FRect{X: hudX + labelPadX + lpLabelW + 10, Y: hudY + labelPadY, W: lpValueW, H: lpH},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),
		ui.NewTextBox(
			player.Name,
			sdl.FRect{X: hudX + labelPadX, Y: hudY + lpH + labelPadY + 10, W: 100, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 200, G: 200, B: 200, A: 255},
			font,
		),

		ui.NewTextBox(
			fmt.Sprintf("IP: %d", player.InvocationPower),
			sdl.FRect{X: hudX + labelPadX + 120, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),

		ui.NewTextBox(
			fmt.Sprintf("D: %d", player.Deck), sdl.FRect{X: hudX + labelPadX + 170, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),

		ui.NewTextBox(
			fmt.Sprintf("S: %d", player.SpellTrapSet),
			sdl.FRect{X: hudX + labelPadX + 220, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),
	}
}

// Creates the right player's HUD.
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
		ui.NewHud(
			sdl.FRect{X: hudX, Y: hudY, W: hudW, H: hudH},
			sdl.Color{R: 0, G: 80, B: 120, A: 255},
		),
		ui.NewTextBox(
			fmt.Sprintf("LP: %d", player.LifePoints),
			sdl.FRect{X: hudX + labelPadX + lpLabelW + 10, Y: hudY + labelPadY, W: lpValueW, H: lpH},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),
		ui.NewTextBox(
			player.Name,
			sdl.FRect{X: hudX + labelPadX, Y: hudY + lpH + labelPadY + 10, W: 100, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 200, G: 200, B: 200, A: 255},
			font,
		),
		ui.NewTextBox(
			fmt.Sprintf("IP: %d", player.InvocationPower),
			sdl.FRect{X: hudX + labelPadX + 120, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),
		ui.NewTextBox(
			fmt.Sprintf("D: %d", player.Deck),
			sdl.FRect{X: hudX + labelPadX + 170, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),
		ui.NewTextBox(
			fmt.Sprintf("S: %d", player.SpellTrapSet),
			sdl.FRect{X: hudX + labelPadX + 220, Y: hudY + lpH + labelPadY + 10, W: 40, H: 28},
			sdl.Color{R: 0, G: 0, B: 0, A: 0},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
		),
	}
}

// Handles pause key press.
func (d *Duel) pausedPressed() {
	if d.IsPaused {
		d.IsPaused = false

	} else {
		d.IsPaused = true
	}
}
