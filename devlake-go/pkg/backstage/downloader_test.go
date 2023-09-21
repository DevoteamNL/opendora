package backstage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const jsonExample = `
[
	{
		"metadata": {
			"namespace": "default",
			"annotations": {
				"backstage.io/managed-by-location": "dora-backstage-plugin/backstage-plugin/examples/org.yaml",
				"backstage.io/managed-by-origin-location": "dora-backstage-plugin/backstage-plugin/examples/org.yaml"
			},
			"name": "GroupD",
			"uid": "0919d912-b2f2-48df-b935-782bacd29fb2",
			"etag": "aae0f2d4f627e24880b70c50a9c7f4bdde6b8819"
		},
		"apiVersion": "backstage.io/v1alpha1",
		"kind": "Group",
		"spec": {
			"type": "team",
			"children": []
		},
		"relations": [
			{
				"type": "childOf",
				"targetRef": "group:default/groupc",
				"target": {
					"kind": "group",
					"namespace": "default",
					"name": "groupc"
				}
			}
		]
	}
]
`

func backstageGetHandler(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/catalog/entities" {
			t.Errorf("Expected to request '/api/catalog/entities', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected a GET request, got: %s", r.Method)
		}
		fmt.Fprint(w, jsonExample)
	})
}

func TestRetrieveTeams(t *testing.T) {
	testServer := httptest.NewServer(backstageGetHandler(t))
	defer testServer.Close()

	teams, err := RetrieveTeams(testServer.URL)
	if err != nil {
		t.Fatalf("unexpected error retrieving teams: %v", err)
	}

	team, exists := teams["group:default/groupd"]

	if !exists || team.Kind != "Group" || team.Metadata.Name != "GroupD" || team.Metadata.UID != "0919d912-b2f2-48df-b935-782bacd29fb2" {
		t.Errorf("incorrectly retrieved or parsed teams: kind %s; name %s; uid %s", team.Kind, team.Metadata.Name, team.Metadata.UID)
	}
}

func TestNoServerGetRequest(t *testing.T) {
	teams, err := RetrieveTeams("http://localhost/no-server")

	if err == nil || teams != nil {
		t.Errorf("Expected no connection to the server to return an error, got: %v", teams)
	}
}

func TestBadClientUrl(t *testing.T) {
	teams, err := RetrieveTeams(":bad-url")

	if err == nil || teams != nil {
		t.Errorf("Expected bad URL for client to return an error, got: %v", teams)
	}
}
