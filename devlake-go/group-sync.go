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
	"slices"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

func retrieveBackstageTeams() []backstage.Entity {
	client, err := backstage.NewClient("http://localhost:7007/", "default", nil)
	backstageTeams, _, err := client.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{
		Filters: []string{
			"kind=group",
		},
		Fields: []string{},
		Order:  []backstage.ListEntityOrder{{Direction: backstage.OrderDescending, Field: "metadata.name"}},
	})

	if err != nil {
		log.Fatal("Cannot retrieve Backstage teams", err)
	}

	return backstageTeams
}

func retrieveDevLakeTeams() ([][]string, []string) {
	resp, err := http.Get("http://localhost:4000/api/plugins/org/teams.csv")
	if err != nil {
		log.Fatal("Cannot retrieve DevLake teams", err)
	}
	csvReader := csv.NewReader(resp.Body)
	devLakeTeams, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Cannot read DevLake team CSV format", err)
	}

	teamNameIndex := slices.Index(devLakeTeams[0], "Name")
	if teamNameIndex == -1 {
		log.Fatal("DevLake team CSV does not contain a column for team names", err)
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

	if err := csvWriter.WriteAll(devLakeTeams); err != nil {
		log.Fatal("Cannot write DevLake teams to CSV format", err)
	}

	multipartBody := &bytes.Buffer{}
	writer := multipart.NewWriter(multipartBody)
	part, _ := writer.CreateFormFile("file", "teams.csv")

	if _, err := io.Copy(part, buf); err != nil {
		log.Fatal("Cannot copy CSV buffer to multipart file: ", err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal("Cannot close CSV writer: ", err)
	}

	req, err := http.NewRequest("PUT", "http://localhost:4000/api/plugins/org/teams.csv", multipartBody)

	if err != nil {
		log.Fatal("Cannot create DevLake PUT request", err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		log.Fatal("Cannot update DevLake teams CSV", err)
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Cannot read response from DevLake teams update request", err)
	}

	log.Printf("Response: %s\n", resBody)
}

func main() {
	backstageTeams := retrieveBackstageTeams()
	devLakeTeams, devLakeTeamNames := retrieveDevLakeTeams()
	devLakeTeams = appendNewTeams(backstageTeams, devLakeTeams, devLakeTeamNames)

	updateDevLakeTeams(devLakeTeams)
}
