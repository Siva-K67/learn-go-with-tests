package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// FileSystemPlayerStore implements PlayerStore backed by JSON data from a file, cached in memory
type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League // ADDED: cached in-memory copy, loaded once at construction
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

	f.database.Encode(f.league)
}

// initialisePlayerDBFile ensures the file has valid JSON content, writing "[]" if it's empty
func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]")) // seed empty file with valid (empty) JSON array
		file.Seek(0, io.SeekStart)
	}

	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDBFile(file) // CHANGED: replaced inline Seek call with this

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
