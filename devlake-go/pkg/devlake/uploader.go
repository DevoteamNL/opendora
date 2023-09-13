package devlake

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func UpdateTeams(devLakeTeams [][]string) {
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)

	if err := csvWriter.WriteAll(devLakeTeams); err != nil {
		log.Fatal("Cannot write DevLake teams to CSV format: ", err)
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

	req, err := http.NewRequest(http.MethodPut, teamsApiUrlFromEnv(), multipartBody)

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
