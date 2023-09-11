package devlake

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
)

func RetrieveTeams(baseUrl string) [][]string {
	if _, ok := os.LookupEnv("REPLACE_DEVLAKE_TEAMS"); ok {
		return [][]string{{"Id", "Name", "Alias", "ParentId", "SortingIndex"}}
	}

	resp, err := http.Get(baseUrl + teamCsvApiPath)
	if err != nil {
		log.Fatal("Cannot retrieve DevLake teams: ", err)
	}
	csvReader := csv.NewReader(resp.Body)
	devLakeTeams, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Cannot read DevLake team CSV format: ", err)
	}

	return devLakeTeams
}
