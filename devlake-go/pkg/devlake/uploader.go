package devlake

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func csvHeader() []string { return []string{"Id", "Name", "Alias", "ParentId", "SortingIndex"} }

func teamMapToCsvSlice(teamMap map[string][]string) [][]string {
	entries := [][]string{csvHeader()}
	for _, entry := range teamMap {
		entries = append(entries, entry)
	}
	return entries
}

func UpdateTeams(baseUrl string, devLakeTeamMap map[string][]string) (response []byte, err error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)

	if err := csvWriter.WriteAll(teamMapToCsvSlice(devLakeTeamMap)); err != nil {
		return nil, fmt.Errorf("cannot write DevLake teams to CSV format: %w", err)
	}

	var multipartBody bytes.Buffer
	writer := multipart.NewWriter(&multipartBody)
	part, err := writer.CreateFormFile("file", "teams.csv")
	if err != nil {
		return nil, fmt.Errorf("cannot create multipart file: %w", err)
	}
	if _, err := io.Copy(part, &buf); err != nil {
		return nil, fmt.Errorf("cannot copy CSV buffer to multipart file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("cannot close CSV writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseUrl+teamCsvApiPath, &multipartBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create DevLake PUT request: %w", err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	var httpClient http.Client
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("cannot update DevLake teams CSV: %w", err)
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response from DevLake teams update request: %w", err)
	}

	return resBody, nil
}
