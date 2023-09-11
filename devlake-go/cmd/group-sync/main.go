package main

import (
	"devlake-go/group-sync/pkg/backstage"
	"devlake-go/group-sync/pkg/conversion"
	"devlake-go/group-sync/pkg/devlake"
)

func main() {
	backstageTeams := backstage.RetrieveTeams()

	devlakeApiUrl := devlake.TeamsApiUrlFromEnv()
	devLakeTeams := devlake.RetrieveTeams(devlakeApiUrl)
	devLakeTeams = conversion.BackstageTeamsToDevLakeTeams(backstageTeams, devLakeTeams)

	devlake.UpdateTeams(devlakeApiUrl, devLakeTeams)
}
