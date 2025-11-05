package deck_builder

import (
	"gocard/data"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

const (
	gap      float32 = 20
	fontSize float32 = 24
)

var (
	deck                  *data.Deck
	lastDeckId            string
	font                  = ui.GetDefaultFont(24)
	scrollableLVHDeckList = ui.NewUIScrollableStackView(sdl.FRect{X: gap, Y: 80, W: 200, H: float32(data.ScreenHeight) - gap - 80}, sdl.Color{R: 100, G: 100, B: 100, A: 50}, 15)
	uiDeckListElements    []ui.Element
	askedDeckId           string

	buttons        []*ui.Button
	uiDeckInfoBtns []*ui.Button
	uiDeckInfo     []ui.Element
)

// Renders the deck menu.
func RenderDeckMenu(renderer *sdl.Renderer, window *sdl.Window, appState *ui.AppState) *ui.AppState {

	for {
		// Update window size
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 255, 165, 0, 255)
		sdl.RenderClear(renderer)

		// Update selected deck if changed
		if askedDeckId != lastDeckId {
			deck = data.GetDeckById(askedDeckId)
			lastDeckId = askedDeckId
		}

		// Display deck info if selected
		if askedDeckId != "" {
			uiDeckInfo, uiDeckInfoBtns = uiGetDeckInfo(deck, scrollableLVHDeckList.GetRect())
			for _, e := range uiDeckInfo {
				e.Draw(renderer)
			}
			for _, e := range uiDeckInfoBtns {
				e.Draw(renderer)
			}
		}
		// Update deck list
		updateDeckListElements()

		// Display deck list, left column, and back button
		scrollableLVHDeckList.GetRect().H = float32(data.ScreenHeight) - gap - scrollableLVHDeckList.GetRect().Y

		scrollableLVHDeckList.Draw(renderer)

		buttons = getDeckMenuButtons(scrollableLVHDeckList.GetRect())
		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		if as := handleEventsDeckMenu(); as != nil {
			return as
		}
		sdl.RenderPresent(renderer)
	}
}

// Handles events for the deck menu.
func handleEventsDeckMenu() *ui.AppState {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		switch event.Type() {
		case sdl.EventQuit:
			return &ui.AppState{State: ui.StateQuit}
		case sdl.EventMouseButtonDown:
			return handleBuilderButtonClick(&event)

		case sdl.EventMouseWheel:
			if scrollableLVHDeckList.OnScroll != nil {
				scrollableLVHDeckList.OnScroll(&event)
			}
		}
	}
	return nil
}

// Handles button and element clicks.
func handleBuilderButtonClick(event *sdl.Event) *ui.AppState {
	x, y := event.Button().X, event.Button().Y

	if as := hitTestButtons(buttons, int32(x), int32(y)); as != nil {
		return as
	}

	if as := hitTestDeckListElements(uiDeckListElements, int32(x), int32(y), float32(y)); as != nil {
		return as
	}

	if as := hitTestButtons(uiDeckInfoBtns, int32(x), int32(y)); as != nil {
		return as
	}

	return nil
}

// Checks clicks on a button list.
func hitTestButtons(btns []*ui.Button, x, y int32) *ui.AppState {
	for _, btn := range btns {
		if ui.HitTest(btn.GetRect(), x, y) {
			return btn.OnClick()
		}
	}
	return nil
}

// Checks clicks on deck list elements (with scrolling).
func hitTestDeckListElements(elements []ui.Element, x, y int32, fy float32) *ui.AppState {
	for _, e := range elements {
		if btn, ok := e.(*ui.Button); ok {
			rect := *btn.GetRect()
			if fy > scrollableLVHDeckList.GetRect().Y && fy < scrollableLVHDeckList.GetRect().Y+scrollableLVHDeckList.GetRect().H {
				rect.Y = rect.Y - scrollableLVHDeckList.GetScrollY()
				if ui.HitTest(&rect, x, y) {
					return btn.OnClick()
				}
			}
		}
	}
	return nil
}

// Creates the deck menu buttons.
func getDeckMenuButtons(parent *sdl.FRect) []*ui.Button {
	return []*ui.Button{
		ui.NewButton(
			"New Deck",
			sdl.FRect{X: parent.X, Y: parent.Y - 30, W: 200, H: 30},
			sdl.Color{R: 0, G: 0, B: 0, A: 100},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
			func() *ui.AppState {
				askedDeckId = ""
				return &ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": "", "action": "new"}}
			},
		),
		ui.NewButton(
			"Back ⬅️",
			sdl.FRect{X: float32(data.ScreenWidth) - 50, Y: 0, W: 50, H: 50},
			sdl.Color{R: 0, G: 255, B: 0, A: 255},
			sdl.Color{R: 255, G: 0, B: 255, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateStartMenu} },
		),
	}
}

// Creates UI elements for the deck list.
func uiGetDeckListElements(decksList []data.Deck, parent *sdl.FRect) []ui.Element {
	elements := make([]ui.Element, len(decksList))
	for i, deck := range decksList {
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = ui.NewButton(
			deck.GetName(),
			sdl.FRect{X: parent.X, Y: parent.Y + float32(i*35), W: parent.W, H: 30},
			sdl.Color{R: r, G: g, B: b, A: 255},
			sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			font,
			func() *ui.AppState {
				askedDeckId = deck.GetId()
				return nil
			},
		)
	}
	return elements
}

// Creates UI for deck info.
func uiGetDeckInfo(deck *data.Deck, parent *sdl.FRect) ([]ui.Element, []*ui.Button) {

	offset := parent.X + gap + parent.W

	// Show the first 3 card names
	elements := make([]ui.Element, min(3, len(deck.GetCards())))

	for i, card := range deck.GetCards() {
		if i > 2 {
			break
		}
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = ui.NewTextBox(
			card.GetName(),
			sdl.FRect{X: offset, Y: float32((i + 5) * 35), W: 200, H: 30},
			sdl.Color{R: r, G: g, B: b, A: 255},
			sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			font,
		)
	}
	buttons := []*ui.Button{
		ui.NewButton(
			"Edit",
			sdl.FRect{X: offset, Y: float32(8 * 35), W: 200, H: 30},
			sdl.Color{R: 100, G: 200, B: 100, A: 255},
			sdl.Color{R: 155, G: 55, B: 155, A: 255},
			font,
			func() *ui.AppState {
				return &ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.GetId(), "action": "edit"}}
			},
		),

		ui.NewButton(
			"Duplicate",
			sdl.FRect{X: offset, Y: float32(9 * 35), W: 200, H: 30},
			sdl.Color{R: 100, G: 100, B: 200, A: 255},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			font,
			func() *ui.AppState {
				data.DuplicateDeckById(deck.GetId())
				updateDeckListElements()
				return nil
			},
		),
		ui.NewButton(
			"Delete",
			sdl.FRect{X: offset, Y: float32(10 * 35), W: parent.W, H: 30},
			sdl.Color{R: 200, G: 100, B: 100, A: 255},
			sdl.Color{R: 55, G: 155, B: 155, A: 255},
			font,
			func() *ui.AppState {
				data.DeleteDeckById(deck.GetId())
				return nil
			},
		),
	}

	return elements, buttons

}

// Updates the deck list elements.
func updateDeckListElements() {
	uiDeckListElements = uiGetDeckListElements(data.GetDeckList(), scrollableLVHDeckList.GetRect())
	scrollableLVHDeckList.SetElements(uiDeckListElements)
	scrollableLVHDeckList.OnScroll = func(event *sdl.Event) {
		y := event.Wheel().Y
		scrollableLVHDeckList.SetScrollY(scrollableLVHDeckList.GetScrollY() - float32(y)*gap)
		if scrollableLVHDeckList.GetScrollY() < 0 {
			scrollableLVHDeckList.SetScrollY(0)
		}
		maxScroll := float32(len(uiDeckListElements))*35 - scrollableLVHDeckList.GetRect().H
		if maxScroll < 0 {
			maxScroll = 0
		}
		if scrollableLVHDeckList.GetScrollY() > maxScroll {
			scrollableLVHDeckList.SetScrollY(maxScroll)
		}
	}
}
