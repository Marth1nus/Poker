package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ================================
// ================================
// ================================

type Rank rune

const (
	RankNone, Ace, Ten, Jack, Queen, King Rank = ' ', 'A', 'X', 'J', 'Q', 'K'
)

func (r Rank) Validate() error {
	switch r {
	case RankNone, Ace, '2', '3', '4', '5', '6', '7', '8', '9', Ten, Jack, Queen, King:
	default:
		return errors.New("Invalid Rank")
	}
	return nil
}

// ================================
// ================================
// ================================

type Suit rune

const (
	SuitNone, Clovers, Diamonds, Hearts, Spades Suit = ' ', '♣', '♦', '♥', '♠'
)

func (s Suit) Validate() error {
	switch s {
	case SuitNone, Clovers, Diamonds, Hearts, Spades:
	default:
		return errors.New("Invalid Suit")
	}
	return nil
}

// ================================
// ================================
// ================================

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) Validate() error {
	if err := c.Rank.Validate(); err != nil {
		return err
	}
	if err := c.Suit.Validate(); err != nil {
		return err
	}
	return nil
}

func (c Card) String() string {
	return fmt.Sprintf("%c%c", c.Rank, c.Suit)
}

func CardFromString(str string) (Card, error) {
	runes := []rune(str)
	if len(runes) != 2 {
		return Card{}, errors.New("Invalid runes length")
	}
	card := Card{Rank: Rank(runes[0]), Suit: Suit(runes[1])}
	if err := card.Validate(); err != nil {
		return card, err
	}
	return card, nil
}

func (c Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Card) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	card, err := CardFromString(str)
	if err != nil {
		return err
	}
	*c = card
	return nil
}

// ================================
// ================================
// ================================

type Player struct {
	ID    string  `json:"id"`
	Token string  `json:"-"`
	Name  string  `json:"name"`
	Bank  uint64  `json:"bank"`
	Bet   uint64  `json:"bet"`
	Hand  [2]Card `json:"-"`
}

// ================================
// ================================
// ================================

type Game struct {
	ID      string   `json:"id"`
	Players []Player `json:"players"`
	MinBet  uint64   `json:"minBet"`
	Pot     uint64   `json:"pot"`
}

// ================================
// ================================
// ================================

var gGames = map[string]Game{}

// ================================
// ================================
// ================================

func routeGame(w http.ResponseWriter, r *http.Request) {
	gameID := strings.TrimPrefix(r.URL.Path, "/api/game/")

	game, gameFound := gGames[gameID]
	if !gameFound {
		http.Error(w, "Game not Found", http.StatusNotFound)
		return
	}

	responseJson, err := json.Marshal(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	http.HandleFunc("/", routeGame)

	fmt.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
