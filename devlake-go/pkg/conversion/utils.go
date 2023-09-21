package conversion

import (
	"strings"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

const backstageTeamIdPrefix = "backstage:"

func backstageToDevLakeId(backStageTeam backstage.Entity) string {
	return backstageTeamIdPrefix + backStageTeam.Metadata.UID
}

func devLakeToBackstageId(key string) (string, bool) {
	return strings.CutPrefix(key, backstageTeamIdPrefix)
}
