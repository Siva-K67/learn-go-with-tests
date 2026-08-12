package main

import (
	"log"
	"net/http"
	"sync" // ADDED: needed for Mutex
)

type InMemoryPlayerStore struct {
	store map[string]int
	lock  sync.Mutex // ADDED: guards `store` against concurrent access
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{store: map[string]int{}}
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	i.lock.Lock()         // ADDED
	defer i.lock.Unlock() // ADDED
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.lock.Lock()         // ADDED
	defer i.lock.Unlock() // ADDED
	i.store[name]++
}

func main() {
	server := NewPlayerServer(NewInMemoryPlayerStore()) // CHANGED: use constructor so router gets built
	log.Fatal(http.ListenAndServe(":5000", server))     // UNCHANGED
}
