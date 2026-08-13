package main

import (
	"encoding/json"
	"io"
	"os"
)

// FileSystemPlayerStore implements PlayerStore backed by JSON data from a file, cached in memory
type FileSystemPlayerStore struct {
	database io.Writer // CHANGED: was io.ReadWriteSeeker — tape now owns the seeking
	league   League    // ADDED: cached in-memory copy, loaded once at construction
}

// NewFileSystemPlayerStore loads the league once from file and wraps writes in a tape
func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore { // ADDED
	file.Seek(0, io.SeekStart)
	league, _ := NewLeague(file)

	return &FileSystemPlayerStore{
		database: &tape{file}, // wrap the file so writes always start clean
		league:   league,
	}
}

// GetLeague returns the cached league — no disk read needed since we loaded it once at startup
func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league // CHANGED: no more Seek + re-parse from disk every call
}

// GetPlayerScore returns a named player's win count, or 0 if they don't exist
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name) // CHANGED: search cached f.league directly, not f.GetLeague()

	if player != nil {
		return player.Wins
	}
	return 0
}

// RecordWin increments a player's win count in memory and persists the updated league via tape
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name) // CHANGED: operate on f.league directly

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1}) // CHANGED: append to f.league directly
	}

	json.NewEncoder(f.database).Encode(f.league) // CHANGED: no manual Seek — tape handles it
}
