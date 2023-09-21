package devlake

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func csvGetHandler(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plugins/org/teams.csv" {
			t.Errorf("Expected to request '/api/plugins/org/teams.csv', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected a GET request, got: %s", r.Method)
		}
		fmt.Fprintln(w, "Id,Name,Alias,ParentId,SortingIndex\n1,Maple Leafs,ML,2,0\n2,Friendly Confines,FC,,1\n3,Blue Jays,BJ,,2")
	})
}

func TestRetrieveTeams(t *testing.T) {
	testServer := httptest.NewServer(csvGetHandler(t))
	defer testServer.Close()

	csv, err := RetrieveTeams(testServer.URL)
	if err != nil {
		t.Fatalf("unexpected error retrieving teams: %v", err)
	}
	want := map[string][]string{
		"1": {"1", "Maple Leafs", "ML", "2", "0"},
		"2": {"2", "Friendly Confines", "FC", "", "1"},
		"3": {"3", "Blue Jays", "BJ", "", "2"},
	}

	if !reflect.DeepEqual(csv, want) {
		t.Errorf("got:\n %v, want:\n %v", csv, want)
	}
}

func TestNoServerGetRequest(t *testing.T) {
	csv, err := RetrieveTeams("http://localhost/no-server")

	if err == nil || csv != nil {
		t.Errorf("Expected no connection to the server to return an error, got: %v", csv)
	}
}
