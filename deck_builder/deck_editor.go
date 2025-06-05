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
	playerCardDict map[int]int

	cardWidth   float32 = 100
	cardHeight  float32 = 150
	cardGap     float32 = 10
	maxColCards int
)

// Déclare le HUD de la listview à l'extérieur
var scrollableLVHRightColumn ui.UIScrollableGridView

func RenderDeckEditor(renderer *sdl.Renderer, window *sdl.Window, deck_id *string) ui.AppState {
	scrollableLVHRightColumn = *(ui.NewUIScrollableGridView(renderer, sdl.FRect{}, sdl.Color{R: 100, G: 100, B: 100, A: 50}, 3, *ui.NewGridConfig(cardWidth, cardHeight, cardGap)))

	selectedDeck = data.GetDeckById(*deck_id)

	playerCardDict := data.GetPlayerCards()
	playerUICards := getPlayerCardListUI(playerCardDict)
	slices.SortFunc(playerUICards, func(a, b ui.UICard) int {
		return a.GetCard().GetId() - b.GetCard().GetId()
	})

	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 0, 165, 0, 100)
		sdl.RenderClear(renderer)

		UIColumn := getColumnUI()
		// Dessine les 3 colonnes
		uiLeftColumn := getLeftColumnUI()
		uiCenterColumn := getDeckCardListUI(UIColumn[1].GetRect())

		//colonne scrollable de droite
		rec := UIColumn[2].GetRect()
		scrollableLVHRightColumn.SetRect(sdl.FRect{X: rec.X + gap/2, Y: rec.Y + 3*gap, W: UIColumn[2].GetRect().W - gap, H: float32(data.ScreenHeight) - scrollableLVHRightColumn.GetRect().Y - gap})

		scrollableLVHRightColumn.SetElements(make([]ui.Element, len(playerUICards)))
		for i, e := range playerUICards {
			scrollableLVHRightColumn.GetElements()[i] = e
		}
		scrollableLVHRightColumn.OnScroll = func(event *sdl.Event) {
			y := event.Wheel().Y
			scrollableLVHRightColumn.SetScrollY(scrollableLVHRightColumn.GetScrollY() - (float32(y)*cardGap)*scrollableLVHRightColumn.GetScrollSpeed())
			if scrollableLVHRightColumn.GetScrollY() < 0 {
				scrollableLVHRightColumn.SetScrollY(0)
			}
			numRows := (len(playerUICards) + maxColCards - 1) / maxColCards
			maxScroll := float32(numRows)*(cardHeight+cardGap) - scrollableLVHRightColumn.GetRect().H
			if maxScroll < 0 {
				maxScroll = 0
			}
			if scrollableLVHRightColumn.GetScrollY() > maxScroll {
				scrollableLVHRightColumn.SetScrollY(maxScroll)
			}
		}

		// 3 colonnes
		for _, e := range UIColumn {
			e.Draw(renderer)
		}
		// colonne de gauche
		for _, e := range uiLeftColumn {
			e.Draw(renderer)
		}
		// colonne de du milieu,
		for _, e := range uiCenterColumn {
			e.Draw(renderer)
		}
		// colonne de droite avec scrollview
		scrollableLVHRightColumn.Draw(renderer)

		// Gestion des événements
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				x, y := event.Button().X, event.Button().Y
				// check les btn de la colonne de gauche
				for _, btn := range uiLeftColumn {
					if btn, ok := btn.(*ui.Button); ok && x > btn.GetRect().X && x < btn.GetRect().X+btn.GetRect().W &&
						y > btn.GetRect().Y && y < btn.GetRect().Y+btn.GetRect().H {
						return btn.OnClick()
					}
				}
				//si clique sur une carte, on la selectionne
				for _, uiCard := range uiCenterColumn {
					if x > uiCard.GetRect().X && x < uiCard.GetRect().X+uiCard.GetRect().W &&
						y > uiCard.GetRect().Y && y < uiCard.GetRect().Y+uiCard.GetRect().H {
						selectedCard = uiCard.GetCard()
					}
				}
				// check les cartes de la colonne de droite
				for _, uiCard := range playerUICards {
					if y > scrollableLVHRightColumn.GetRect().Y && y < scrollableLVHRightColumn.GetRect().Y+scrollableLVHRightColumn.GetRect().H {
						if x > uiCard.GetRect().X && x < uiCard.GetRect().X+uiCard.GetRect().W &&
							y > uiCard.GetRect().Y && y < uiCard.GetRect().Y+uiCard.GetRect().H {
							selectedCard = uiCard.GetCard()
						}
					}
				}
			case sdl.EventMouseWheel:
				scrollableLVHRightColumn.OnScroll(&event)
			}
		}

		sdl.RenderPresent(renderer)
	}
}

// rectengles des 3 colonnes
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

// retourne UI carte selectionnée et les boutons d'ajout/retrait
func getLeftColumnUI() []ui.Element {
	if selectedCard == nil {
		return nil
	}

	elements := make([]ui.Element, 3)

	gap := float32(data.ScreenWidth / 48.0)
	cardWidth := (float32(data.ScreenWidth)*103.0)/480.0 - 2*gap
	cardHeight := cardWidth * 1.5
	cardRect := sdl.FRect{X: float32(2 * data.ScreenWidth / 48), Y: 40, W: cardWidth, H: cardHeight}

	uiCard := ui.CreateUICard(selectedCard, cardRect, selectedDeck.CountCard(selectedCard))
	elements[0] = uiCard

	currentDeckSize := len(selectedDeck.GetCards())
	currentCardCountInDeck := selectedDeck.CountCard(selectedCard)
	currentPlayerCardCount := playerCardDict[selectedCard.GetId()]

	// est-ce que je gere l'érreur de carte null ?
	if currentDeckSize >= 40 || currentCardCountInDeck >= 3 || currentPlayerCardCount > 3-currentCardCountInDeck {
		elements[1] = ui.NewTextBox("Ajouter au deck", sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30}, sdl.Color{R: 80, G: 80, B: 80, A: 255}, sdl.Color{R: 255, G: 0, B: 0, A: 255}, ui.GetDefaultFont(20))
	} else {
		elements[1] = ui.NewButton(
			"Ajouter au deck",
			sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30},
			sdl.Color{R: 20, G: 20, B: 20, A: 100},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			ui.GetDefaultFont(20),
			func() ui.AppState {
				selectedDeck.SetCards(append(selectedDeck.GetCards(), selectedCard))
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": selectedDeck.GetId()}}
			},
		)
	}
	if currentCardCountInDeck == 0 {

		elements[2] = ui.NewTextBox("Retirer du deck", sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30}, sdl.Color{R: 80, G: 80, B: 80, A: 255}, sdl.Color{R: 255, G: 0, B: 0, A: 255}, ui.GetDefaultFont(20))
	} else {

		elements[2] = ui.NewButton(
			"Retirer du deck",
			sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30},
			sdl.Color{R: 20, G: 20, B: 20, A: 100},
			sdl.Color{R: 255, G: 255, B: 255, A: 255},
			ui.GetDefaultFont(20),
			func() ui.AppState {
				selectedDeck.RemoveCard(selectedCard)
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": selectedDeck.GetId()}}
			},
		)
	}
	return elements

}

// retourne la liste des cartes du deck à afficher dans la colonne centrale
func getDeckCardListUI(centerColRect *sdl.FRect) []ui.UICard {
	uiCenterColumn := make([]ui.UICard, 0, len(selectedDeck.GetCards()))

	// Calcul du nombre de cartes qui tiennent dans la largeur de la colonne
	maxColCards = int((centerColRect.W + cardGap) / (cardWidth + cardGap))
	// maxRowCards := int(centerColRect.H / (cardHeight + gap))

	startX := centerColRect.X + (centerColRect.W-float32(maxColCards)*(cardWidth+cardGap)+cardGap)/2
	y := 2 * cardGap

	for i, card := range selectedDeck.GetCards() {

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

// retourne la liste totale des cartes du joueur à afficher
func getPlayerCardListUI(playerCardDict map[int]int) []ui.UICard {
	elements := make([]ui.UICard, 0, len(playerCardDict))
	for cardId, qty := range playerCardDict {
		card := data.GetAllCards()[cardId]
		elements = append(elements, ui.CreateUICard(card, sdl.FRect{}, qty))
	}
	return elements
}
