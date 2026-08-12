package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
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
	server := &PlayerServer{&store}

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
	server := &PlayerServer{&store}

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

// TestLeague checks the server responds to GET /league with 200 OK
func TestLeague(t *testing.T) {
	store := StubPlayerStore{} // empty stub, don't care about scores yet
	server := &PlayerServer{&store}

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil) // build GET /league request
		response := httptest.NewRecorder()                            // records the response

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK) // check we got 200
	})
}
