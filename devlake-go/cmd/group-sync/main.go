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

	devLakeApiUrl := config.LookupEnvDefault("DEVLAKE_URL", "http://localhost:4000/")
	devLakeAdminUser := config.LookupEnvDefault("DEVLAKE_ADMIN_USER", "devlake")
	devLakeAdminPass := config.LookupEnvDefault("DEVLAKE_ADMIN_PASS", "merico")
	devLakeTeamMap, err := devlake.RetrieveTeams(devLakeApiUrl, devLakeAdminUser, devLakeAdminPass)
	if err != nil {
		log.Fatal(err)
	}

	conversion.BackstageTeamsToDevLakeTeams(backstageTeamMap, devLakeTeamMap)

	response, err := devlake.UpdateTeams(devLakeApiUrl, devLakeAdminUser, devLakeAdminPass, devLakeTeamMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %s\n", response)
}
