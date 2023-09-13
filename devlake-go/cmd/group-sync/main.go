package main

import (
	"devlake-go/group-sync/pkg/backstage"
	"devlake-go/group-sync/pkg/config"
	"devlake-go/group-sync/pkg/conversion"
	"devlake-go/group-sync/pkg/devlake"
	"log"
)

func main() {
	backstageTeamMap, err := backstage.RetrieveTeams(config.LookupEnvDefault("BACKSTAGE_URL", "http://localhost:7007/"))
	if err != nil {
		log.Fatal(err)
	}

	devlakeApiUrl := config.LookupEnvDefault("DEVLAKE_URL", "http://localhost:4000/")
	devLakeTeamMap, err := devlake.RetrieveTeams(devlakeApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	conversion.BackstageTeamsToDevLakeTeams(backstageTeamMap, devLakeTeamMap)

	response, err := devlake.UpdateTeams(devlakeApiUrl, devLakeTeamMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %s\n", response)
}
