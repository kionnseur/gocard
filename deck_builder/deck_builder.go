package deck_builder

import (
	"gocard/data"
	"gocard/ui"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/jupiterrider/purego-sdl3/ttf"
)

const (
	gap      float32 = 20
	fontSize float32 = 24
)

var (
	// deckBuilderData map[string]string
	scrollLevel float32 = 0.0
	_, font             = ttf.Init(), ttf.OpenFont("assets/fonts/arial.ttf", fontSize)
)

// Déclare le HUD de la listview à l'extérieur
var listViewHud = &ui.Hud{
	Rect:  sdl.FRect{X: gap, Y: 80, W: 200, H: float32(data.ScreenHeight) - gap},
	Color: sdl.Color{R: 100, G: 100, B: 100, A: 50},
}

func RenderDeckBuilder(renderer *sdl.Renderer, window *sdl.Window) ui.AppState {
	var buttons []ui.Button
	var overlayDeckId string
	var overlayActive bool
	appState := ui.AppState{State: ui.StateDeckBuilder}

	ui_deckList := getDeckListElements(data.GetDeckList())

	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)

		// contient le btn nouveau deck et retour
		buttons = getDeckBuilderButtons()

		// Gestion de l'overlay selon appState.Data
		if appState.Data != nil {
			if deckId, ok := appState.Data["deckId"]; ok {
				switch appState.Data["action"] {
				case "ask":
					overlayDeckId = deckId
					overlayActive = true
				case "delete":
					data.DeleteDeckById(deckId)
					// Recharge la liste après suppression
					ui_deckList = getDeckListElements(data.GetDeckList())
					appState = ui.AppState{State: ui.StateDeckBuilder}
					overlayActive = false
				}
			}
			// Reset Data pour éviter de boucler sur l'action
			appState.Data = nil
		}

		var event sdl.Event
		for sdl.PollEvent(&event) {

			switch event.Type() {
			case sdl.EventQuit:
				return ui.AppState{State: ui.StateQuit}
			case sdl.EventMouseButtonDown:
				x, y := event.Button().X, event.Button().Y
				if overlayActive {
					// Clique sur overlay : détecte si on clique sur un bouton de l'overlay
					// (à adapter selon la position réelle des boutons overlay)
					// Si clic en dehors, ferme l'overlay
					if x < 0 || x > listViewHud.Rect.W || y < 0 || y > float32(6*35) {
						overlayActive = false
					}
				} else {
					for _, btn := range buttons {
						if x > btn.Rect.X && x < btn.Rect.X+btn.Rect.W &&
							y > btn.Rect.Y && y < btn.Rect.Y+btn.Rect.H {
							return btn.OnClick()
						}
					}
					for _, e := range ui_deckList {
						if tb, ok := e.(*ui.Button); ok {
							rect := tb.Rect
							rect.Y -= scrollLevel
							if x > rect.X && x < rect.X+rect.W &&
								y > rect.Y && y < rect.Y+rect.H {
								appState = tb.OnClick()
							}
						}
					}
				}
			case sdl.EventMouseWheel:
				if !overlayActive {
					y := event.Wheel().Y
					scrollLevel -= float32(y) * gap
					if scrollLevel < 0 {
						scrollLevel = 0
					}
					maxScroll := float32(len(ui_deckList))*35 - listViewHud.Rect.H
					if scrollLevel > maxScroll {
						scrollLevel = maxScroll
					}
				}
			}
		}

		sdl.SetRenderDrawColor(renderer, 255, 165, 0, 255)
		sdl.RenderClear(renderer)

		// Dessine le HUD de la listview
		// met à jour la taille du HUD selon la taille de la fenêtre à chaque frame
		listViewHud.Rect.H = float32(data.ScreenHeight) - gap - listViewHud.Rect.Y
		listViewHud.Draw(renderer)

		// Dessine les éléments de la liste, scroll/clipping selon le HUD
		for _, e := range ui_deckList {
			if tb, ok := e.(*ui.Button); ok {
				rect := tb.Rect
				rect.Y -= scrollLevel
				if rect.Y+rect.H > listViewHud.Rect.Y && rect.Y < listViewHud.Rect.Y+listViewHud.Rect.H {
					tmp := *tb
					tmp.Rect = rect
					tmp.Draw(renderer)
				}
			}
		}
		for _, btn := range buttons {
			btn.Draw(renderer)
		}

		// Affiche l'overlay si besoin
		if overlayActive && overlayDeckId != "" {
			showOverlay(renderer, overlayDeckId)
		}

		sdl.RenderPresent(renderer)
		sdl.Delay(16)
	}
}

func getDeckBuilderButtons() []ui.Button {

	return []ui.Button{
		{
			Rect:      sdl.FRect{X: listViewHud.Rect.X, Y: listViewHud.Rect.Y - 30, W: 200, H: 30},
			Color:     sdl.Color{R: 0, G: 0, B: 0, A: 100},
			Text:      "Nouveau Deck",
			TextColor: sdl.Color{R: 255, G: 255, B: 255, A: 255},
			Font:      font,
			OnClick:   func() ui.AppState { return ui.AppState{State: ui.StateStartMenu} },
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

func getDeckListElements(decksList []data.Deck) []ui.Element {
	// Crée une liste d'éléments de type TextBox représentent chaque deck
	elements := make([]ui.Element, len(decksList))
	for i, deck := range decksList {
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = &ui.Button{
			Rect:      sdl.FRect{X: listViewHud.Rect.X, Y: listViewHud.Rect.Y + float32(i*35), W: listViewHud.Rect.W, H: 30},
			Color:     sdl.Color{R: r, G: g, B: b, A: 255},
			Text:      deck.Name,
			TextColor: sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.ID, "action": "ask"}}
			},
		}
	}
	return elements
}

func showOverlay(renderer *sdl.Renderer, deckID string) {
	// Crée une overlay avec le texte passé en paramètre
	deck := data.GetDeckById(deckID)
	overlay := &ui.Hud{
		Rect:  sdl.FRect{X: 0, Y: 0, W: float32(data.ScreenWidth), H: float32(data.ScreenHeight)},
		Color: sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
	// Dessine l'overlay
	overlay.Draw(renderer)

	// affiche le nom des 3 premiere cartes du deck
	elements := make([]ui.Element, 3)
	for i, card := range deck.Cards {
		if i > 2 {
			break
		}
		var r, g, b = ui.ColorBreathSin(i * 10)
		elements[i] = &ui.TextBox{
			Rect:      sdl.FRect{X: 0, Y: float32(i * 35), W: listViewHud.Rect.W, H: 30},
			Color:     sdl.Color{R: r, G: g, B: b, A: 255},
			Text:      card.GetName(),
			TextColor: sdl.Color{R: 255 - r, G: 255 - g, B: 255 - b, A: 255},
			Font:      font,
		}
	}

	buttons := []*ui.Button{
		{
			Rect:      sdl.FRect{X: 0, Y: float32(3 * 35), W: listViewHud.Rect.W, H: 30},
			Color:     sdl.Color{R: 100, G: 200, B: 100, A: 255},
			Text:      "Editer",
			TextColor: sdl.Color{R: 155, G: 55, B: 155, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.ID, "action": "delete"}}
			},
		},
		{
			Rect:      sdl.FRect{X: 0, Y: float32(4 * 35), W: listViewHud.Rect.W, H: 30},
			Color:     sdl.Color{R: 200, G: 100, B: 100, A: 255},
			Text:      "Supprimer",
			TextColor: sdl.Color{R: 55, G: 155, B: 155, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.ID, "action": "delete"}}
			},
		},
		{
			Rect:      sdl.FRect{X: 0, Y: float32(5 * 35), W: listViewHud.Rect.W, H: 30},
			Color:     sdl.Color{R: 100, G: 100, B: 200, A: 255},
			Text:      "Annuler",
			TextColor: sdl.Color{R: 155, G: 155, B: 55, A: 255},
			Font:      font,
			OnClick: func() ui.AppState {
				return ui.AppState{State: ui.StateDeckBuilder}
			},
		},
	}

	// Dessine les éléments de l'overlay

	for _, e := range elements {
		e.Draw(renderer)
	}
	for _, btn := range buttons {
		btn.Draw(renderer)
	}

}
