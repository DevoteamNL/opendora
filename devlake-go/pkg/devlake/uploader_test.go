package devlake

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func exampleTeamsInput() map[string][]string {
	return map[string][]string{
		"1": {"1", "Maple Leafs", "ML", "2", "0"},
		"2": {"2", "Friendly Confines", "FC", "", "1"},
		"3": {"3", "Blue Jays", "BJ", "", "2"},
	}
}

func exampleTeamsCsvLines() []string {
	return []string{"Id,Name,Alias,ParentId,SortingIndex\n", "1,Maple Leafs,ML,2,0\n", "2,Friendly Confines,FC,,1", "3,Blue Jays,BJ,,2\n"}
}

func csvPutHandler(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plugins/org/teams.csv" {
			t.Errorf("Expected to request '/api/plugins/org/teams.csv', got: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected a PUT request, got: %s", r.Method)
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unexpected error reading request body: %v", err)
		}

		for _, line := range exampleTeamsCsvLines() {
			if !strings.Contains(string(reqBody), line) {
				t.Errorf("Expected line in body %s, got: %s", line, string(reqBody))
			}
		}

		fmt.Fprintln(w, "Success")
	})
}

func TestUpdateTeams(t *testing.T) {
	testServer := httptest.NewServer(csvPutHandler(t))
	defer testServer.Close()

	response, err := UpdateTeams(testServer.URL, exampleTeamsInput())
	if err != nil {
		t.Fatalf("unexpected error retrieving teams: %v", err)
	}
	want := "Success\n"

	if string(response) != want {
		t.Errorf("got %s, want %s", response, want)
	}
}

func TestNoServerPutRequest(t *testing.T) {
	response, err := UpdateTeams("http://localhost/no-server", exampleTeamsInput())

	if err == nil || response != nil {
		t.Errorf("Expected no connection to the server to return an error, got: %v", response)
	}
}

func TestErrorResponsePutRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer testServer.Close()

	response, err := UpdateTeams(testServer.URL, exampleTeamsInput())
	if err != nil {
		t.Fatalf("unexpected error retrieving teams: %v", err)
	}
	want := "Not found\n"

	if string(response) != want {
		t.Errorf("got %s, want %s", response, want)
	}
}
