package backstage

import (
	"context"
	"fmt"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func RetrieveTeams(baseUrl string) (teams []backstage.Entity, err error) {
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

	return backstageTeams, nil
}
