package conversion

import (
	"devlake-go/group-sync/pkg/devlake"
	"log"
	"slices"
	"strconv"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func BackstageTeamsToDevLakeTeams(backstageTeams []backstage.Entity, devLakeTeams [][]string) [][]string {
	lastId := devlake.LargestTeamId(devLakeTeams)

	for _, backStageTeam := range backstageTeams {
		currentIndex := slices.IndexFunc(devLakeTeams, devlake.TeamNamePredicate(backStageTeam.Metadata.Name))

		if currentIndex != -1 {
			log.Printf("Team already exists in DevLake: %s\n", backStageTeam.Metadata.Name)
		} else {
			lastId += 1
			devLakeTeams = append(devLakeTeams, []string{strconv.Itoa(lastId), backStageTeam.Metadata.Name, "", "", ""})
			currentIndex = len(devLakeTeams) - 1
		}

		createRelationships(backStageTeam, devLakeTeams, currentIndex)
	}
	return devLakeTeams
}

func createRelationships(backStageTeam backstage.Entity, devLakeTeams [][]string, sourceIndex int) {
	for _, relation := range backStageTeam.Relations {
		targetIndex := slices.IndexFunc(devLakeTeams, devlake.TeamNamePredicate(relation.Target.Name))

		if targetIndex == -1 {
			continue
		}
		if relation.Type == "childOf" {
			devLakeTeams[sourceIndex][devlake.TeamParentIdColumn] = devLakeTeams[targetIndex][devlake.TeamIdColumn]
		} else if relation.Type == "parentOf" {
			devLakeTeams[targetIndex][devlake.TeamParentIdColumn] = devLakeTeams[sourceIndex][devlake.TeamIdColumn]
		}
	}
}
