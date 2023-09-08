package main

import (
	"devlake-go/group-sync/pkg/backstage"
	"devlake-go/group-sync/pkg/conversion"
	"devlake-go/group-sync/pkg/devlake"
)

func main() {
	backstageTeams := backstage.RetrieveTeams()
	devLakeTeams := devlake.RetrieveTeams()
	devLakeTeams = conversion.BackstageTeamsToDevLakeTeams(backstageTeams, devLakeTeams)

	devlake.UpdateTeams(devLakeTeams)
}
