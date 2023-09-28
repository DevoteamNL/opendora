package devlake

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

func RetrieveTeams(baseUrl string, user string, pass string) (teams map[string][]string, err error) {
	req, err := http.NewRequest(http.MethodGet, baseUrl+teamCsvApiPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create DevLake GET request: %w", err)
	}
	req.SetBasicAuth(user, pass)

	var httpClient http.Client
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("cannot retrieve DevLake teams: %w", err)
	}
	csvReader := csv.NewReader(resp.Body)
	devLakeTeams, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read DevLake team CSV format: %w", err)
	}

	devLakeTeamMap := make(map[string][]string)

	for _, team := range devLakeTeams[1:] {
		devLakeTeamMap[team[TeamIdColumn]] = team
	}

	return devLakeTeamMap, nil
}
