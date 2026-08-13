package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := NewPlayerServer(&store) // CHANGED: use constructor so router gets built

	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Body.String() != "20" {
			t.Errorf("got %q want %q", response.Body.String(), "20")
		}
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Apollo", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("got status %d want %d", response.Code, http.StatusNotFound)
		}
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{},
	}
	server := NewPlayerServer(&store) // CHANGED: use constructor so router gets built

	t.Run("it records wins on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Code != http.StatusAccepted {
			t.Errorf("got status %d want %d", response.Code, http.StatusAccepted)
		}

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != "Pepper" {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], "Pepper")
		}
	})
}

// TestLeague checks GET /league returns the correct player data as JSON
func TestLeague(t *testing.T) {

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{ // ADDED: the data we expect back
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague} // CHANGED: stub now seeded with wantedLeague
		server := NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(got, wantedLeague) { // ADDED: compare actual vs expected
			t.Errorf("got %v want %v", got, wantedLeague)
		}
	})
}

// GetLeague returns whatever league data the test seeded into the stub
func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}
