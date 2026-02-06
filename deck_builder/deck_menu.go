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

	// Caching to avoid per-frame allocations
	cachedMenuButtons      []*ui.Button
	cachedMenuButtonsWidth float32
	lastDeckListCount      int
	cachedDeckListElements []ui.Element
	cachedDeckInfoDeckId   string
	cachedDeckInfoElements []ui.Element
	cachedDeckInfoButtons  []*ui.Button
)

// Renders the deck menu.
func RenderDeckMenu(renderer *sdl.Renderer, window *sdl.Window, appState *ui.AppState) *ui.AppState {

	for {
		// Update window size
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 35, 35, 45, 255)
		sdl.RenderClear(renderer)

		// Update selected deck if changed
		if askedDeckId != lastDeckId {
			deck = data.GetDeckById(askedDeckId)
			lastDeckId = askedDeckId
			// Invalidate deck info cache when selection changes
			cachedDeckInfoDeckId = ""
		}

		// Display deck info if selected (with caching)
		if askedDeckId != "" {
			if cachedDeckInfoDeckId != askedDeckId {
				cachedDeckInfoElements, cachedDeckInfoButtons = uiGetDeckInfo(deck, scrollableLVHDeckList.GetRect())
				cachedDeckInfoDeckId = askedDeckId
				// Update global vars used by handleBuilderButtonClick
				uiDeckInfo = cachedDeckInfoElements
				uiDeckInfoBtns = cachedDeckInfoButtons
			}
			for _, e := range cachedDeckInfoElements {
				e.Draw(renderer)
			}
			for _, e := range cachedDeckInfoButtons {
				e.Draw(renderer)
			}
		}

		// Update deck list only when deck count changes
		currentDeckCount := len(data.GetDeckList())
		if currentDeckCount != lastDeckListCount {
			updateDeckListElements()
			lastDeckListCount = currentDeckCount
		}

		// Display deck list, left column, and back button
		scrollableLVHDeckList.GetRect().H = float32(data.ScreenHeight) - gap - scrollableLVHDeckList.GetRect().Y

		scrollableLVHDeckList.Draw(renderer)

		// Cache menu buttons and only recreate if window width changes
		if cachedMenuButtons == nil || cachedMenuButtonsWidth != float32(data.ScreenWidth) {
			cachedMenuButtons = getDeckMenuButtons(scrollableLVHDeckList.GetRect())
			cachedMenuButtonsWidth = float32(data.ScreenWidth)
		}
		buttons = cachedMenuButtons
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
			sdl.Color{R: 60, G: 120, B: 80, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState {
				askedDeckId = ""
				return &ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": "", "action": "new"}}
			},
		),
		ui.NewButton(
			"Back ⬅️",
			sdl.FRect{X: float32(data.ScreenWidth) - 50, Y: 0, W: 50, H: 50},
			sdl.Color{R: 80, G: 80, B: 80, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateStartMenu} },
		),
	}
}

// Creates UI elements for the deck list.
func uiGetDeckListElements(decksList []data.Deck, parent *sdl.FRect) []ui.Element {
	elements := make([]ui.Element, len(decksList))
	for i, deck := range decksList {
		elements[i] = ui.NewButton(
			deck.GetName(),
			sdl.FRect{X: parent.X, Y: parent.Y + float32(i*35), W: parent.W, H: 30},
			sdl.Color{R: 60, G: 60, B: 70, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
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
		elements[i] = ui.NewTextBox(
			card.GetName(),
			sdl.FRect{X: offset, Y: float32((i + 5) * 35), W: 200, H: 30},
			sdl.Color{R: 50, G: 50, B: 60, A: 255},
			sdl.Color{R: 220, G: 220, B: 220, A: 255},
			font,
		)
	}
	buttons := []*ui.Button{
		ui.NewButton(
			"Edit",
			sdl.FRect{X: offset, Y: float32(8 * 35), W: 200, H: 30},
			sdl.Color{R: 60, G: 120, B: 80, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState {
				return &ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.GetId(), "action": "edit"}}
			},
		),

		ui.NewButton(
			"Duplicate",
			sdl.FRect{X: offset, Y: float32(9 * 35), W: 200, H: 30},
			sdl.Color{R: 70, G: 100, B: 150, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
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
			sdl.Color{R: 160, G: 60, B: 60, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
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
