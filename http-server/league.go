package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// NewLeague parses JSON player data from a Reader into a league, wrapping any parse error with context
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league) // attempt decode
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err) // add context to raw decode error
	}
	return league, err
}

// League is a collection of Player with lookup behavior attached
type League []Player

// Find returns a pointer to the player with the given name, or nil if not present
func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
