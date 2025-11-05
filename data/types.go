package data

type Card interface {
	GetId() int
	GetName() string
	GetImage() string
	GetDescription() string
}

type MonsterCard struct {
	id          int
	name        string
	image       string
	description string
	level       int
	attack      int
	defense     int
}

type SpellTrapCard struct {
	id          int
	name        string
	image       string
	description string
}

func NewMonsterCard(id int, name, image, description string, level, attack, defense int) *MonsterCard {
	return &MonsterCard{
		id:          id,
		name:        name,
		image:       image,
		description: description,
		level:       level,
		attack:      attack,
		defense:     defense,
	}
}

func NewSpellTrapCard(id int, name, image, description string) *SpellTrapCard {
	return &SpellTrapCard{
		id:          id,
		name:        name,
		image:       image,
		description: description,
	}
}

func (m *MonsterCard) GetId() int {
	return m.id
}

func (m *MonsterCard) GetName() string {
	return m.name
}

func (m *MonsterCard) GetImage() string {
	return m.image
}

func (m *MonsterCard) GetDescription() string {
	return m.description
}

func (s *SpellTrapCard) GetName() string {
	return s.name
}

func (s *SpellTrapCard) GetImage() string {
	return s.image
}

func (s *SpellTrapCard) GetDescription() string {
	return s.description
}

func (s *SpellTrapCard) GetId() int {
	return s.id
}
