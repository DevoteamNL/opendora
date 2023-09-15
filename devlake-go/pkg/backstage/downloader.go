package backstage

import (
	"context"
	"fmt"
	"strings"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func RetrieveTeams(baseUrl string) (teams map[string]backstage.Entity, err error) {
	client, err := backstage.NewClient(baseUrl, "default", nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create Backstage client: %w", err)
	}

	backstageTeams, _, err := client.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{
		Filters: []string{
			"kind=group",
		},
		Fields: []string{},
		Order:  []backstage.ListEntityOrder{{Direction: backstage.OrderDescending, Field: "metadata.name"}},
	})

	if err != nil {
		return nil, fmt.Errorf("cannot retrieve Backstage teams: %w", err)
	}

	backstageTeamMap := make(map[string]backstage.Entity)

	for _, team := range backstageTeams {
		backstageTeamMap[strings.ToLower(fmt.Sprintf("%s:%s/%s", team.Kind, team.Metadata.Namespace, team.Metadata.Name))] = team
	}

	return backstageTeamMap, nil
}
