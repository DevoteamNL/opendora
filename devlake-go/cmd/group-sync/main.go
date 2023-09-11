package main

import (
	"devlake-go/group-sync/pkg/backstage"
	"devlake-go/group-sync/pkg/conversion"
	"devlake-go/group-sync/pkg/devlake"
	"log"
)

func main() {
	backstageTeams := backstage.RetrieveTeams()

	devlakeApiUrl := devlake.ApiUrlFromEnv()
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
