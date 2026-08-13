package main

import (
	"encoding/json"
	"io"
)

// FileSystemPlayerStore implements PlayerStore backed by JSON data from a Reader
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker // CHANGED: was io.Reader — need Seek to rewind before each read
}

// GetLeague reads and parses the current league data from the file
func (f *FileSystemPlayerStore) GetLeague() League {
	f.database.Seek(0, io.SeekStart)   // rewind so we read from the start every time
	league, _ := NewLeague(f.database) // parse JSON into League
	return league
}

// GetPlayerScore returns a named player's win count, or 0 if they don't exist
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.GetLeague().Find(name) // look up player by name

	if player != nil {
		return player.Wins
	}
	return 0 // no such player
}

// RecordWin increments a player's win count and persists the updated league to disk
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name) // pointer into league, or nil

	if player != nil {
		player.Wins++ // mutate through the pointer, updates league in place
	} else {
		league = append(league, Player{name, 1}) //if palyer does not exist in records yet and he wins, then add this player into the record.
	}

	f.database.Seek(0, io.SeekStart)           // rewind before overwriting
	json.NewEncoder(f.database).Encode(league) // write updated league back
}
