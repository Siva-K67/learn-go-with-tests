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
