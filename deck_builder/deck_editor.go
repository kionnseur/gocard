package deck_builder

import (
	"gocard/data"
	"gocard/ui"
	"slices"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

var (
	selectedCard        data.Card
	selectedDeck        *data.Deck
	scrollLevelRightCol float32 = 0.0
	cardWidth           float32 = 100
	cardHeight          float32 = 150
	cardGap             float32 = 10
	maxColCards         int
)

// Déclare le HUD de la listview à l'extérieur
var scrollableLVHRightColumn = &ui.Hud{
	Rect:  sdl.FRect{X: 0, Y: 0, W: 0, H: float32(data.ScreenHeight) - gap},
	Color: sdl.Color{R: 100, G: 100, B: 100, A: 50},
}

func RenderDeckEditor(renderer *sdl.Renderer, window *sdl.Window, deck_id string) ui.AppState {
	selectedDeck = data.GetDeckById(deck_id)

	playerCardDict := data.GetPlayerCards()

	playerCards := make([]data.Card, 0, len(playerCardDict))
	for cardId, count := range playerCardDict {
		for i := 0; i < count; i++ {
			playerCards = append(playerCards, data.GetAllCards()[cardId])
		}
	}

	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 0, 165, 0, 100)
		sdl.RenderClear(renderer)

		UIColumn := getColumnUI()

		//colonne scrolable de droite
		rec := UIColumn[2].GetRect()
		scrollableLVHRightColumn.Rect.X = rec.X + gap/2
		scrollableLVHRightColumn.Rect.Y = rec.Y + 3*gap // remplacer par hauter des menus
		scrollableLVHRightColumn.Rect.W = UIColumn[2].GetRect().W - gap
		scrollableLVHRightColumn.Rect.H = float32(data.ScreenHeight) - scrollableLVHRightColumn.Rect.Y - gap
		// Dessine les 3 colonnes
		uiLeftColumn := getLeftColumnUI()
		uiCenterColumn := getCenterColumnUI(UIColumn[1].GetRect())
		uiRightColumn := getRightColumnUI(playerCards)

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

		// Dessine les éléments de la liste, scroll/clipping selon le HUD
		for _, e := range uiRightColumn {

			e.GetRect().Y -= scrollLevelRightCol
			if e.GetRect().Y+e.GetRect().H > scrollableLVHRightColumn.Rect.Y && e.GetRect().Y < scrollableLVHRightColumn.Rect.Y+scrollableLVHRightColumn.Rect.H {
				e.Draw(renderer)
			}
		}

		// Gestion des événements
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				// check les btn de la colonne de gauche
				for _, btn := range uiLeftColumn {
					if btn, ok := btn.(*ui.Button); ok && event.Button().X > btn.Rect.X && event.Button().X < btn.Rect.X+btn.Rect.W &&
						event.Button().Y > btn.Rect.Y && event.Button().Y < btn.Rect.Y+btn.Rect.H {
						return btn.OnClick()
					}
				}
				//si clique sur une carte, on la selectionne
				for _, uiCard := range uiCenterColumn {
					if event.Button().X > uiCard.GetRect().X && event.Button().X < uiCard.GetRect().X+uiCard.GetRect().W &&
						event.Button().Y > uiCard.GetRect().Y && event.Button().Y < uiCard.GetRect().Y+uiCard.GetRect().H {
						selectedCard = uiCard.GetCard()
					}
				}
				// check les cartes de la colonne de droite
				for _, uiCard := range uiRightColumn {
					if event.Button().X > uiCard.GetRect().X && event.Button().X < uiCard.GetRect().X+uiCard.GetRect().W &&
						event.Button().Y > uiCard.GetRect().Y && event.Button().Y < uiCard.GetRect().Y+uiCard.GetRect().H {
						selectedCard = uiCard.GetCard()
					}
				}
			case sdl.EventMouseWheel:
				y := event.Wheel().Y
				scrollLevelRightCol -= float32(y) * cardGap
				if scrollLevelRightCol < 0 {
					scrollLevelRightCol = 0
				}
				numRows := (len(uiRightColumn) + maxColCards - 1) / maxColCards
				maxScroll := float32(numRows)*(cardHeight+cardGap) - scrollableLVHRightColumn.Rect.H
				if maxScroll < 0 {
					maxScroll = 0
				}
				if scrollLevelRightCol > maxScroll {
					scrollLevelRightCol = maxScroll
				}

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
		&ui.Hud{Rect: sdl.FRect{X: gap, Y: 0, W: widthColA, H: float32(data.ScreenHeight)}, Color: sdl.Color{R: 255, G: 165, B: 0, A: 255}},
		&ui.Hud{Rect: sdl.FRect{X: widthColA + (2 * gap), Y: 0, W: widthColB, H: float32(data.ScreenHeight)}, Color: sdl.Color{R: 0, G: 255, B: 165, A: 255}},
		&ui.Hud{Rect: sdl.FRect{X: widthColA + widthColB + 3*gap, Y: 0, W: widthColC, H: float32(data.ScreenHeight)}, Color: sdl.Color{R: 165, G: 0, B: 255, A: 255}},
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

	uiCard := ui.CreateUICard(selectedCard, cardRect)
	elements[0] = uiCard
	// est-ce que je gere l'érreur de carte null ?
	if len(selectedDeck.Cards) >= 40 || selectedDeck.CountCard(selectedCard) >= 3 {
		elements[1] = &ui.TextBox{
			Rect:      sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30},
			Color:     sdl.Color{R: 80, G: 80, B: 80, A: 255},
			Text:      "Ajouter au deck",
			TextColor: sdl.Color{R: 255, G: 0, B: 0, A: 255},
			Font:      ui.GetDefaultFont(20),
		}
		elements[2] = &ui.Button{
			Rect:      sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30},
			Color:     sdl.Color{R: 20, G: 20, B: 20, A: 100},
			Text:      "Retirer du deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      ui.GetDefaultFont(20),
			OnClick: func() ui.AppState {
				selectedDeck.RemoveCard(selectedCard)
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": selectedDeck.ID}}
			},
		}
	} else if selectedDeck.CountCard(selectedCard) == 0 {
		elements[1] = &ui.Button{
			Rect:      sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30},
			Color:     sdl.Color{R: 20, G: 20, B: 20, A: 100},
			Text:      "Ajouter au deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      ui.GetDefaultFont(20),
			OnClick: func() ui.AppState {
				selectedDeck.Cards = append(selectedDeck.Cards, selectedCard)
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": selectedDeck.ID}}
			},
		}
		elements[2] = &ui.TextBox{
			Rect:      sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30},
			Color:     sdl.Color{R: 80, G: 80, B: 80, A: 255},
			Text:      "Retirer du deck",
			TextColor: sdl.Color{R: 255, G: 0, B: 0, A: 255},
			Font:      ui.GetDefaultFont(20),
		}
	} else {
		elements[1] = &ui.Button{
			Rect:      sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 10, W: cardWidth, H: 30},
			Color:     sdl.Color{R: 20, G: 20, B: 20, A: 100},
			Text:      "Ajouter au deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      ui.GetDefaultFont(20),
			OnClick: func() ui.AppState {
				selectedDeck.Cards = append(selectedDeck.Cards, selectedCard)
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": selectedDeck.ID}}
			},
		}
		elements[2] = &ui.Button{
			Rect:      sdl.FRect{X: cardRect.X, Y: cardRect.Y + cardRect.H + 50, W: cardWidth, H: 30},
			Color:     sdl.Color{R: 20, G: 20, B: 20, A: 100},
			Text:      "Retirer du deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      ui.GetDefaultFont(20),
			OnClick: func() ui.AppState {
				selectedDeck.RemoveCard(selectedCard)
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": selectedDeck.ID}}
			},
		}
	}
	return elements

}

// retourne la liste des cartes du deck à afficher dans la colonne centrale
func getCenterColumnUI(centerColRect *sdl.FRect) []ui.UICard {
	uiCenterColumn := make([]ui.UICard, 0, len(selectedDeck.Cards))

	// Calcul du nombre de cartes qui tiennent dans la largeur de la colonne
	maxColCards = int((centerColRect.W + cardGap) / (cardWidth + cardGap))
	// maxRowCards := int(centerColRect.H / (cardHeight + gap))

	startX := centerColRect.X + (centerColRect.W-float32(maxColCards)*(cardWidth+cardGap)+cardGap)/2
	y := 2 * cardGap

	for i, card := range deck.Cards {

		x := startX + float32(i%maxColCards)*(cardWidth+cardGap)
		if i%maxColCards == 0 && i != 0 {
			y += cardHeight + cardGap
		}
		cardRect := sdl.FRect{X: x, Y: y, W: cardWidth, H: cardHeight}
		uiCard := ui.CreateUICard(card, cardRect)
		uiCenterColumn = append(uiCenterColumn, uiCard)
	}
	return uiCenterColumn
}

// retourne la liste totale des cartes
func getRightColumnUI(playerCards []data.Card) []ui.UICard {

	maxColCards := int((scrollableLVHRightColumn.Rect.W + cardGap) / (cardWidth + cardGap))
	startX := scrollableLVHRightColumn.Rect.X + (scrollableLVHRightColumn.Rect.W-float32(maxColCards)*(cardWidth+cardGap)+cardGap)/2
	y := scrollableLVHRightColumn.Rect.Y + cardGap
	slices.SortFunc(playerCards, func(a, b data.Card) int {
		return a.GetId() - b.GetId()
	})
	elements := make([]ui.UICard, len(playerCards))
	for i, card := range playerCards {
		x := startX + float32(i%maxColCards)*(cardWidth+cardGap)
		if i%maxColCards == 0 && i != 0 {
			y += cardHeight + cardGap
		}
		cardRect := sdl.FRect{X: x, Y: y, W: cardWidth, H: cardHeight}
		uiCard := ui.CreateUICard(card, cardRect)
		elements[i] = uiCard

	}
	return elements
}
