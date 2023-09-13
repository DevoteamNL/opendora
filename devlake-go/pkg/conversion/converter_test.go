package conversion

import (
	"devlake-go/group-sync/pkg/test"
	"reflect"
	"testing"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func createBackstageTeamWithName(name string) backstage.Entity {
	return backstage.Entity{
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
	}
}

func createTestBackstageTeamsWithNames(names []string) []backstage.Entity {
	var teams []backstage.Entity

	for _, name := range names {
		teams = append(teams, createBackstageTeamWithName(name))
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
			Input{createTestBackstageTeamsWithNames([]string{"TeamA", "TeamB", "TeamC"}), test.ExampleCsvWithColumnHeaders(test.CreateTestDevLakeTeamsWithIds([]string{"3", "5", "4"}))},
			test.ExampleCsvWithColumnHeaders([][]string{{"3", "B", "C", "D", "E"}, {"5", "B", "C", "D", "E"}, {"4", "B", "C", "D", "E"}, {"6", "TeamA", "", "", ""}, {"7", "TeamB", "", "", ""}, {"8", "TeamC", "", "", ""}}),
		},
		{"should skip teams to already in table",
			Input{createTestBackstageTeamsWithNames([]string{"TeamA", "TeamB", "TeamC"}), [][]string{{"6", "TeamA", "", "", ""}, {"7", "TeamB", "", "", ""}}},
			[][]string{{"6", "TeamA", "", "", ""}, {"7", "TeamB", "", "", ""}, {"8", "TeamC", "", "", ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := BackstageTeamsToDevLakeTeams(tt.input.backstageTeams, tt.input.devLakeTeams)
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

type TestRelation struct {
	name         string
	relationType string
}

func createTestBackstageTeamWithRelations(relations []TestRelation) backstage.Entity {
	team := createBackstageTeamWithName("BackstageTeam")

	for _, relation := range relations {
		team.Relations = append(team.Relations, backstage.EntityRelation{
			Type:      relation.relationType,
			TargetRef: "",
			Target: backstage.EntityRelationTarget{
				Name:      relation.name,
				Kind:      "",
				Namespace: "",
			},
		})
	}

	return team
}

func TestCreateRelationshipsTableDriven(t *testing.T) {
	type Input struct {
		backstageTeam backstage.Entity
		devLakeTeams  [][]string
		sourceIndex   int
	}
	var tests = []struct {
		name  string
		input Input
		want  [][]string
	}{
		{
			"should set the parentId of the target for each source parentOf relation",
			Input{
				createTestBackstageTeamWithRelations([]TestRelation{{name: "TeamB", relationType: "parentOf"}, {name: "TeamC", relationType: "parentOf"}}),
				test.CreateTestDevLakeTeamsWithNames([]string{"TeamA", "TeamB", "TeamC"}),
				0,
			},
			[][]string{{"0", "TeamA", "C", "D", "E"}, {"1", "TeamB", "C", "0", "E"}, {"2", "TeamC", "C", "0", "E"}},
		},
		{
			"should set the parentId of the source for each source childOf relation",
			Input{
				createTestBackstageTeamWithRelations([]TestRelation{{name: "TeamC", relationType: "childOf"}}),
				test.CreateTestDevLakeTeamsWithNames([]string{"TeamA", "TeamB", "TeamC"}),
				1,
			},
			[][]string{{"0", "TeamA", "C", "D", "E"}, {"1", "TeamB", "C", "2", "E"}, {"2", "TeamC", "C", "D", "E"}},
		},
		{
			"should not set anything if the target does not exist in the list of DevLake teams",
			Input{
				createTestBackstageTeamWithRelations([]TestRelation{{name: "TeamC", relationType: "parentOf"}}),
				test.CreateTestDevLakeTeamsWithNames([]string{"TeamA", "TeamB"}),
				0,
			},
			[][]string{{"0", "TeamA", "C", "D", "E"}, {"1", "TeamB", "C", "D", "E"}},
		},
		{
			"should not set anything if the target does not exist in the list of DevLake teams",
			Input{
				createTestBackstageTeamWithRelations([]TestRelation{{name: "TeamC", relationType: "childOf"}}),
				test.CreateTestDevLakeTeamsWithNames([]string{"TeamA", "TeamB"}),
				0,
			},
			[][]string{{"0", "TeamA", "C", "D", "E"}, {"1", "TeamB", "C", "D", "E"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createRelationships(tt.input.backstageTeam, tt.input.devLakeTeams, tt.input.sourceIndex)
			ans := tt.input.devLakeTeams
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
