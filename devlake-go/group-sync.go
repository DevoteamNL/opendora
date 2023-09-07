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
	"strings"

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

func retrieveDevLakeTeams() [][]string {
	resp, err := http.Get("http://localhost:4000/api/plugins/org/teams.csv")
	if err != nil {
		log.Fatal("Cannot retrieve DevLake teams", err)
	}
	csvReader := csv.NewReader(resp.Body)
	devLakeTeams, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Cannot read DevLake team CSV format", err)
	}

	return devLakeTeams
}

func devLakeTeamFinder(teamName string) func(devLakeTeam []string) bool {
	const DevLakeTeamNameColumn = 1

	return func(devLakeTeam []string) bool {
		return strings.EqualFold(devLakeTeam[DevLakeTeamNameColumn], teamName)
	}
}

func appendNewTeams(backstageTeams []backstage.Entity, devLakeTeams [][]string) [][]string {
	const DevLakeTeamIdColumn = 0
	const DevLakeTeamParentIdColumn = 3

	for _, backStageTeam := range backstageTeams {
		currentIndex := slices.IndexFunc(devLakeTeams, devLakeTeamFinder(backStageTeam.Metadata.Name))
		log.Printf("Team %v %s: %v\n", currentIndex, backStageTeam.Metadata.Name, backStageTeam.Spec["children"])
		if currentIndex != -1 {
			log.Printf("Team already exists in DevLake: %s\n", backStageTeam.Metadata.Name)
		} else {
			devLakeTeams = append(devLakeTeams, []string{fmt.Sprint(len(devLakeTeams)), backStageTeam.Metadata.Name, "", "", ""})
			currentIndex = len(devLakeTeams) - 1
		}

		for _, relation := range backStageTeam.Relations {
			targetIndex := slices.IndexFunc(devLakeTeams, devLakeTeamFinder(relation.Target.Name))

			if targetIndex == -1 {
				continue
			}
			if relation.Type == "childOf" {
				devLakeTeams[currentIndex][DevLakeTeamParentIdColumn] = devLakeTeams[targetIndex][DevLakeTeamIdColumn]
			} else if relation.Type == "parentOf" {
				devLakeTeams[targetIndex][DevLakeTeamParentIdColumn] = devLakeTeams[currentIndex][DevLakeTeamIdColumn]
			}
		}

	}
	return devLakeTeams
}

func updateDevLakeTeams(devLakeTeams [][]string) {
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	csvWriter.WriteAll(devLakeTeams)

	if err := csvWriter.Error(); err != nil {
		log.Fatal("Cannot write DevLake teams to CSV format", err)
	}

	multipartBody := &bytes.Buffer{}
	writer := multipart.NewWriter(multipartBody)
	part, _ := writer.CreateFormFile("file", "teams.csv")
	io.Copy(part, buf)
	writer.Close()

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
	devLakeTeams := retrieveDevLakeTeams()
	devLakeTeams = appendNewTeams(backstageTeams, devLakeTeams)

	updateDevLakeTeams(devLakeTeams)
}
