package backstage

import (
	"context"
	"devlake-go/group-sync/pkg/config"
	"log"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func RetrieveTeams() []backstage.Entity {
	client, err := backstage.NewClient(config.LookupEnvDefault("BACKSTAGE_URL", "http://localhost:7007/"), "default", nil)
	backstageTeams, _, err := client.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{
		Filters: []string{
			"kind=group",
		},
		Fields: []string{},
		Order:  []backstage.ListEntityOrder{{Direction: backstage.OrderDescending, Field: "metadata.name"}},
	})

	if err != nil {
		log.Fatal("Cannot retrieve Backstage teams: ", err)
	}

	return backstageTeams
}
