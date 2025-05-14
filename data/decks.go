package data

import (
	"math/rand"
	"strconv"
	"sync"
)

var (
	deckList     []Deck
	deckListOnce sync.Once
)

func loadDeckList() {
	count := 20
	cards := make([]Card, 0, count)
	for i := 0; i < count; i++ {
		cards = append(cards, &MonsterCard{Name: "Card 1", Image: "card1.png", Description: "Description 1", Level: 1, Attack: 100, Defense: 50})
	}

	defaultDeck := Deck{
		ID:    "default",
		Name:  "Default Deck",
		Cards: cards,
	}

	deckList = make([]Deck, 0, 30)
	deckList = append(deckList, defaultDeck)
	for i := 0; i < 30; i++ {
		deckList = append(deckList, debug_create_random_deck())
	}
	print("DeckList loaded")

	// file, err := os.Open("decks.json")
	// if err != nil {
	// 	deckList = []Deck{}
	// 	return
	// }
	// defer file.Close()
	// json.NewDecoder(file).Decode(&deckList)
}

func GetDeckList() []Deck {
	deckListOnce.Do(loadDeckList)
	return deckList
}

func GetDeckById(id string) Deck {
	deckListOnce.Do(loadDeckList)
	for _, deck := range deckList {
		if deck.ID == id {
			return deck
		}
	}
	return Deck{}
}
func DeleteDeckById(id string) {
	deckListOnce.Do(loadDeckList)
	for i, deck := range deckList {
		if deck.ID == id {
			deckList = append(deckList[:i], deckList[i+1:]...)
			return
		}
	}
}

func debug_create_random_deck() Deck {
	// retour un deck de 40 cartes
	// avec un nom random

	count := 40
	cards := make([]Card, 0, count)
	for i := 0; i < count; i++ {
		cards = append(cards, &MonsterCard{Name: "Card 1", Image: "card1.png", Description: "Description 1", Level: 1, Attack: 100, Defense: 50})
	}
	id := "" + strconv.Itoa(rand.Intn(1000))
	randomName := "Deck " + id
	return Deck{
		ID:    id,
		Name:  randomName,
		Cards: cards,
	}
}
