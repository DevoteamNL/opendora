package conversion

import (
	"devlake-go/group-sync/pkg/devlake"
	"log"
	"slices"
	"strings"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func BackstageTeamsToDevLakeTeams(backstageTeamMap map[string]backstage.Entity, devLakeTeamMap map[string][]string) {
	var backstageTeamIds []string
	for _, backStageTeam := range backstageTeamMap {
		backstageTeamIds = append(backstageTeamIds, backStageTeam.Metadata.UID)
		id := devlake.BackstageTeamIdPrefix + backStageTeam.Metadata.UID
		devLakeTeam, exists := devLakeTeamMap[id]

		if exists {
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
		backstageTeamId, found := strings.CutPrefix(key, devlake.BackstageTeamIdPrefix)
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
			sourceId := devlake.BackstageTeamIdPrefix + backStageTeam.Metadata.UID
			target, exists := backstageTeamMap[relation.TargetRef]

			if !exists {
				continue
			}

			targetId := devlake.BackstageTeamIdPrefix + target.Metadata.UID
			if relation.Type == "childOf" {
				devLakeTeamMap[sourceId][devlake.TeamParentIdColumn] = targetId
			} else if relation.Type == "parentOf" {
				devLakeTeamMap[targetId][devlake.TeamParentIdColumn] = sourceId
			}
		}
	}
}
