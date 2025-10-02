package data

import (
	"math/rand"
	"slices"
	"strconv"
	"sync"
)

type Deck struct {
	id    string
	name  string
	cards []Card
}

func (d *Deck) GetName() string {
	return d.name
}
func (d *Deck) GetId() string {
	return d.id
}
func (d *Deck) GetCards() []Card {
	return d.cards
}
func (d *Deck) SetCards(cards []Card) {
	d.cards = cards
}

func (d *Deck) RemoveCard(selectedCard Card) {
	for i, card := range d.cards {
		if card.GetName() == selectedCard.GetName() {
			// d.Cards = slices.Delete(d.Cards, i, i+1)
			d.cards = append(d.cards[:i], d.cards[i+1:]...)
			break
		}
	}
}

func (d *Deck) CountCard(selectedCard Card) int {
	count := 0
	for _, card := range d.cards {
		if card.GetName() == selectedCard.GetName() {
			count++
		}
	}
	return count
}

// ////////////
// fonctionnel
// ////////////

var (
	deckList     []Deck
	deckListOnce sync.Once
)

func loadDeckList() {
	// Génère les decks par défaut en mémoire
	count := 20
	cards := make([]Card, 0, count)
	for i := 0; i < count; i++ {
		cards = append(cards, &MonsterCard{name: "Card 1", image: "card1.png", description: "Description 1", level: 1, attack: 100, defense: 50})
	}

	defaultDeck := Deck{
		id:    "default",
		name:  "Default Deck",
		cards: cards,
	}

	deckList = make([]Deck, 0, 30)
	deckList = append(deckList, defaultDeck)
	for i := 0; i < 30; i++ {
		deckList = append(deckList, debug_create_random_deck())
	}
}

func GetDeckList() []Deck {
	deckListOnce.Do(loadDeckList)
	return deckList
}

func GetDeckById(id string) *Deck {
	deckListOnce.Do(loadDeckList)
	for i := range deckList {
		if deckList[i].id == id {
			return &deckList[i]
		}
	}
	return nil
}

func DeleteDeckById(id string) {
	deckListOnce.Do(loadDeckList)
	for i, deck := range deckList {
		if deck.id == id {
			deckList = append(deckList[:i], deckList[i+1:]...)
			return
		}
	}
}

func DuplicateDeckById(id string) {
	deckListOnce.Do(loadDeckList)

	for i, deck := range deckList {
		if deck.id == id {
			var newDeck Deck
			newDeck.id = strconv.Itoa(rand.Intn(1000))
			newDeck.name = "Copy of " + deck.name
			newDeck.cards = append(newDeck.cards, deck.cards...)
			deckList = append(deckList[:i+1], append([]Deck{newDeck}, deckList[i+1:]...)...)
			return
		}
	}
}

func SaveDeck(deck Deck) {
	deckListOnce.Do(loadDeckList)
	for i, d := range deckList {
		if d.id == deck.id {
			deckList[i] = deck
			return
		}
	}
	deckList = append(deckList, deck)
}

// ////////////
// Fonction temporaire pour créer un deck random
// ////////////
func debug_create_random_deck() Deck {
	// retour un deck de 40 cartes
	// avec un nom random

	count := 40
	cards := make([]Card, 0, count)
	for i := 0; i < count; i++ {
		cards = append(cards, debug_get_random_player_card())
	}
	id := "" + strconv.Itoa(rand.Intn(1000))
	randomName := "Deck " + id

	//sort cards by id
	slices.SortFunc(cards, func(a, b Card) int {
		return a.GetId() - b.GetId()
	})

	return Deck{
		id:    id,
		name:  randomName,
		cards: cards,
	}
}
