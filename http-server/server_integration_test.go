package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRecordingWinsAndRetrievingThem drives the real server + real store through actual HTTP requests
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	newPostWinRequest := func(name string) *http.Request {
		req, _ := http.NewRequest(http.MethodPost, "/players/"+name, nil)
		return req
	}
	newGetScoreRequest := func(name string) *http.Request {
		req, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
		return req
	}

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) { // CHANGED: wrapped in subtest
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(t, response.Code, http.StatusOK) // CHANGED: uses shared helper from server_test.go

		if response.Body.String() != "3" {
			t.Errorf("got %q want %q", response.Body.String(), "3")
		}
	})

	t.Run("get league", func(t *testing.T) { // ADDED: new subtest checking /league end-to-end
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest()) // reuses helper from server_test.go
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body) // reuses helper from server_test.go
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
}
