package data

type Deck struct {
	ID    string
	Name  string
	Cards []Card
}

func (d *Deck) RemoveCard(selectedCard Card) {
	for i, card := range d.Cards {
		if card.GetName() == selectedCard.GetName() {
			// d.Cards = slices.Delete(d.Cards, i, i+1)
			d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)
			break
		}
	}
}

func (d *Deck) CountCard(selectedCard Card) int {
	count := 0
	for _, card := range d.Cards {
		if card.GetName() == selectedCard.GetName() {
			count++
		}
	}
	return count
}

type Card interface {
	GetName() string
	GetImage() string
	GetDescription() string
}

type MonsterCard struct {
	ID          string
	Name        string
	Image       string
	Description string
	Level       int
	Attack      int
	Defense     int
}

type SpellTrapCard struct {
	ID          string
	Name        string
	Image       string
	Description string
}

func (m *MonsterCard) GetName() string {
	return m.Name
}
func (m *MonsterCard) GetImage() string {
	return m.Image
}
func (m *MonsterCard) GetDescription() string {
	return m.Description
}
func (s *SpellTrapCard) GetName() string {
	return s.Name
}
func (s *SpellTrapCard) GetImage() string {
	return s.Image
}
func (s *SpellTrapCard) GetDescription() string {
	return s.Description
}
