package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"slices"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func lookupEnvDefault(envKey string, envDefaultValue string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return envDefaultValue
}

func retrieveBackstageTeams() []backstage.Entity {
	client, err := backstage.NewClient(lookupEnvDefault("BACKSTAGE_URL", "http://localhost:7007/"), "default", nil)
	backstageTeams, _, err := client.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{
		Filters: []string{
			"kind=group",
		},
		Fields: []string{},
		Order:  []backstage.ListEntityOrder{{Direction: backstage.OrderDescending, Field: "metadata.name"}},
	})

	if err != nil {
		log.Fatal("Cannot retrieve Backstage teams: ", err)
	}

	return backstageTeams
}

func devLakeTeamsApiUrlFromEnv() string {
	return lookupEnvDefault("DEVLAKE_URL", "http://localhost:4000/") + "api/plugins/org/teams.csv"
}

func retrieveDevLakeTeams() ([][]string, []string) {
	if _, ok := os.LookupEnv("REPLACE_DEVLAKE_TEAMS"); ok {
		return [][]string{{"Id", "Name", "Alias", "ParentId", "SortingIndex"}}, []string{}
	}

	resp, err := http.Get(devLakeTeamsApiUrlFromEnv())
	if err != nil {
		log.Fatal("Cannot retrieve DevLake teams: ", err)
	}
	csvReader := csv.NewReader(resp.Body)
	devLakeTeams, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Cannot read DevLake team CSV format: ", err)
	}

	teamNameIndex := slices.Index(devLakeTeams[0], "Name")
	if teamNameIndex == -1 {
		log.Fatal("DevLake team CSV does not contain a column for team names: ", err)
	}

	devLakeTeamNames := []string{}
	for _, team := range devLakeTeams[1:] {
		devLakeTeamNames = append(devLakeTeamNames, team[teamNameIndex])
	}

	return devLakeTeams, devLakeTeamNames
}

func appendNewTeams(backstageTeams []backstage.Entity, devLakeTeams [][]string, devLakeTeamNames []string) [][]string {
	for _, team := range backstageTeams {
		if slices.Contains(devLakeTeamNames, team.Metadata.Name) {
			log.Printf("Team already exists in DevLake: %s\n", team.Metadata.Name)
		} else {
			devLakeTeams = append(devLakeTeams, []string{fmt.Sprint(len(devLakeTeams)), team.Metadata.Name, "", "", ""})
		}
	}
	return devLakeTeams
}

func updateDevLakeTeams(devLakeTeams [][]string) {
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	csvWriter.WriteAll(devLakeTeams)

	if err := csvWriter.Error(); err != nil {
		log.Fatal("Cannot write DevLake teams to CSV format: ", err)
	}

	multipartBody := &bytes.Buffer{}
	writer := multipart.NewWriter(multipartBody)
	part, _ := writer.CreateFormFile("file", "teams.csv")
	io.Copy(part, buf)
	writer.Close()

	req, err := http.NewRequest("PUT", devLakeTeamsApiUrlFromEnv(), multipartBody)

	if err != nil {
		log.Fatal("Cannot create DevLake PUT request: ", err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		log.Fatal("Cannot update DevLake teams CSV: ", err)
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Cannot read response from DevLake teams update request: ", err)
	}

	log.Printf("Response: %s\n", resBody)
}

func main() {
	backstageTeams := retrieveBackstageTeams()
	devLakeTeams, devLakeTeamNames := retrieveDevLakeTeams()
	devLakeTeams = appendNewTeams(backstageTeams, devLakeTeams, devLakeTeamNames)

	updateDevLakeTeams(devLakeTeams)
}
