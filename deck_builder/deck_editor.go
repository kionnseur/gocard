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
	editingDeck    data.Deck // Working copy of deck being edited
	playerCardDict map[int]int

	deckId string

	cardWidth  float32 = 100
	cardHeight float32 = 150
	cardGap    float32 = 10

	uiElements               []ui.Element
	scrollableLVHRightColumn ui.UIScrollableGridView
	scrollableCenterColumn   ui.UIScrollableGridView
	returnBtns               []*ui.Button
	playerUICards            []ui.UICard

	// Caching to avoid per-frame allocations
	cachedEditorButtons        []*ui.Button
	cachedEditorButtonsWidth   float32
	cachedLeftColumnCardId     int
	cachedLeftColumnCardCount  int
	cachedLeftColumnElements   []ui.Element
	lastEditingDeckSize        int
	cachedCenterColumnElements []ui.Element
)

// Renders the deck editor.
func RenderDeckEditor(renderer *sdl.Renderer, window *sdl.Window, appState *ui.AppState) *ui.AppState {
	scrollableLVHRightColumn = *(ui.NewUIScrollableGridView(renderer, sdl.FRect{}, sdl.Color{R: 100, G: 100, B: 100, A: 50}, 3, *ui.NewGridConfig(cardWidth, cardHeight, cardGap)))
	scrollableCenterColumn = *(ui.NewUIScrollableGridView(renderer, sdl.FRect{}, sdl.Color{R: 45, G: 45, B: 55, A: 255}, 3, *ui.NewGridConfig(cardWidth, cardHeight, cardGap)))

	deckId = appState.Data["deckId"]
	selectedDeck = data.GetDeckById(deckId)
	editingDeck = data.CloneDeckById(deckId)
	// Force refresh of center column on deck open
	lastEditingDeckSize = -1

	playerCardDict := data.GetPlayerCards()
	playerUICards = getPlayerCardListUI(playerCardDict)
	slices.SortFunc(playerUICards, func(a, b ui.UICard) int {
		return a.GetCard().GetId() - b.GetCard().GetId()
	})

	for {
		// Update window size
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 30, 30, 40, 255)
		sdl.RenderClear(renderer)

		// Build column UI (columns change with window size)
		columnUI := getColumnUI()

		// Cache left column UI - rebuild if selected card changes OR deck content changes
		if selectedCard != nil {
			currentCardCount := editingDeck.CountCard(selectedCard)
			if selectedCard.GetId() != cachedLeftColumnCardId || currentCardCount != cachedLeftColumnCardCount {
				cachedLeftColumnElements = getLeftColumnUI()
				cachedLeftColumnCardId = selectedCard.GetId()
				cachedLeftColumnCardCount = currentCardCount
			}
		} else if selectedCard == nil {
			cachedLeftColumnElements = nil
			cachedLeftColumnCardId = 0
			cachedLeftColumnCardCount = 0
		}

		uiElements = uiElements[:0]
		uiElements = append(columnUI, cachedLeftColumnElements...)

		// Update center scrollable column only when deck size changes
		currentDeckSize := len(editingDeck.GetCards())
		if currentDeckSize != lastEditingDeckSize {
			setScrollableCenterColumn(uiElements[1].GetRect())
			lastEditingDeckSize = currentDeckSize
		}

		// Update right scrollable column
		setScrollableLVHRightColumn(&playerUICards)

		// Draw UI elements, deck cards, buttons, scrollable views
		for _, e := range uiElements {
			e.Draw(renderer)
		}
		scrollableCenterColumn.Draw(renderer)

		// Cache editor buttons - only recreate when window width changes
		if cachedEditorButtons == nil || cachedEditorButtonsWidth != float32(data.ScreenWidth) {
			cachedEditorButtons = getDeckEditorButtons()
			cachedEditorButtonsWidth = float32(data.ScreenWidth)
		}
		returnBtns = cachedEditorButtons
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
			// Get mouse position to determine which column to scroll
			var mx, my float32
			sdl.GetMouseState(&mx, &my)

			// Check if mouse is over center column
			centerRect := scrollableCenterColumn.GetRect()
			if mx >= centerRect.X && mx <= centerRect.X+centerRect.W &&
				my >= centerRect.Y && my <= centerRect.Y+centerRect.H {
				if scrollableCenterColumn.OnScroll != nil {
					scrollableCenterColumn.OnScroll(&event)
				}
			}

			// Check if mouse is over right column
			rightRect := scrollableLVHRightColumn.GetRect()
			if mx >= rightRect.X && mx <= rightRect.X+rightRect.W &&
				my >= rightRect.Y && my <= rightRect.Y+rightRect.H {
				if scrollableLVHRightColumn.OnScroll != nil {
					scrollableLVHRightColumn.OnScroll(&event)
				}
			}
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
	// Check center column cards (with scrolling)
	centerRect := scrollableCenterColumn.GetRect()
	if float32(y) > centerRect.Y && float32(y) < centerRect.Y+centerRect.H {
		for _, elem := range scrollableCenterColumn.GetElements() {
			if uiCard, ok := elem.(ui.UICard); ok {
				rect := *uiCard.GetRect()
				rect.Y -= scrollableCenterColumn.GetScrollY()
				if ui.HitTest(&rect, int32(x), int32(y)) {
					selectedCard = uiCard.GetCard()
					return nil
				}
			}
		}
	}
	// Check right column cards (with scrolling)
	rightRect := scrollableLVHRightColumn.GetRect()
	if float32(y) > rightRect.Y && float32(y) < rightRect.Y+rightRect.H {
		for _, uiCard := range playerUICards {
			rect := *uiCard.GetRect()
			rect.Y -= scrollableLVHRightColumn.GetScrollY()
			if ui.HitTest(&rect, int32(x), int32(y)) {
				selectedCard = uiCard.GetCard()
				return nil
			}
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
		ui.NewHud(sdl.FRect{X: gap, Y: 0, W: widthColA, H: float32(data.ScreenHeight)}, sdl.Color{R: 40, G: 40, B: 50, A: 255}),
		ui.NewHud(sdl.FRect{X: widthColA + (2 * gap), Y: 0, W: widthColB, H: float32(data.ScreenHeight)}, sdl.Color{R: 45, G: 45, B: 55, A: 255}),
		ui.NewHud(sdl.FRect{X: widthColA + widthColB + 3*gap, Y: 0, W: widthColC, H: float32(data.ScreenHeight)}, sdl.Color{R: 40, G: 40, B: 50, A: 255}),
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

	// Calculate cards per row
	maxCols := int((rec.W - gap) / (cardWidth + cardGap))
	if maxCols < 1 {
		maxCols = 1
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
		numRows := (length + maxCols - 1) / maxCols
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

	uiCard := ui.CreateUICard(selectedCard, cardRect, editingDeck.CountCard(selectedCard))
	elements[0] = uiCard

	currentDeckSize := len(editingDeck.GetCards())
	currentCardCountInDeck := editingDeck.CountCard(selectedCard)
	currentPlayerCardCount := playerCardDict[selectedCard.GetId()]

	// Add button (disabled if limits reached)
	if currentDeckSize >= 40 || currentCardCountInDeck >= 3 || currentPlayerCardCount > 3-currentCardCountInDeck {
		elements[1] = ui.NewTextBox("Add to deck", sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30}, sdl.Color{R: 60, G: 60, B: 60, A: 255}, sdl.Color{R: 160, G: 160, B: 160, A: 255}, ui.GetDefaultFont(20))
	} else {
		elements[1] = ui.NewButton(
			"Add to deck",
			sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30},
			sdl.Color{R: 20, G: 20, B: 20, A: 100},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			ui.GetDefaultFont(20),
			func() *ui.AppState {
				editingDeck.SetCards(append(editingDeck.GetCards(), selectedCard))
				return nil
			},
		)
	}
	// Remove button (disabled if none in deck)
	if currentCardCountInDeck == 0 {
		elements[2] = ui.NewTextBox("Remove from deck", sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30}, sdl.Color{R: 60, G: 60, B: 60, A: 255}, sdl.Color{R: 160, G: 160, B: 160, A: 255}, ui.GetDefaultFont(20))
	} else {
		elements[2] = ui.NewButton(
			"Remove from deck",
			sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30},
			sdl.Color{R: 120, G: 60, B: 60, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			ui.GetDefaultFont(20),
			func() *ui.AppState {
				editingDeck.RemoveCard(selectedCard)
				return nil
			},
		)
	}
	return elements
}

// Configures the center scrollable column with unique cards from deck.
func setScrollableCenterColumn(centerColRect *sdl.FRect) {
	// Get unique cards and their counts
	cardCounts := make(map[int]int)
	uniqueCards := make([]data.Card, 0)
	seenIds := make(map[int]bool)

	for _, card := range editingDeck.GetCards() {
		cardCounts[card.GetId()]++
		if !seenIds[card.GetId()] {
			seenIds[card.GetId()] = true
			uniqueCards = append(uniqueCards, card)
		}
	}

	// Create UI elements for unique cards
	elements := make([]ui.Element, len(uniqueCards))
	for i, card := range uniqueCards {
		elements[i] = ui.CreateUICard(card, sdl.FRect{}, cardCounts[card.GetId()])
	}

	scrollableCenterColumn.SetElements(elements)

	// Calculate cards per row
	maxCols := int((centerColRect.W - cardGap) / (cardWidth + cardGap))
	if maxCols < 1 {
		maxCols = 1
	}

	// Set the scrollable area rectangle (inside the center column)
	scrollYOrigin := centerColRect.Y + cardGap
	height := centerColRect.H - cardGap
	if height < 0 {
		height = 0
	}

	scrollableCenterColumn.SetRect(sdl.FRect{
		X: centerColRect.X + cardGap/2,
		Y: scrollYOrigin,
		W: centerColRect.W - cardGap,
		H: height,
	})

	// Configure scroll handler
	scrollableCenterColumn.OnScroll = func(event *sdl.Event) {
		y := event.Wheel().Y
		scrollableCenterColumn.SetScrollY(scrollableCenterColumn.GetScrollY() - (float32(y)*cardGap)*scrollableCenterColumn.GetScrollSpeed())
		if scrollableCenterColumn.GetScrollY() < 0 {
			scrollableCenterColumn.SetScrollY(0)
		}
		numRows := (len(uniqueCards) + maxCols - 1) / maxCols
		maxScroll := float32(numRows)*(cardHeight+cardGap) - scrollableCenterColumn.GetRect().H
		if maxScroll < 0 {
			maxScroll = 0
		}
		if scrollableCenterColumn.GetScrollY() > maxScroll {
			scrollableCenterColumn.SetScrollY(maxScroll)
		}
	}
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
			sdl.Color{R: 80, G: 80, B: 80, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateDeckMenu} },
		),
		ui.NewButton(
			"Save",
			sdl.FRect{X: float32(data.ScreenWidth) - 100, Y: 0, W: 50, H: 50},
			sdl.Color{R: 60, G: 120, B: 80, A: 255},
			sdl.Color{R: 240, G: 240, B: 240, A: 255},
			font,
			func() *ui.AppState {
				if deckId == "" {
					// New deck - assign a new ID first
					deckId = data.ID()
					editingDeck.SetId(deckId)
				}
				data.SaveDeck(editingDeck)
				selectedDeck = data.GetDeckById(deckId)
				return nil
			},
		),
	}
}
