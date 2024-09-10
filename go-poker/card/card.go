package card

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

/*==============
==== Rank ======
==============*/

type Rank int

var Ranks = [RanksLen]string{" ", "A", "2", "3", "4", "5", "6", "7", "8", "9", "X", "J", "Q", "K"}

const (
	RankNull Rank = iota
	RankA
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	Rank9
	RankX
	RankJ
	RankQ
	RankK
	RanksLen
)

func (r Rank) Valid() bool {
	return RankNull <= r && r < RanksLen
}

func (r Rank) String() string {
	if !r.Valid() {
		return "Invalid"
	}
	return Ranks[r]
}

func (r *Rank) LoadString(s string) bool {
	for i, v := range Ranks {
		if s == v {
			*r = Rank(i)
			return true
		}
	}
	*r = Rank(-1)
	return false
}

func RankFromString(s string) Rank {
	r := Rank(-1)
	r.LoadString(s)
	return r
}

func (r Rank) MarshalJSON() ([]byte, error) {
	if !r.Valid() {
		return nil, errors.New("Invalid")
	}
	return json.Marshal(r.String())
}

func (r *Rank) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if !r.LoadString(s) {
		return errors.New("Invalid")
	}
	return nil
}

/*==============
==== Suit ======
==============*/

type Suit int

var Suits = [SuitsLen]string{" ", "♣", "♦", "♥", "♠"}

const (
	SuitNull Suit = iota
	SuitClovers
	SuitDiamonds
	SuitHearts
	SuitSpades
	SuitsLen
)

func (s Suit) Valid() bool {
	return SuitNull <= s && s < SuitsLen
}

func (s Suit) String() string {
	if !s.Valid() {
		return "Invalid"
	}
	return Suits[s]
}

func (s *Suit) LoadString(str string) bool {
	for i, v := range Suits {
		if str == v {
			*s = Suit(i)
			return true
		}
	}
	*s = Suit(-1)
	return false
}

func SuitFromString(str string) Suit {
	s := Suit(-1)
	s.LoadString(str)
	return s
}

func (s Suit) MarshalJSON() ([]byte, error) {
	if !s.Valid() {
		return nil, errors.New("Invalid")
	}
	return json.Marshal(s.String())
}

func (s *Suit) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if !s.LoadString(str) {
		return errors.New("Invalid")
	}
	return nil
}

/*==============
==== Card ======
==============*/

type Card struct {
	Rank Rank
	Suit Suit
}

var BlankCard = Card{RankNull, SuitNull}

const CardsLen = int(RanksLen-1)*int(SuitsLen-1) + 1

func (c Card) Valid() bool {
	return c.Rank.Valid() && c.Suit.Valid()
}

func (c *Card) LoadIndex(i int) {
	i -= 1
	if 0 <= i && i < CardsLen-1 {
		*c = Card{Rank(i >> 2), Suit(i & 0b11)}
	} else {
		*c = BlankCard
	}
}

func (c *Card) LoadString(s string) error {
	if !c.Rank.LoadString(s[:1]) {
		return errors.New("InvalidRank")
	}
	if !c.Suit.LoadString(s[1:]) {
		return errors.New("InvalidSuit")
	}
	return nil
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank.String(), c.Suit.String())
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

/*==============
==== Deck ======
==============*/

func FullDeck() [CardsLen - 1]Card {
	res := [CardsLen - 1]Card{}
	for i := 0; i < CardsLen - 1; i++ {
		res[i].LoadIndex(i + 1)
	}
	return res
}

func Shuffle[T any](slice []T) []T {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
