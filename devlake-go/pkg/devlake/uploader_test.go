package devlake

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func csvPutHandler(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plugins/org/teams.csv" {
			t.Errorf("Expected to request '/api/plugins/org/teams.csv', got: %s", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Errorf("Expected a PUT request, got: %s", r.Method)
		}

		resBody, _ := io.ReadAll(r.Body)
		fmt.Println(string(resBody))
		fmt.Fprintln(w, "Success")
	})
}

func TestUpdateTeams(t *testing.T) {
	testServer := httptest.NewServer(csvPutHandler(t))
	defer testServer.Close()

	response, err := UpdateTeams(testServer.URL, [][]string{})
	if err != nil {
		t.Fatalf("unexpected error retrieving teams: %v", err)
	}
	want := "Success\n"

	if string(response) != want {
		t.Errorf("got %s, want %s", response, want)
	}
}

func TestNoServerPutRequest(t *testing.T) {
	response, err := UpdateTeams("http://localhost/no-server", [][]string{})

	if err == nil || response != nil {
		t.Errorf("Expected no connection to the server to return an error, got: %v", response)
	}
}

func TestErrorResponsePutRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer testServer.Close()

	response, err := UpdateTeams(testServer.URL, [][]string{})
	if err != nil {
		t.Fatalf("unexpected error retrieving teams: %v", err)
	}
	want := "Not found\n"

	if string(response) != want {
		t.Errorf("got %s, want %s", response, want)
	}
}
