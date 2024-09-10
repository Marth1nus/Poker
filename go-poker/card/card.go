package card

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

// ================
// ================

type Rank rune

var ValidRanks = [13]Rank{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'X', 'J', 'Q', 'K'}

func (r Rank) Valid() bool {
	for _, v := range ValidRanks {
		if v == r {
			return true
		}
	}
	return false
}

// ================
// ================

type Suit rune

var ValidSuits = [4]Suit{'♣', '♦', '♥', '♠'}

func (s Suit) Valid() bool {
	for _, v := range ValidSuits {
		if v == s {
			return true
		}
	}
	return false
}

// ================
// ================

type Card struct {
	Rank Rank
	Suit Suit
}

var BlankCard = Card{' ', ' '}

func (c Card) Valid() bool {
	return c == BlankCard || c.Rank.Valid() && c.Suit.Valid()
}

func (c *Card) LoadIndex(i int) {
	i -= 1
	if 0 <= i && i < 52 {
		*c = Card{ValidRanks[i>>2], ValidSuits[i&0b11]}
	} else {
		*c = BlankCard
	}
}

func (c *Card) LoadString(s string) error {
	runes := []rune(s)
	*c = Card{}
	if len(runes) != 2 {
		return errors.New("bad string length")
	}
	rank, suit := Rank(runes[0]), Suit(runes[1])
	*c = Card{rank, suit}
	if !c.Valid() {
		return errors.New("card invalid")
	}
	return nil
}

func (c Card) String() string {
	return fmt.Sprintf("%c%c", c.Rank, c.Suit)
}

func (c *Card) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if err := c.LoadString(str); err != nil {
		return err
	}
	return nil
}

func (c Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// ================
// ================

func FullDeck() [52]Card {
	res := [52]Card{}
	for i := 0; i < len(res); i++ {
		res[i].LoadIndex(i + 1)
	}
	return res
}

func Shuffle(cards []Card) []Card {
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i)
		cards[i], cards[j] = cards[j], cards[i]
	}
	return cards
}
