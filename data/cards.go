package data

import (
	"math/rand"
	"strconv"
	"sync"
)

var (
	allCards    []Card
	playerCards map[int]int
	initOnce    sync.Once
)

func init() {
	initOnce.Do(func() {
		loadAllCards()
		loadPlayerCards()
	})
}

func loadAllCards() {
	allCards = make([]Card, 0, 1000)
	for i := 0; i < 1000; i++ {
		allCards = append(allCards, debug_create_random_card())
	}
}

func loadPlayerCards() {
	playerCards = make(map[int]int)
	for _, c := range debug_get_x_random_card(200) {
		playerCards[c.GetId()]++
	}
}

func GetAllCards() []Card {
	return allCards
}

func GetPlayerCards() map[int]int {
	return playerCards
}

func AddCardToPlayer(card Card) {
	playerCards[card.GetId()] += 1
}

func debug_get_x_random_card(nbr int) []Card {
	// retourne nbr carte random
	cards := make([]Card, nbr)
	for i := 0; i < nbr; i++ {
		cards[i] = GetAllCards()[rand.Intn(len(GetAllCards()))]
	}
	return cards
}

func debug_create_random_card() Card {
	// créé une carte random
	// avec un nom random
	id := rand.Intn(1000)
	randomName := "Card " + strconv.Itoa(id)
	return &MonsterCard{
		ID:          id,
		Name:        randomName,
		Image:       "card" + strconv.Itoa(id) + ".png",
		Description: "Description " + strconv.Itoa(id),
		Level:       rand.Intn(10) + 1,
		Attack:      rand.Intn(1000),
		Defense:     rand.Intn(1000),
	}
}
