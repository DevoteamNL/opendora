package devlake

import (
	"strconv"
	"strings"
)

func TeamNamePredicate(teamName string) func(devLakeTeam []string) bool {
	return func(devLakeTeam []string) bool {
		return strings.EqualFold(devLakeTeam[teamNameColumn], teamName)
	}
}

func LargestTeamId(devLakeTeams [][]string) int {
	latestId := 0
	for _, devLakeTeam := range devLakeTeams {
		idAsInt, err := strconv.Atoi(devLakeTeam[TeamIdColumn])
		if err == nil && latestId < idAsInt {
			latestId = idAsInt
		}
	}
	return latestId
}
