package main

import (
	"devlake-go/group-sync/pkg/backstage"
	"devlake-go/group-sync/pkg/conversion"
	"devlake-go/group-sync/pkg/devlake"
	"log"
)

func main() {
	backstageTeams := backstage.RetrieveTeams()

	devlakeApiUrl := devlake.TeamsApiUrlFromEnv()
	devLakeTeams, err := devlake.RetrieveTeams(devlakeApiUrl)
	if err != nil {
		log.Fatal(err)
	}
	devLakeTeams = conversion.BackstageTeamsToDevLakeTeams(backstageTeams, devLakeTeams)

	devlake.UpdateTeams(devlakeApiUrl, devLakeTeams)
}
