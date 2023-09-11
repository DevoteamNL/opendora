package devlake

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
)

func RetrieveTeams(baseUrl string) (teams [][]string, err error) {
	if _, ok := os.LookupEnv("REPLACE_DEVLAKE_TEAMS"); ok {
		return [][]string{{"Id", "Name", "Alias", "ParentId", "SortingIndex"}}, nil
	}

	resp, err := http.Get(baseUrl + teamCsvApiPath)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve DevLake teams: %w", err)
	}
	csvReader := csv.NewReader(resp.Body)
	devLakeTeams, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read DevLake team CSV format: %w", err)
	}

	return devLakeTeams, nil
}
