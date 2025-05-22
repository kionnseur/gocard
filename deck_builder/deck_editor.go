package deck_builder

import (
	"gocard/data"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

var (
	selectedCard data.Card
	selectedDeck *data.Deck
)

func RenderDeckEditor(renderer *sdl.Renderer, window *sdl.Window, deck_id string) ui.AppState {
	selectedDeck = data.GetDeckById(deck_id)

	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 0, 165, 0, 100)
		sdl.RenderClear(renderer)

		UIColumn, centerColRect := getColumnUI()
		uiLeftColumn := getLeftColomnUI()
		uiCenterColumn := getCenterColomnUI(centerColRect)

		// Affichage de la liste des decks
		for _, e := range UIColumn {
			e.Draw(renderer)
		}
		for _, e := range uiLeftColumn {
			e.Draw(renderer)
		}
		for _, e := range uiCenterColumn {
			e.Draw(renderer)
		}

		// Gestion des événements
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				//si clique sur une carte, on la selectionne
				for _, uiCard := range uiCenterColumn {
					if event.Button().X > uiCard.GetRect().X && event.Button().X < uiCard.GetRect().X+uiCard.GetRect().W &&
						event.Button().Y > uiCard.GetRect().Y && event.Button().Y < uiCard.GetRect().Y+uiCard.GetRect().H {
						selectedCard = uiCard.GetCard()
					}
				}
				// check less btn de la colonne de gauche
				for _, btn := range uiLeftColumn {
					if btn, ok := btn.(*ui.Button); ok && event.Button().X > btn.Rect.X && event.Button().X < btn.Rect.X+btn.Rect.W &&
						event.Button().Y > btn.Rect.Y && event.Button().Y < btn.Rect.Y+btn.Rect.H {
						return btn.OnClick()
					}
				}
			}
		}

		sdl.RenderPresent(renderer)
	}
}

// rectengles des 3 colonnes
func getColumnUI() ([]ui.Element, sdl.FRect) {
	gap := float32(data.ScreenWidth / 48.0)
	widthColA := float32(data.ScreenWidth * 103 / 480)
	widthColB := float32(data.ScreenWidth * 197 / 480)
	widthColC := float32(data.ScreenWidth * 7 / 24)

	centerColRect := sdl.FRect{
		X: widthColA + (2 * gap),
		Y: 0,
		W: widthColB,
		H: float32(data.ScreenHeight),
	}

	return []ui.Element{
		&ui.Hud{Rect: sdl.FRect{X: gap, Y: 0, W: widthColA, H: float32(data.ScreenHeight)}, Color: sdl.Color{R: 255, G: 165, B: 0, A: 255}},
		&ui.Hud{Rect: centerColRect, Color: sdl.Color{R: 0, G: 255, B: 165, A: 255}},
		&ui.Hud{Rect: sdl.FRect{X: widthColA + widthColB + 3*gap, Y: 0, W: widthColC, H: float32(data.ScreenHeight)}, Color: sdl.Color{R: 165, G: 0, B: 255, A: 255}},
	}, centerColRect
}

// retourne la liste des cartes à afficher dans la colonne centrale
func getCenterColomnUI(centerColRect sdl.FRect) []ui.UICard {
	uiCenterColumn := make([]ui.UICard, 0, len(selectedDeck.Cards))
	cardWidth := float32(100)
	cardHeight := float32(150)
	gap := float32(10)

	// Calcul du nombre de cartes qui tiennent dans la largeur de la colonne
	maxColCards := int((centerColRect.W + gap) / (cardWidth + gap))
	// maxRowCards := int(centerColRect.H / (cardHeight + gap))

	startX := centerColRect.X + (centerColRect.W-float32(maxColCards)*(cardWidth+gap)+gap)/2
	y := 2 * gap

	for i, card := range deck.Cards {

		x := startX + float32(i%maxColCards)*(cardWidth+gap)
		if i%maxColCards == 0 && i != 0 {
			y += cardHeight + gap
		}
		cardRect := sdl.FRect{X: x, Y: y, W: cardWidth, H: cardHeight}
		uiCard := ui.CreateUICard(card, cardRect)
		uiCenterColumn = append(uiCenterColumn, uiCard)
	}
	return uiCenterColumn
}

func getLeftColomnUI() []ui.Element {
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
