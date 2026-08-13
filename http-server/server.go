package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json" // ADDED: avoids repeating this string everywhere

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

// PlayerServer holds dependencies needed to serve requests
type PlayerServer struct {
	store        PlayerStore
	http.Handler // CHANGED: embedded interface instead of named *http.ServeMux field
}

// Player represents a single entry in the league table
type Player struct {
	Name string
	Wins int
}

// NewPlayerServer builds a PlayerServer with routes registered once, not per-request
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer) // creates zero-valued *PlayerServer
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router // ADDED: fills in the embedded http.Handler with our router

	return p
}
func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

// leagueHandler responds to GET /league with the league table as JSON
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType) // CHANGED: was hard-coded string
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

// ADDED: pulled out of the old inline func for /players/
func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}
