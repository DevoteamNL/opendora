package main

import (
	"reflect"
	"testing"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func exampleCsvWithColumnHeaders(otherRows [][]string) [][]string {
	return append([][]string{{"Id", "Name", "Alias", "ParentId", "SortingIndex"}}, otherRows...)
}

func createTestDevLakeTeamsWithIds(ids []string) [][]string {
	teams := [][]string{}

	for _, id := range ids {
		teams = append(teams, []string{id, "B", "C", "D", "E"})
	}

	return teams
}

func TestLargestTeamIdTableDriven(t *testing.T) {

	var tests = []struct {
		name  string
		input [][]string
		want  int
	}{
		{"should use 0 for empty table", [][]string{}, 0},
		{"should use 0 for table without numbers in the first column", exampleCsvWithColumnHeaders([][]string{}), 0},
		{"should use 0 for table without numbers in the first column", exampleCsvWithColumnHeaders(createTestDevLakeTeamsWithIds([]string{"A"})), 0},
		{"should use the largest number found in the first column", createTestDevLakeTeamsWithIds([]string{"3", "4", "5"}), 5},
		{"should use the largest number found in the first column", createTestDevLakeTeamsWithIds([]string{"50", "40", "30"}), 50},
		{"should use the largest number found in the first column", createTestDevLakeTeamsWithIds([]string{"40", "A", "30"}), 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := largestTeamId(tt.input)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func createTestBackstageTeamsWithNames(names []string) []backstage.Entity {
	teams := []backstage.Entity{}

	for _, name := range names {
		teams = append(teams, backstage.Entity{
			ApiVersion: "",
			Kind:       "",
			Metadata: backstage.EntityMeta{
				UID:         "",
				Etag:        "",
				Name:        name,
				Namespace:   "",
				Title:       "",
				Description: "",
				Labels:      map[string]string{},
				Annotations: map[string]string{},
				Tags:        []string{},
				Links:       []backstage.EntityLink{},
			},
			Spec:      map[string]interface{}{},
			Relations: []backstage.EntityRelation{},
		})
	}

	return teams
}

func TestAppendNewTeamsTableDriven(t *testing.T) {
	type Input struct {
		backstageTeams []backstage.Entity
		devLakeTeams   [][]string
	}
	var tests = []struct {
		name  string
		input Input
		want  [][]string
	}{
		{"should append nothing to empty table", Input{[]backstage.Entity{}, [][]string{}}, [][]string{}},
		{
			"should append a single team to empty table",
			Input{createTestBackstageTeamsWithNames([]string{"TeamA"}), [][]string{}},
			[][]string{{"1", "TeamA", "", "", ""}},
		},
		{"should append teams to empty table",
			Input{createTestBackstageTeamsWithNames([]string{"TeamA", "TeamB", "TeamC"}), [][]string{}},
			[][]string{{"1", "TeamA", "", "", ""}, {"2", "TeamB", "", "", ""}, {"3", "TeamC", "", "", ""}},
		},
		{"should append teams to populated table",
			Input{createTestBackstageTeamsWithNames([]string{"TeamA", "TeamB", "TeamC"}), exampleCsvWithColumnHeaders(createTestDevLakeTeamsWithIds([]string{"3", "5", "4"}))},
			exampleCsvWithColumnHeaders([][]string{{"3", "B", "C", "D", "E"}, {"5", "B", "C", "D", "E"}, {"4", "B", "C", "D", "E"}, {"6", "TeamA", "", "", ""}, {"7", "TeamB", "", "", ""}, {"8", "TeamC", "", "", ""}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := appendNewTeams(tt.input.backstageTeams, tt.input.devLakeTeams)
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
