package main

import (
	"encoding/json"
	"io"
)

// FileSystemPlayerStore implements PlayerStore backed by JSON data from a Reader
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker // CHANGED: was io.Reader — need Seek to rewind before each read
}

// GetLeague decodes the JSON in database into a slice of Player
func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, io.SeekStart) // ADDED: move the read cursor back to byte 0
	league, _ := NewLeague(f.database)
	return league
}

// GetPlayerScore returns the win count for a named player, or 0 if not found
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	var wins int
	for _, player := range f.GetLeague() { // reuse GetLeague — it already handles Seek + parsing
		if player.Name == name {
			wins = player.Wins
			break
		}
	}
	return wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
