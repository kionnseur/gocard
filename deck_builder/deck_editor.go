package deck_builder

import (
	"gocard/data"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

var (
	selectedCard data.Card
)

func RenderDeckEditor(renderer *sdl.Renderer, window *sdl.Window, deck data.Deck) ui.AppState {
	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 0, 165, 0, 100)
		sdl.RenderClear(renderer)

		UIColumn, centerColRect := GetColumnUI()
		uiCenterColumn := GetCenterColomnUI(deck, centerColRect)

		// Affichage de la liste des decks
		for _, e := range UIColumn {
			e.Draw(renderer)
		}
		for _, e := range uiCenterColumn {
			e.Draw(renderer)
		}

		if selectedCard != nil {
			DrawSelectedCard(renderer, selectedCard)
		}

		// Gestion des événements
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				//si clique sur une carte, on la selectionne
				for _, UICard := range uiCenterColumn {
					if monsterCard, ok := UICard.(*data.UIMonsterCard); ok {
						if event.Button().X > monsterCard.Rect.X && event.Button().X < monsterCard.Rect.X+monsterCard.Rect.W &&
							event.Button().Y > monsterCard.Rect.Y && event.Button().Y < monsterCard.Rect.Y+monsterCard.Rect.H {
							selectedCard = monsterCard.Card
						}
					}
				}
			}
		}

		sdl.RenderPresent(renderer)
	}
}

// get colonnes elements
func GetColumnUI() ([]ui.Element, sdl.FRect) {
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

func GetCenterColomnUI(deck data.Deck, centerColRect sdl.FRect) []data.UICard {
	uiCenterColumn := make([]data.UICard, 0, len(deck.Cards))
	cardWidth := float32(50)
	cardHeight := float32(75)
	gap := float32(20)

	// Calcul du nombre de cartes qui tiennent dans la largeur de la colonne
	maxColCards := int((centerColRect.W + gap) / (cardWidth + gap))
	maxRowCards := int(centerColRect.H / (cardHeight + gap))
	maxCards := maxColCards * maxRowCards

	startX := centerColRect.X + (centerColRect.W-float32(maxColCards)*(cardWidth+gap)+gap)/2
	y := 2 * gap

	for i, card := range deck.Cards {
		if i >= maxCards {
			break
		}
		x := startX + float32(i%maxColCards)*(cardWidth+gap)
		if i%maxColCards == 0 && i != 0 {
			y += cardHeight + gap
		}
		cardRect := sdl.FRect{X: x, Y: y, W: cardWidth, H: cardHeight}
		uiCard := &data.UIMonsterCard{
			Card: card,
			Rect: cardRect,
		}

		uiCenterColumn = append(uiCenterColumn, uiCard)
	}
	return uiCenterColumn
}

func DrawSelectedCard(renderer *sdl.Renderer, card data.Card) {
	cardWidth := float32(300)
	cardHeight := float32(450)

	cardRect := sdl.FRect{X: 40, Y: 40, W: cardWidth, H: cardHeight}
	uiCard := &data.UIMonsterCard{
		Card: card,
		Rect: cardRect,
	}
	uiCard.Draw(renderer)
}
