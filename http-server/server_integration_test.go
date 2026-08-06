package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := &PlayerServer{store}
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

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	if response.Code != http.StatusOK {
		t.Errorf("got status %d want %d", response.Code, http.StatusOK)
	}

	if response.Body.String() != "3" {
		t.Errorf("got %q want %q", response.Body.String(), "3")
	}
}
