package devlake

import (
	"devlake-go/group-sync/pkg/test"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRetrieveTeams(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plugins/org/teams.csv" {
			t.Errorf("Expected to request '/api/plugins/org/teams.csv'/, got: %s", r.URL.Path)
		}
		fmt.Fprintln(w, "Id,Name,Alias,ParentId,SortingIndex\n1,Maple Leafs,ML,2,0\n2,Friendly Confines,FC,,1\n3,Blue Jays,BJ,,2")
	}))

	defer testServer.Close()

	csv := RetrieveTeams(testServer.URL)
	want := test.ExampleCsvWithColumnHeaders([][]string{{"1", "Maple Leafs", "ML", "2", "0"}, {"2", "Friendly Confines", "FC", "", "1"}, {"3", "Blue Jays", "BJ", "", "2"}})

	if !reflect.DeepEqual(csv, want) {
		t.Errorf("got %v, want %v", csv, want)
	}
}
