package deck_builder

import (
	"gocard/data"
	"gocard/ui"
	"slices"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

var (
	selectedCard   data.Card
	selectedDeck   *data.Deck
	dollyDeck      data.Deck // Copy of deck being edited
	playerCardDict map[int]int

	deckId string

	cardWidth   float32 = 100
	cardHeight  float32 = 150
	cardGap     float32 = 10
	maxColCards int

	uiElements               []ui.Element
	scrollableLVHRightColumn ui.UIScrollableGridView
	uiCenterColumn           []ui.UICard
	returnBtns               []*ui.Button
	playerUICards            []ui.UICard
)

// Renders the deck editor.
func RenderDeckEditor(renderer *sdl.Renderer, window *sdl.Window, appState *ui.AppState) *ui.AppState {
	scrollableLVHRightColumn = *(ui.NewUIScrollableGridView(renderer, sdl.FRect{}, sdl.Color{R: 100, G: 100, B: 100, A: 50}, 3, *ui.NewGridConfig(cardWidth, cardHeight, cardGap)))

	deckId := appState.Data["deckId"]
	selectedDeck = data.GetDeckById(deckId)
	dollyDeck = data.CloneDeckById(deckId)

	playerCardDict := data.GetPlayerCards()
	playerUICards = getPlayerCardListUI(playerCardDict)
	slices.SortFunc(playerUICards, func(a, b ui.UICard) int {
		return a.GetCard().GetId() - b.GetCard().GetId()
	})

	for {
		// Update window size
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 0, 165, 0, 100)
		sdl.RenderClear(renderer)

		uiElements = append(getColumnUI(), getLeftColumnUI()...)

		uiCenterColumn = getDeckCardListUI(uiElements[1].GetRect())

		// Update right scrollable column
		setScrollableLVHRightColumn(&playerUICards)

		// Draw UI elements, deck cards, buttons, scrollable views
		for _, e := range uiElements {
			e.Draw(renderer)
		}
		for _, e := range uiCenterColumn {
			e.Draw(renderer)
		}
		returnBtns = getDeckEditorButtons()
		for _, btn := range returnBtns {
			btn.Draw(renderer)
		}
		scrollableLVHRightColumn.Draw(renderer)

		// Handle events
		if state := handleEvents(); state != nil {
			return state
		}

		sdl.RenderPresent(renderer)
	}
}

// Handles events for the deck editor.
func handleEvents() *ui.AppState {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		switch event.Type() {
		case sdl.EventQuit:
			return &ui.AppState{State: ui.StateQuit}
		case sdl.EventMouseButtonDown:
			if state := handleButtonClick(event); state != nil {
				return state
			}
		case sdl.EventMouseWheel:
			scrollableLVHRightColumn.OnScroll(&event)
		}
	}
	return nil
}

// Handles button and card clicks.
func handleButtonClick(event sdl.Event) *ui.AppState {
	x, y := event.Button().X, event.Button().Y
	// Check left column buttons
	for _, btn := range uiElements {
		if btn, ok := btn.(*ui.Button); ok && ui.HitTest(btn.GetRect(), int32(x), int32(y)) {
			return btn.OnClick()
		}
	}
	// If click on card, select it
	for _, uiCard := range uiCenterColumn {
		if ui.HitTest(uiCard.GetRect(), int32(x), int32(y)) {
			selectedCard = uiCard.GetCard()
			return nil
		}
	}
	// Check right column cards
	for _, uiCard := range playerUICards {
		if y > scrollableLVHRightColumn.GetRect().Y && y < scrollableLVHRightColumn.GetRect().Y+scrollableLVHRightColumn.GetRect().H &&
			ui.HitTest(uiCard.GetRect(), int32(x), int32(y)) {
			selectedCard = uiCard.GetCard()
			return nil
		}
	}
	// Check top right buttons
	for _, btn := range returnBtns {
		if ui.HitTest(btn.GetRect(), int32(x), int32(y)) {
			return btn.OnClick()
		}
	}
	return nil
}

// Returns the three column rectangles.
func getColumnUI() []ui.Element {
	gap := float32(data.ScreenWidth / 48.0)
	widthColA := float32(data.ScreenWidth * 103 / 480)
	widthColB := float32(data.ScreenWidth * 197 / 480)
	widthColC := float32(data.ScreenWidth * 7 / 24)

	return []ui.Element{
		ui.NewHud(sdl.FRect{X: gap, Y: 0, W: widthColA, H: float32(data.ScreenHeight)}, sdl.Color{R: 255, G: 165, B: 0, A: 255}),
		ui.NewHud(sdl.FRect{X: widthColA + (2 * gap), Y: 0, W: widthColB, H: float32(data.ScreenHeight)}, sdl.Color{R: 0, G: 255, B: 165, A: 255}),
		ui.NewHud(sdl.FRect{X: widthColA + widthColB + 3*gap, Y: 0, W: widthColC, H: float32(data.ScreenHeight)}, sdl.Color{R: 165, G: 0, B: 255, A: 255}),
	}
}

// Configures the right scrollable column.
func setScrollableLVHRightColumn(playerCards *[]ui.UICard) {
	length := len(*playerCards)

	scrollableLVHRightColumn.SetElements(make([]ui.Element, length))
	for i, e := range *playerCards {
		scrollableLVHRightColumn.GetElements()[i] = e
	}

	rec := uiElements[2].GetRect()
	scrollYOrigin := rec.Y + 3*gap
	height := float32(data.ScreenHeight) - scrollYOrigin - gap
	if height < 0 {
		height = 0
	}

	scrollableLVHRightColumn.SetRect(sdl.FRect{
		X: rec.X + gap/2,
		Y: scrollYOrigin,
		W: uiElements[2].GetRect().W - gap,
		H: height,
	})

	// Configure scroll handler after rect is set
	scrollableLVHRightColumn.OnScroll = func(event *sdl.Event) {
		y := event.Wheel().Y
		scrollableLVHRightColumn.SetScrollY(scrollableLVHRightColumn.GetScrollY() - (float32(y)*cardGap)*scrollableLVHRightColumn.GetScrollSpeed())
		if scrollableLVHRightColumn.GetScrollY() < 0 {
			scrollableLVHRightColumn.SetScrollY(0)
		}
		numRows := (length + maxColCards - 1) / maxColCards
		maxScroll := float32(numRows)*(cardHeight+cardGap) - scrollableLVHRightColumn.GetRect().H
		if maxScroll < 0 {
			maxScroll = 0
		}
		if scrollableLVHRightColumn.GetScrollY() > maxScroll {
			scrollableLVHRightColumn.SetScrollY(maxScroll)
		}
	}
}

// Returns UI for selected card and add/remove buttons.
func getLeftColumnUI() []ui.Element {
	if selectedCard == nil {
		return nil
	}

	elements := make([]ui.Element, 3)

	gap := float32(data.ScreenWidth / 48.0)
	cardWidth := (float32(data.ScreenWidth)*103.0)/480.0 - 2*gap
	cardHeight := cardWidth * 1.5
	cardRect := sdl.FRect{X: float32(2 * data.ScreenWidth / 48), Y: 40, W: cardWidth, H: cardHeight}

	uiCard := ui.CreateUICard(selectedCard, cardRect, dollyDeck.CountCard(selectedCard))
	elements[0] = uiCard

	currentDeckSize := len(dollyDeck.GetCards())
	currentCardCountInDeck := dollyDeck.CountCard(selectedCard)
	currentPlayerCardCount := playerCardDict[selectedCard.GetId()]

	// Add button (disabled if limits reached)
	if currentDeckSize >= 40 || currentCardCountInDeck >= 3 || currentPlayerCardCount > 3-currentCardCountInDeck {
		elements[1] = ui.NewTextBox("Add to deck", sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30}, sdl.Color{R: 80, G: 80, B: 80, A: 255}, sdl.Color{R: 255, G: 0, B: 0, A: 255}, ui.GetDefaultFont(20))
	} else {
		elements[1] = ui.NewButton(
			"Add to deck",
			sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30},
			sdl.Color{R: 20, G: 20, B: 20, A: 100},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			ui.GetDefaultFont(20),
			func() *ui.AppState {
				dollyDeck.SetCards(append(dollyDeck.GetCards(), selectedCard))
				return nil
			},
		)
	}
	// Remove button (disabled if none in deck)
	if currentCardCountInDeck == 0 {
		elements[2] = ui.NewTextBox("Remove from deck", sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30}, sdl.Color{R: 80, G: 80, B: 80, A: 255}, sdl.Color{R: 255, G: 0, B: 0, A: 255}, ui.GetDefaultFont(20))
	} else {
		elements[2] = ui.NewButton(
			"Remove from deck",
			sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30},
			sdl.Color{R: 20, G: 20, B: 20, A: 100},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			ui.GetDefaultFont(20),
			func() *ui.AppState {
				dollyDeck.RemoveCard(selectedCard)
				return nil
			},
		)
	}
	return elements
}

// Returns UI cards for the deck in center column.
func getDeckCardListUI(centerColRect *sdl.FRect) []ui.UICard {
	uiCenterColumn := make([]ui.UICard, 0, len(dollyDeck.GetCards()))

	// Calculate cards that fit in column width
	maxColCards = int((centerColRect.W + cardGap) / (cardWidth + cardGap))

	startX := centerColRect.X + (centerColRect.W-float32(maxColCards)*(cardWidth+cardGap)+cardGap)/2
	y := 2 * cardGap

	for i, card := range dollyDeck.GetCards() {

		x := startX + float32(i%maxColCards)*(cardWidth+cardGap)
		if i%maxColCards == 0 && i != 0 {
			y += cardHeight + cardGap
		}
		cardRect := sdl.FRect{X: x, Y: y, W: cardWidth, H: cardHeight}
		uiCard := ui.CreateUICard(card, cardRect, deck.CountCard(card))
		uiCenterColumn = append(uiCenterColumn, uiCard)
	}
	return uiCenterColumn
}

// Returns UI cards for player's collection.
func getPlayerCardListUI(playerCardDict map[int]int) []ui.UICard {
	elements := make([]ui.UICard, 0, len(playerCardDict))
	for cardId, qty := range playerCardDict {
		card := data.GetAllCards()[cardId]
		elements = append(elements, ui.CreateUICard(card, sdl.FRect{}, qty))
	}
	return elements
}

// Returns the deck editor buttons.
func getDeckEditorButtons() []*ui.Button {
	return []*ui.Button{
		ui.NewButton(
			"Back",
			sdl.FRect{X: float32(data.ScreenWidth) - 50, Y: 0, W: 50, H: 50},
			sdl.Color{R: 0, G: 255, B: 0, A: 255},
			sdl.Color{R: 255, G: 0, B: 255, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateDeckMenu} },
		),
		ui.NewButton(
			"Save✅✔️☑️",
			sdl.FRect{X: float32(data.ScreenWidth) - 100, Y: 0, W: 50, H: 50},
			sdl.Color{R: 0, G: 255, B: 0, A: 255},
			sdl.Color{R: 255, G: 0, B: 255, A: 255},
			font,
			func() *ui.AppState {
				if deckId != "" {
					*selectedDeck = dollyDeck
				}
				data.SaveDeck(dollyDeck)
				return nil
			},
		),
	}
}
