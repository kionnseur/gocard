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
	deck       *data.Deck
	lastDeckId string
	font       = ui.GetDefaultFont(24)

	scrollableLVHDeckList = &ui.UIScrollableStackView{
		Rect:        sdl.FRect{X: gap, Y: 80, W: 200, H: float32(data.ScreenHeight) - gap - 80}, // 80 pour le Y
		Color:       sdl.Color{R: 100, G: 100, B: 100, A: 50},
		ScrollY:     0,
		ScrollSpeed: 15,
	}
)

func RenderDeckMenu(renderer *sdl.Renderer, window *sdl.Window, appState ui.AppState) ui.AppState {
	var buttons []ui.Button
	var uiDeckInfoBtn []ui.Button
	var uiDeckInfo []ui.Element

	// boutons de la liste des decks
	uiDeckListElements := uiGetDeckListElements(data.GetDeckList(), scrollableLVHDeckList.GetRect())
	scrollableLVHDeckList.Elements = uiDeckListElements
	scrollableLVHDeckList.OnScroll = func(event *sdl.Event) {
		y := event.Wheel().Y
		scrollableLVHDeckList.ScrollY -= float32(y) * gap
		if scrollableLVHDeckList.ScrollY < 0 {
			scrollableLVHDeckList.ScrollY = 0
		}
		maxScroll := float32(len(uiDeckListElements))*35 - scrollableLVHDeckList.Rect.H
		if scrollableLVHDeckList.ScrollY > maxScroll {
			scrollableLVHDeckList.ScrollY = maxScroll
		}
	}
	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 255, 165, 0, 255)
		sdl.RenderClear(renderer)

		action := appState.Data["action"]
		deckId := appState.Data["deckId"]
		if deckId != lastDeckId {
			deck = data.GetDeckById(deckId)
			lastDeckId = deckId
		}
		// supprime avant d'afficher la liste
		if action == "ask" {
			uiDeckInfo, uiDeckInfoBtn = uiGetDeckInfo(deck, scrollableLVHDeckList.GetRect())
			for _, e := range uiDeckInfo {
				e.Draw(renderer)
			}
			for _, e := range uiDeckInfoBtn {
				e.Draw(renderer)
			}
		} else if action == "edit" {
			return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.ID, "action": "edit"}}
		} else if action == "delete" && deck.ID != "" {
			data.DeleteDeckById(deck.ID)
			// Réinitialise l'état pour revenir à la liste
			appState.Data["action"] = ""
			uiDeckListElements = uiGetDeckListElements(data.GetDeckList(), scrollableLVHDeckList.GetRect())
		} else if action == "duplicate" && deck.ID != "" {
			data.DuplicateDeckById(deck.ID)
			appState.Data["action"] = ""
			uiDeckListElements = uiGetDeckListElements(data.GetDeckList(), scrollableLVHDeckList.GetRect())
		}

		// Affiche la liste des decks, colonne de gauche et btn retour
		scrollableLVHDeckList.Rect.H = float32(data.ScreenHeight) - gap - scrollableLVHDeckList.Rect.Y

		scrollableLVHDeckList.Draw(renderer)

		buttons = getDeckMenuButtons(scrollableLVHDeckList.GetRect())
		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		sdl.RenderPresent(renderer)

		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				x, y := event.Button().X, event.Button().Y
				// nouveau deck & retour
				for _, btn := range buttons {
					if x > btn.Rect.X && x < btn.Rect.X+btn.Rect.W &&
						y > btn.Rect.Y && y < btn.Rect.Y+btn.Rect.H {
						return btn.OnClick()
					}
				}
				// liste de deck
				for _, e := range uiDeckListElements {
					if btn, ok := e.(*ui.Button); ok {
						rect := btn.Rect
						rect.Y -= scrollableLVHDeckList.ScrollY
						if x > rect.X && x < rect.X+rect.W &&
							y > rect.Y && y < rect.Y+rect.H {
							return btn.OnClick()
						}
					}
				}
				// deck info
				if action == "ask" {
					for _, btn := range uiDeckInfoBtn {
						if x > btn.Rect.X && x < btn.Rect.X+btn.Rect.W &&
							y > btn.Rect.Y && y < btn.Rect.Y+btn.Rect.H {
							return btn.OnClick()
						}
					}
				}

			case sdl.EventMouseWheel:
				if scrollableLVHDeckList.OnScroll != nil {
					scrollableLVHDeckList.OnScroll(&event)
				}
			}
		}
	}
}

func getDeckMenuButtons(parent *sdl.FRect) []ui.Button {

	return []ui.Button{
		{
			Rect:      sdl.FRect{X: parent.X, Y: parent.Y - 30, W: 200, H: 30},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 100},
			Text:      "Nouveau Deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": "", "action": "new"}}
			},
		},
		{
			Rect:      sdl.FRect{X: float32(data.ScreenWidth) - 50, Y: 0, W: 50, H: 50},
			Color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
			Text:      "Retour ⬅️",
			TextColor: sdl.Color{R: 255, G: 0, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.AppState{State: ui.StateStartMenu} },
		},
	}
}

func uiGetDeckListElements(decksList []data.Deck, parent *sdl.FRect) []ui.Element {
	// Crée une liste d'éléments de type Button représentent chaque deck
	elements := make([]ui.Element, len(decksList))
	for i, deck := range decksList {
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = &ui.Button{
			Rect:      sdl.FRect{X: parent.X, Y: parent.Y + float32(i*35), W: parent.W, H: 30},
			Color:     sdl.Color{R: r, G: g, B: b, A: 255},
			Text:      deck.Name,
			TextColor: sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckMenu, Data: map[string]string{"deckId": deck.ID, "action": "ask"}}
			},
		}
	}
	return elements
}

func uiGetDeckInfo(deck *data.Deck, parent *sdl.FRect) ([]ui.Element, []ui.Button) {

	offset := parent.X + gap + parent.W
	// affiche le nom des 3 premiere cartes du deck

	elements := make([]ui.Element, 3)
	for i, card := range deck.Cards {
		if i > 2 {
			break
		}
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = &ui.TextBox{
			Rect:      sdl.FRect{X: offset, Y: float32((i + 5) * 35), W: 200, H: 30},
			Color:     sdl.Color{R: r, G: g, B: b, A: 255},
			Text:      card.GetName(),
			TextColor: sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			Font:      font,
		}
	}
	buttons := []ui.Button{
		{
			Rect:      sdl.FRect{X: offset, Y: float32(8 * 35), W: 200, H: 30},
			Color:     sdl.Color{R: 100, G: 200, B: 100, A: 255},
			Text:      "Editer",
			TextColor: sdl.Color{R: 155, G: 55, B: 155, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckMenu, Data: map[string]string{"deckId": deck.ID, "action": "edit"}}
			},
		},
		{
			Rect:      sdl.FRect{X: offset, Y: float32(9 * 35), W: 200, H: 30},
			Color:     sdl.Color{R: 100, G: 100, B: 200, A: 255},
			Text:      "Dupliquer",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckMenu, Data: map[string]string{"deckId": deck.ID, "action": "duplicate"}}
			},
		},
		{
			Rect:      sdl.FRect{X: offset, Y: float32(10 * 35), W: parent.W, H: 30},
			Color:     sdl.Color{R: 200, G: 100, B: 100, A: 255},
			Text:      "Supprimer",
			TextColor: sdl.Color{R: 55, G: 155, B: 155, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckMenu, Data: map[string]string{"deckId": deck.ID, "action": "delete"}}
			},
		},
	}

	return elements, buttons

}
