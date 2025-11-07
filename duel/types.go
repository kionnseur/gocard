package duel

type Player struct {
	Name             string
	LifePoints       uint32
	InvocationPower  uint8
	InvocationNumber uint8
	SpellTrapSet     uint8
	Deck             int
}

type Duel struct {
	LeftPlayer  Player
	RightPlayer Player
	IsPaused    bool
	Timer       int
}
