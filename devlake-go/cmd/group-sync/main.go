package main

import (
	"devlake-go/group-sync/pkg/backstage"
	"devlake-go/group-sync/pkg/config"
	"devlake-go/group-sync/pkg/conversion"
	"devlake-go/group-sync/pkg/devlake"
	"log"
)

func main() {
	backstageTeams, err := backstage.RetrieveTeams(config.LookupEnvDefault("BACKSTAGE_URL", "http://localhost:7007/"))
	if err != nil {
		log.Fatal(err)
	}

	devlakeApiUrl := config.LookupEnvDefault("DEVLAKE_URL", "http://localhost:4000/")
	devLakeTeams, err := devlake.RetrieveTeams(devlakeApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	devLakeTeams = conversion.BackstageTeamsToDevLakeTeams(backstageTeams, devLakeTeams)

	response, err := devlake.UpdateTeams(devlakeApiUrl, devLakeTeams)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %s\n", response)
}
