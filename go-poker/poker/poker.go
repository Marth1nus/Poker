package poker

import (
	"errors"
	"go-poker/card"
)

type Card = card.Card

type Phase string
type Action string
type Cents int

// ================
// = Types ========
// ================

type Board struct {
	Pot    Cents   `json:"pot"`
	MinBet Cents   `json:"minBet"`
	Round  int     `json:"round"`
	Cards  [5]Card `json:"cards"`

	Deck []Card `json:"-"`
}

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Cards     [2]Card  `json:"cards"`
	Score     int      `json:"score"`
	HandTypes []string `json:"handTypes"`

	BoughtIn bool   `json:"boughtIn"`
	Folded   bool   `json:"folded"`
	Action   Action `json:"action"`

	Bet  Cents `json:"bet"`
	Bank Cents `json:"bank"`
}

type Game struct {
	Id          string   `json:"id"`
	Board       Board    `json:"board"`
	Players     []Player `json:"players"`
	PlayerTurnI int      `json:"playerTurnI"`
	Phase       Phase    `json:"phase"`
}

// ================
// = Phase ========
// ================

var Phases = [4]Phase{"BuyIn", "Check", "Bet", "Winner"}

func (p Phase) Valid() bool {
	switch p {
	case "BuyIn", "Check", "Bet", "Winner":
		return true
	default:
		return false
	}
}

func (p Phase) Next() Phase {
	for i, v := range Phases {
		if v == p {
			return Phases[(i+1)%len(Phases)]
		}
	}
	return Phases[0]
}

func (p Phase) ValidActions() []Action {
	switch p {
	case "BuyIn":
		// None
	case "Check":
		return []Action{"Check", "Fold", "Raise"}
	case "Bet":
		return []Action{"Fold", "Call", "Raise"}
	case "Winner":
		// None
	default:
		// None
	}
	return []Action{}
}

func (p Phase) ValidAction(a Action) bool {
	for _, v := range p.ValidActions() {
		if a == v {
			return true
		}
	}
	return false
}

// ================
// = Action =======
// ================

var Actions = [4]Action{"Check", "Fold", "Call", "Raise"}

func (a Action) Valid() bool {
	switch a {
	case "Check", "Fold", "Call", "Raise":
		return true
	default:
		return false
	}
}

// ================
// = Board ========
// ================

func (b Board) View() Board {
	b.Deck = []Card{}
	return b
}

func (b *Board) PopDeck() Card {
	c := b.Deck[len(b.Deck)-1]
	b.Deck = b.Deck[:len(b.Deck)-1]
	return c
}

func (b *Board) RevealCard() bool {
	for i := 0; i < len(b.Cards); i++ {
		if b.Cards[i] == card.BlankCard {
			b.Cards[i] = b.PopDeck()
			return true
		}
	}
	return false
}

func (b *Board) NewDeck() {
	fullDeck := card.FullDeck()
	b.Deck = card.Shuffle(fullDeck[:])
	b.Cards = [5]Card{card.BlankCard, card.BlankCard, card.BlankCard, card.BlankCard, card.BlankCard}
}

// ================
// = Player =======
// ================

var BlankPlayerHand = [2]Card{card.BlankCard, card.BlankCard}

func (p Player) View(selfPov bool) Player {
	if !selfPov { // Private
		p.Cards = BlankPlayerHand
		p.Score = 0
		p.HandTypes = []string{}
	}
	return p
}

func (p *Player) PlaceBet(bet Cents) bool {
	if bet < p.Bank {
		return false
	}
	p.Bet = bet
	return true
}

// ================
// = Game =========
// ================

func (g Game) View(povPlayerI int) Game {
	g.Board = g.Board.View()
	players := make([]Player, len(g.Players))
	povPlayerJ := povPlayerI
	if !(0 <= povPlayerJ && povPlayerJ < len(g.Players)) {
		povPlayerJ = 0
	}
	for i, j := 0, povPlayerJ; i < len(players); i++ {
		players[i] = g.Players[j].View(j == povPlayerI)
		j = (j + 1) % len(g.Players)
	}
	g.Players = players
	g.PlayerTurnI -= povPlayerI
	if g.PlayerTurnI < 0 {
		g.PlayerTurnI += len(g.Players)
	}
	return g
}

func (g Game) MaxPlayers() int {
	return 8
}

func (g Game) GetPlayerI(playerId string) int {
	for i, v := range g.Players {
		if v.Id == playerId {
			return i
		}
	}
	return -1
}

func (g *Game) Join(player Player) error {
	if len(g.Players) >= g.MaxPlayers() {
		return errors.New("Full")
	}
	if g.GetPlayerI(player.Id) != -1 {
		return errors.New("Present")
	}
	g.Players = append(g.Players, player)
	return nil
}

func (g *Game) Leave(playerId string) error {
	playerI := g.GetPlayerI(playerId)
	if playerI == -1 {
		return errors.New("NotFound")
	}
	if playerI == len(g.Players)-1 {
		g.PlayerTurnI = 0
	} else if playerI < g.PlayerTurnI {
		g.PlayerTurnI--
	}
	g.Players = append(g.Players[:playerI], g.Players[playerI+1:]...)
	return nil
}

func (g *Game) Call(playerId string) error {
	playerI := g.GetPlayerI(playerId)
	if playerI == -1 {
		return errors.New("NotFound")
	}
	player := &g.Players[playerI]
	if player.Folded {
		return errors.New("Folded")
	}
	if player.Bank < g.Board.MinBet {
		return errors.New("Funds")
	}
	player.BoughtIn = true
	player.Bet = g.Board.MinBet
	return nil
}

func (g *Game) Action(playerId string, action Action) error {
	if !g.Phase.ValidAction(action) {
		return errors.New("Phase")
	}

	playerI := g.GetPlayerI(playerId)
	if playerI == -1 {
		return errors.New("NotFound")
	}
	player := &g.Players[playerI]

	if !player.BoughtIn {
		return errors.New("BuyIn")
	}
	if player.Folded {
		return errors.New("Folded")
	}

	switch action {
	case "Check":
		// Move On
	case "Fold":
		player.Folded = true
	case "Call":
		g.Call()
	case "Raise":
	default:
		return errors.New("Invalid")
	}

	player.Action = action
	return nil
}

// ================
// ================

// ================
// ================

// ================
// ================
