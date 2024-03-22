package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

func TestMetricMttrService_ServeRequest(t *testing.T) {

	exampleWeeklyDataPoints := []models.DataPoint{{Key: "202338", Value: 0}, {Key: "202337", Value: 1}, {Key: "202336", Value: 2}}
	exampleMonthlyDataPoints := []models.DataPoint{{Key: "23/04", Value: 6}, {Key: "23/03", Value: 5}, {Key: "23/02", Value: 4}}
	exampleQuarterlyDataPoints := []models.DataPoint{{Key: "2024-01-01", Value: 6}, {Key: "2023-10-01", Value: 5}, {Key: "2023-07-01", Value: 4}}

	dataMockMap := map[string]sql_client.MockDeploymentsDataReturn{
		sql_queries.WeeklyMttrSql:    {Data: exampleWeeklyDataPoints},
		sql_queries.MonthlyMttrSql:   {Data: exampleMonthlyDataPoints},
		sql_queries.QuarterlyMttrSql: {Data: exampleQuarterlyDataPoints}}

	errorMockMap := map[string]sql_client.MockDeploymentsDataReturn{
		sql_queries.WeeklyMttrSql:    {Err: fmt.Errorf("error from weekly query")},
		sql_queries.MonthlyMttrSql:   {Err: fmt.Errorf("error from monthly query")},
		sql_queries.QuarterlyMttrSql: {Err: fmt.Errorf("error from quarterly query")},
	}

	tests := []struct {
		name           string
		params         ServiceParameters
		mockClient     sql_client.MockClient
		expectResponse models.MetricResponse
		expectError    string
	}{
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "mttr", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDeploymentsDataMap: errorMockMap},
			expectResponse: models.MetricResponse{Aggregation: "weekly", DataPoints: nil},
			expectError:    "error from weekly query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "mttr", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDeploymentsDataMap: errorMockMap},
			expectResponse: models.MetricResponse{Aggregation: "monthly", DataPoints: nil},
			expectError:    "error from monthly query",
		},
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "mttr", Aggregation: "quarterly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDeploymentsDataMap: errorMockMap},
			expectResponse: models.MetricResponse{Aggregation: "quarterly", DataPoints: nil},
			expectError:    "error from quarterly query",
		},
		{
			name:           "should return weekly data points from the database",
			params:         ServiceParameters{TypeQuery: "mttr", Aggregation: "weekly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDeploymentsDataMap: dataMockMap},
			expectResponse: models.MetricResponse{Aggregation: "weekly", DataPoints: exampleWeeklyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return monthly data points from the database",
			params:         ServiceParameters{TypeQuery: "mttr", Aggregation: "monthly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDeploymentsDataMap: dataMockMap},
			expectResponse: models.MetricResponse{Aggregation: "monthly", DataPoints: exampleMonthlyDataPoints},
			expectError:    "",
		},
		{
			name:           "should return quarterly data points from the database",
			params:         ServiceParameters{TypeQuery: "mttr", Aggregation: "quarterly", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockDeploymentsDataMap: dataMockMap},
			expectResponse: models.MetricResponse{Aggregation: "quarterly", DataPoints: exampleQuarterlyDataPoints},
			expectError:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mttrService := MetricMttrService{Client: tt.mockClient}
			got, err := mttrService.ServeRequest(tt.params)

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !reflect.DeepEqual(got, tt.expectResponse) {
				t.Errorf("MttrService.ServeRequest() = %v, want %v", got, tt.expectResponse)
			}
		})
	}
}
