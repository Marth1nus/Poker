package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Suit rune
type Rank rune
type Card struct { 
	Rank Rank
	Suit Suit
}

type Player struct {
	ID    string
	Token string
}

type Game struct {
	ID      string
	Players []Player
}

var g_Games = map[string]Game{

}

func routeGame(w http.ResponseWriter, r *http.Request) {
	gameID := strings.TrimPrefix(r.URL.Path, "/api/game/")

	found := gameID != "null"
	if !found {
		http.Error(w, "Game not Found", http.StatusNotFound)
		return
	}

	responseBody := map[string]string{
		"gameID": gameID,
	}

	responseJson, err := json.Marshal(responseBody)
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

	http.HandleFunc("/api/game/", routeGame)

	fmt.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
