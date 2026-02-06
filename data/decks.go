package data

import (
	"slices"
	"strconv"
	"sync"
)

// ID generator for unique deck IDs.
var _ID = 1

// Returns a unique string ID.
func ID() string {
	i := _ID
	_ID++
	return strconv.Itoa(i)
}

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

func (d *Deck) SetId(id string) {
	d.id = id
}

func (d *Deck) GetCards() []Card {
	return d.cards
}

func (d *Deck) SetCards(cards []Card) {
	d.cards = cards
}

// Removes a specific card from the deck.
func (d *Deck) RemoveCard(selectedCard Card) {
	for i, card := range d.cards {
		if card.GetName() == selectedCard.GetName() {
			d.cards = append(d.cards[:i], d.cards[i+1:]...)
			break
		}
	}
}

// Counts the occurrences of a card in the deck.
func (d *Deck) CountCard(selectedCard Card) int {
	if d == nil {
		return 0
	}
	count := 0
	for _, card := range d.cards {
		if card.GetName() == selectedCard.GetName() {
			count++
		}
	}
	return count
}

// ////////////
// Operational functions
// ////////////

// Global deck list.
var (
	deckList     []Deck
	deckListOnce sync.Once
)

// Loads the initial deck list (with default and random decks).
func loadDeckList() {
	count := 20
	cards := make([]Card, 0, count)
	for i := 0; i < count; i++ {
		cards = append(cards, &MonsterCard{name: "Card 1", image: "card1.png", description: "Description 1", level: 1, attack: 100, defense: 50})
	}

	defaultDeck := Deck{
		id:    ID(),
		name:  "Default Deck",
		cards: cards,
	}

	deckList = make([]Deck, 0, 30)
	deckList = append(deckList, defaultDeck)
	for i := 0; i < 30; i++ {
		deckList = append(deckList, debugCreateRandomDeck())
	}
}

// Returns all decks.
func GetDeckList() []Deck {
	deckListOnce.Do(loadDeckList)
	return deckList
}

// Returns a deck by ID.
func GetDeckById(id string) *Deck {
	deckListOnce.Do(loadDeckList)
	for i := range deckList {
		if deckList[i].id == id {
			return &deckList[i]
		}
	}
	return nil
}

// Deletes a deck by ID.
func DeleteDeckById(id string) {
	deckListOnce.Do(loadDeckList)
	for i, deck := range deckList {
		if deck.id == id {
			deckList = append(deckList[:i], deckList[i+1:]...)
			return
		}
	}
}

// Clones an existing deck for editing.
func CloneDeckById(id string) Deck {
	deckListOnce.Do(loadDeckList)
	if GetDeckById(id) == nil {
		// Return empty deck without ID - ID will be assigned on save
		return Deck{name: "New Deck", cards: []Card{}}
	}
	var toCopy Deck = *GetDeckById(id)
	var newDeck Deck
	newDeck.id = toCopy.id
	newDeck.name = toCopy.name
	newDeck.cards = append(newDeck.cards, toCopy.cards...)
	return newDeck
}

// Duplicates a deck with a new name.
func DuplicateDeckById(id string) {
	deckListOnce.Do(loadDeckList)

	for i, deck := range deckList {
		if deck.id == id {
			var newDeck Deck
			newDeck.id = ID()
			newDeck.name = "Copy of " + deck.name
			newDeck.cards = append(newDeck.cards, deck.cards...)
			deckList = append(deckList[:i+1], append([]Deck{newDeck}, deckList[i+1:]...)...)
		}
	}
}

// Saves or updates a deck.
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
// Temporary function to create a random deck
// ////////////
func debugCreateRandomDeck() Deck {
	count := 40
	cards := make([]Card, 0, count)
	for i := 0; i < count; i++ {
		cards = append(cards, debugGetRandomPlayerCard())
	}
	id := ID()
	randomName := "Deck " + id

	// Sort cards by ID
	slices.SortFunc(cards, func(a, b Card) int {
		return a.GetId() - b.GetId()
	})

	return Deck{
		id:    id,
		name:  randomName,
		cards: cards,
	}
}
