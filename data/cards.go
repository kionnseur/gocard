package data

import (
	"math/rand"
	"strconv"
	"sync"
)

// Global collection of all available cards.
var (
	allCards    []Card
	playerCards map[int]int // id -> quantity
	initOnce    sync.Once
)

// Init performs one-time initialization of cards.
func init() {
	initOnce.Do(func() {
		loadAllCards()
		loadPlayerCards()
	})
}

// Loads all available cards.
func loadAllCards() {
	allCards = make([]Card, 1000)
	for i := 0; i < 1000; i++ {
		allCards[i] = debugCreateRandomCard(i)
	}
}

// Loads the player's initial random card collection.
func loadPlayerCards() {
	playerCards = make(map[int]int)
	for _, c := range debugGetXRandomCard(200) {
		playerCards[c.GetId()]++
	}
}

// Returns all available cards.
func GetAllCards() []Card {
	return allCards
}

// Returns the player's card collection.
func GetPlayerCards() map[int]int {
	return playerCards
}

// Adds a card to the player's collection.
func AddCardToPlayer(card Card) {
	playerCards[card.GetId()] += 1
}

// Debug function that returns n random cards.
func debugGetXRandomCard(nbr int) []Card {
	cards := make([]Card, nbr)
	for i := 0; i < nbr; i++ {
		cards[i] = GetAllCards()[rand.Intn(len(GetAllCards()))]
	}
	return cards
}

// Debug function that creates a random monster card.
func debugCreateRandomCard(id int) Card {
	randomName := "Card " + strconv.Itoa(id)
	return NewMonsterCard(id, randomName, "card"+strconv.Itoa(id)+".png", "Description "+strconv.Itoa(id), rand.Intn(10)+1, rand.Intn(1000), rand.Intn(1000))
}

// Debug function that returns a random card from the player's collection.
func debugGetRandomPlayerCard() Card {
	keyList := mapKeys(GetPlayerCards())
	size := len(keyList)
	if size == 0 {
		return nil
	}

	randomKey := keyList[rand.Intn(size)]
	return GetAllCards()[randomKey]
}

// Utility function that returns the keys of a map.
func mapKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
