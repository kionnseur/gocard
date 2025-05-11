package data

type Deck struct {
	Name  string
	Cards []Card
}
type Card interface {
	GetName() string
	GetImage() string
	GetDescription() string
}

type MonsterCard struct {
	Name        string
	Image       string
	Description string
	Level       int
	Attack      int
	Defense     int
}

type SpellTrapCard struct {
	Name        string
	Image       string
	Description string
}

func (m MonsterCard) GetName() string {
	return m.Name
}
func (m MonsterCard) GetImage() string {
	return m.Image
}
func (m MonsterCard) GetDescription() string {
	return m.Description
}
func (s SpellTrapCard) GetName() string {
	return s.Name
}
func (s SpellTrapCard) GetImage() string {
	return s.Image
}
func (s SpellTrapCard) GetDescription() string {
	return s.Description
}
