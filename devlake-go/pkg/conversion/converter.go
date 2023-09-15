package conversion

import (
	"devlake-go/group-sync/pkg/devlake"
	"log"
	"slices"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func BackstageTeamsToDevLakeTeams(backstageTeamMap map[string]backstage.Entity, devLakeTeamMap map[string][]string) {
	var backstageTeamIds []string
	for _, backStageTeam := range backstageTeamMap {
		backstageTeamIds = append(backstageTeamIds, backStageTeam.Metadata.UID)
		id := backstageToDevLakeId(backStageTeam)

		if devLakeTeam, exists := devLakeTeamMap[id]; exists {
			devLakeTeam[devlake.TeamNameColumn] = backStageTeam.Metadata.Name
			log.Printf("Team already exists in DevLake, updating name: %s\n", backStageTeam.Metadata.Name)
		} else {
			devLakeTeamMap[id] = []string{id, backStageTeam.Metadata.Name, "", "", ""}
		}
	}

	removeDeletedBackstageTeams(backstageTeamIds, devLakeTeamMap)
	createRelationships(backstageTeamMap, devLakeTeamMap)
}

func removeDeletedBackstageTeams(backstageTeamIds []string, devLakeTeamMap map[string][]string) {
	for key := range devLakeTeamMap {
		backstageTeamId, found := devLakeToBackstageId(key)
		if !found {
			continue
		}

		if !slices.Contains(backstageTeamIds, backstageTeamId) {
			log.Printf("Team no longer exists in Backstage, removing: %s\n", backstageTeamId)
			delete(devLakeTeamMap, key)
		}
	}
}

func createRelationships(backstageTeamMap map[string]backstage.Entity, devLakeTeamMap map[string][]string) {
	for _, backStageTeam := range backstageTeamMap {
		for _, relation := range backStageTeam.Relations {
			sourceId := backstageToDevLakeId(backStageTeam)
			target, exists := backstageTeamMap[relation.TargetRef]

			if !exists {
				continue
			}

			targetId := backstageToDevLakeId(target)
			if relation.Type == "childOf" {
				devLakeTeamMap[sourceId][devlake.TeamParentIdColumn] = targetId
			} else if relation.Type == "parentOf" {
				devLakeTeamMap[targetId][devlake.TeamParentIdColumn] = sourceId
			}
		}
	}
}
