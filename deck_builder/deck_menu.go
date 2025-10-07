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
)

func RenderDeckMenu(renderer *sdl.Renderer, window *sdl.Window, appState *ui.AppState) ui.AppState {
	var buttons []*ui.Button
	var uiDeckInfoBtns []*ui.Button
	var uiDeckInfo []ui.Element

	for {
		sdl.GetWindowSize(window, &data.ScreenWidth, &data.ScreenHeight)
		sdl.SetRenderDrawColor(renderer, 255, 165, 0, 255)
		sdl.RenderClear(renderer)

		if askedDeckId != lastDeckId {
			deck = data.GetDeckById(askedDeckId)
			lastDeckId = askedDeckId
		}

		// supprime avant d'afficher la liste
		if askedDeckId != "" {
			uiDeckInfo, uiDeckInfoBtns = uiGetDeckInfo(deck, scrollableLVHDeckList.GetRect())
			for _, e := range uiDeckInfo {
				e.Draw(renderer)
			}
			for _, e := range uiDeckInfoBtns {
				e.Draw(renderer)
			}
		}
		// boutons de la liste des decks
		updateDeckListElements()

		// Affiche la liste des decks, colonne de gauche et btn retour
		scrollableLVHDeckList.GetRect().H = float32(data.ScreenHeight) - gap - scrollableLVHDeckList.GetRect().Y

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
					if x > btn.GetRect().X && x < btn.GetRect().X+btn.GetRect().W &&
						y > btn.GetRect().Y && y < btn.GetRect().Y+btn.GetRect().H {
						as := btn.OnClick()
						if as != nil {
							return *as
						}
					}
				}
				// liste de deck
				for _, e := range uiDeckListElements {
					if y > scrollableLVHDeckList.GetRect().Y && y < scrollableLVHDeckList.GetRect().Y+scrollableLVHDeckList.GetRect().H {
						if btn, ok := e.(*ui.Button); ok {
							rect := btn.GetRect()
							rect.Y -= scrollableLVHDeckList.GetScrollY()
							if x > rect.X && x < rect.X+rect.W &&
								y > rect.Y && y < rect.Y+rect.H {
								as := btn.OnClick()
								if as != nil {
									return *as
								}
							}
						}
					}
				}
				// deck info, uiDeckInfoBtn est vide si action != ask
				for _, btn := range uiDeckInfoBtns {
					if x > btn.GetRect().X && x < btn.GetRect().X+btn.GetRect().W &&
						y > btn.GetRect().Y && y < btn.GetRect().Y+btn.GetRect().H {
						as := btn.OnClick()
						if as != nil {
							return *as
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

func getDeckMenuButtons(parent *sdl.FRect) []*ui.Button {
	return []*ui.Button{
		ui.NewButton(
			"Nouveau Deck",
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
			"Retour ⬅️",
			sdl.FRect{X: float32(data.ScreenWidth) - 50, Y: 0, W: 50, H: 50},
			sdl.Color{R: 0, G: 255, B: 0, A: 255},
			sdl.Color{R: 255, G: 0, B: 255, A: 255},
			font,
			func() *ui.AppState { return &ui.AppState{State: ui.StateStartMenu} },
		),
	}
}

func uiGetDeckListElements(decksList []data.Deck, parent *sdl.FRect) []ui.Element {
	// Crée une liste d'éléments de type Button représentant chaque deck
	elements := make([]ui.Element, len(decksList))
	for i, deck := range decksList {
		var r, g, b = ui.ColorBreathSin(i * 10)
		// Utilise le constructeur NewButton avec les bons paramètres
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

func uiGetDeckInfo(deck *data.Deck, parent *sdl.FRect) ([]ui.Element, []*ui.Button) {

	offset := parent.X + gap + parent.W
	// affiche le nom des 3 premiere cartes du deck

	elements := make([]ui.Element, 3)
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
			"Editer",
			sdl.FRect{X: offset, Y: float32(8 * 35), W: 200, H: 30},
			sdl.Color{R: 100, G: 200, B: 100, A: 255},
			sdl.Color{R: 155, G: 55, B: 155, A: 255},
			font,
			func() *ui.AppState {
				return &ui.AppState{State: ui.StateDeckBuilder, Data: map[string]string{"deckId": deck.GetId(), "action": "edit"}}
			},
		),

		ui.NewButton(
			"Dupliquer",
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
			"Supprimer",
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
		if scrollableLVHDeckList.GetScrollY() > maxScroll {
			scrollableLVHDeckList.SetScrollY(maxScroll)
		}
	}
}
