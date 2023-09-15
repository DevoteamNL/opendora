package conversion

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func createTestDevLakeTeamsWithIds(ids []string) map[string][]string {
	teams := make(map[string][]string)

	for _, id := range ids {
		teams[id] = []string{id, "B", "C", "D", "E"}
	}

	return teams
}

func createBackstageTeam(name string, uid string) backstage.Entity {
	return backstage.Entity{
		ApiVersion: "",
		Kind:       "Group",
		Metadata: backstage.EntityMeta{
			UID:         uid,
			Etag:        "",
			Name:        name,
			Namespace:   "default",
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

func createTestBackstageTeamsWithUids(uids []string) map[string]backstage.Entity {
	teams := make(map[string]backstage.Entity)

	for i, uid := range uids {
		teams["group:default/team"+strconv.Itoa(i)] = createBackstageTeam("Team"+strconv.Itoa(i), uid)
	}

	return teams
}

func TestAppendNewTeamsTableDriven(t *testing.T) {
	type Input struct {
		backstageTeams map[string]backstage.Entity
		devLakeTeams   map[string][]string
	}
	var tests = []struct {
		name  string
		input Input
		want  map[string][]string
	}{
		{"should append nothing to empty table", Input{make(map[string]backstage.Entity), make(map[string][]string)}, make(map[string][]string)},
		{
			"should append a single team to empty table",
			Input{createTestBackstageTeamsWithUids([]string{"uid1"}), make(map[string][]string)},
			map[string][]string{
				"backstage:uid1": {"backstage:uid1", "Team0", "", "", ""},
			},
		},
		{"should append teams to empty table",
			Input{createTestBackstageTeamsWithUids([]string{"uid1", "uid2", "uid3"}), make(map[string][]string)},
			map[string][]string{
				"backstage:uid1": {"backstage:uid1", "Team0", "", "", ""},
				"backstage:uid2": {"backstage:uid2", "Team1", "", "", ""},
				"backstage:uid3": {"backstage:uid3", "Team2", "", "", ""},
			},
		},
		{"should append teams to populated table",
			Input{createTestBackstageTeamsWithUids([]string{"uid1", "uid2", "uid3"}), createTestDevLakeTeamsWithIds([]string{"3", "5", "4"})},
			map[string][]string{
				"3":              {"3", "B", "C", "D", "E"},
				"4":              {"4", "B", "C", "D", "E"},
				"5":              {"5", "B", "C", "D", "E"},
				"backstage:uid1": {"backstage:uid1", "Team0", "", "", ""},
				"backstage:uid2": {"backstage:uid2", "Team1", "", "", ""},
				"backstage:uid3": {"backstage:uid3", "Team2", "", "", ""},
			},
		},
		{"should update the name for teams already in devlake",
			Input{createTestBackstageTeamsWithUids([]string{"uid1", "uid2", "uid3"}), createTestDevLakeTeamsWithIds([]string{"backstage:uid2", "5", "4"})},
			map[string][]string{
				"backstage:uid2": {"backstage:uid2", "Team1", "C", "D", "E"},
				"5":              {"5", "B", "C", "D", "E"},
				"4":              {"4", "B", "C", "D", "E"},
				"backstage:uid1": {"backstage:uid1", "Team0", "", "", ""},
				"backstage:uid3": {"backstage:uid3", "Team2", "", "", ""},
			},
		},
		{"should remove teams that are no-longer in backstage",
			Input{createTestBackstageTeamsWithUids([]string{"uid1", "uid3"}), createTestDevLakeTeamsWithIds([]string{"backstage:uid2", "5", "4"})},
			map[string][]string{
				"5":              {"5", "B", "C", "D", "E"},
				"4":              {"4", "B", "C", "D", "E"},
				"backstage:uid1": {"backstage:uid1", "Team0", "", "", ""},
				"backstage:uid3": {"backstage:uid3", "Team1", "", "", ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BackstageTeamsToDevLakeTeams(tt.input.backstageTeams, tt.input.devLakeTeams)
			ans := tt.input.devLakeTeams
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got:\n %v, want:\n %v", ans, tt.want)
			}
		})
	}
}

type TestRelation struct {
	targetRef    string
	relationType string
}

func createTestBackstageTeamWithRelations(name string, uid string, relations []TestRelation) backstage.Entity {
	team := createBackstageTeam(name, uid)

	for _, relation := range relations {
		team.Relations = append(team.Relations, backstage.EntityRelation{
			Type:      relation.relationType,
			TargetRef: relation.targetRef,
			Target: backstage.EntityRelationTarget{
				Name:      "",
				Kind:      "",
				Namespace: "",
			},
		})
	}

	return team
}

func TestCreateRelationshipsParentOf(t *testing.T) {
	backstageTeams := map[string]backstage.Entity{
		"group:default/teama": createTestBackstageTeamWithRelations("teama", "uid1", []TestRelation{{targetRef: "group:default/teamb", relationType: "parentOf"}, {targetRef: "group:default/teamc", relationType: "parentOf"}}),
		"group:default/teamb": createBackstageTeam("teamb", "uid2"),
		"group:default/teamc": createBackstageTeam("teamc", "uid3"),
	}

	devlakeTeams := createTestDevLakeTeamsWithIds([]string{"backstage:uid1", "backstage:uid2", "backstage:uid3"})

	createRelationships(backstageTeams, devlakeTeams)

	want := map[string][]string{
		"backstage:uid1": {"backstage:uid1", "B", "C", "D", "E"},
		"backstage:uid2": {"backstage:uid2", "B", "C", "backstage:uid1", "E"},
		"backstage:uid3": {"backstage:uid3", "B", "C", "backstage:uid1", "E"},
	}

	if !reflect.DeepEqual(devlakeTeams, want) {
		t.Errorf("got:\n %v, want:\n %v", devlakeTeams, want)
	}
}

func TestCreateRelationshipsChildOf(t *testing.T) {
	backstageTeams := map[string]backstage.Entity{
		"group:default/teama": createTestBackstageTeamWithRelations("teama", "uid1", []TestRelation{{targetRef: "group:default/teamb", relationType: "childOf"}}),
		"group:default/teamb": createTestBackstageTeamWithRelations("teamb", "uid2", []TestRelation{{targetRef: "group:default/teamc", relationType: "childOf"}}),
		"group:default/teamc": createBackstageTeam("teamc", "uid3"),
	}

	devlakeTeams := createTestDevLakeTeamsWithIds([]string{"backstage:uid1", "backstage:uid2", "backstage:uid3"})

	createRelationships(backstageTeams, devlakeTeams)

	want := map[string][]string{
		"backstage:uid1": {"backstage:uid1", "B", "C", "backstage:uid2", "E"},
		"backstage:uid2": {"backstage:uid2", "B", "C", "backstage:uid3", "E"},
		"backstage:uid3": {"backstage:uid3", "B", "C", "D", "E"},
	}

	if !reflect.DeepEqual(devlakeTeams, want) {
		t.Errorf("got:\n %v, want:\n %v", devlakeTeams, want)
	}
}

func TestCreateRelationshipsMissingTarget(t *testing.T) {
	backstageTeams := map[string]backstage.Entity{
		"group:default/teama": createTestBackstageTeamWithRelations("teama", "uid1", []TestRelation{{targetRef: "group:default/team-non-existent", relationType: "parentOf"}}),
		"group:default/teamb": createTestBackstageTeamWithRelations("teamb", "uid2", []TestRelation{{targetRef: "group:default/team-non-existent", relationType: "childOf"}}),
		"group:default/teamc": createBackstageTeam("teamc", "uid3"),
	}

	devlakeTeams := createTestDevLakeTeamsWithIds([]string{"backstage:uid1", "backstage:uid2", "backstage:uid3"})

	createRelationships(backstageTeams, devlakeTeams)

	want := map[string][]string{
		"backstage:uid1": {"backstage:uid1", "B", "C", "D", "E"},
		"backstage:uid2": {"backstage:uid2", "B", "C", "D", "E"},
		"backstage:uid3": {"backstage:uid3", "B", "C", "D", "E"},
	}

	if !reflect.DeepEqual(devlakeTeams, want) {
		t.Errorf("got:\n %v, want:\n %v", devlakeTeams, want)
	}
}

func TestFullConversion(t *testing.T) {
	backstageTeams := map[string]backstage.Entity{
		"group:default/teama": createTestBackstageTeamWithRelations("teama", "uid1", []TestRelation{{targetRef: "group:default/teamb", relationType: "childOf"}}),
		"group:default/teamb": createTestBackstageTeamWithRelations("teamb", "uid2", []TestRelation{{targetRef: "group:default/teamc", relationType: "parentOf"}}),
		"group:default/teamc": createBackstageTeam("teamc", "uid3"),
		"group:default/teame": createTestBackstageTeamWithRelations("teame", "uid5", []TestRelation{{targetRef: "group:default/teamd", relationType: "childOf"}}),
	}

	devlakeTeams := createTestDevLakeTeamsWithIds([]string{"devlake-1", "devlake-2", "backstage:uid3", "backstage:uid4", "backstage:uid5"})

	BackstageTeamsToDevLakeTeams(backstageTeams, devlakeTeams)

	want := map[string][]string{
		"backstage:uid1": {"backstage:uid1", "teama", "", "backstage:uid2", ""},
		"backstage:uid2": {"backstage:uid2", "teamb", "", "", ""},
		"backstage:uid3": {"backstage:uid3", "teamc", "C", "backstage:uid2", "E"},
		"backstage:uid5": {"backstage:uid5", "teame", "C", "D", "E"},
		"devlake-1":      {"devlake-1", "B", "C", "D", "E"},
		"devlake-2":      {"devlake-2", "B", "C", "D", "E"},
	}

	if !reflect.DeepEqual(devlakeTeams, want) {
		t.Errorf("got:\n %v, want:\n %v", devlakeTeams, want)
	}
}
