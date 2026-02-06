package start_menu

import (
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

// Displays the start menu and handles events.
func RenderStartMenu(renderer *sdl.Renderer) *ui.AppState {
	buttons := getStartMenuButtons()

	for {
		// Handle events
		if state := handleEvents(buttons); state != nil {
			return state
		}

		// Clear screen
		sdl.SetRenderDrawColor(renderer, 35, 35, 45, 255)
		sdl.RenderClear(renderer)

		// Draw buttons
		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		sdl.RenderPresent(renderer)
	}
}

// Handles events for the start menu.
func handleEvents(buttons []*ui.Button) *ui.AppState {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		switch event.Type() {
		case sdl.EventQuit:
			return &ui.AppState{State: ui.StateQuit}
		case sdl.EventMouseButtonDown:
			if state := handleButtonClick(event, buttons); state != nil {
				return state
			}
		}
	}
	return nil
}

// Handles button clicks.
func handleButtonClick(event sdl.Event, buttons []*ui.Button) *ui.AppState {
	x, y := event.Button().X, event.Button().Y
	for _, btn := range buttons {
		rect := btn.GetRect()
		if x > rect.X && x < rect.X+rect.W &&
			y > rect.Y && y < rect.Y+rect.H {
			if as := btn.OnClick(); as != nil {
				return as
			}
		}
	}
	return nil
}

// Creates the start menu buttons.
func getStartMenuButtons() []*ui.Button {
	font := ui.GetDefaultFont(24)

	return []*ui.Button{
		ui.NewButton(
			"Deck Builder",
			sdl.FRect{X: 140.0, Y: 80.0, W: 200.0, H: 50.0},
			sdl.Color{R: 45, G: 85, B: 135, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateDeckMenu} },
		),
		ui.NewButton(
			"Duel",
			sdl.FRect{X: 140, Y: 180, W: 200, H: 50},
			sdl.Color{R: 135, G: 45, B: 45, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateDuel} },
		),
	}
}
