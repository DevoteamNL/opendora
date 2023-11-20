package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

func TestBenchmarkDfService_ServeRequest(t *testing.T) {

	dataMockMap := map[string]sql_client.MockBenchmarkDataReturn{
		sql_queries.BenchmarkDfSql: {Data: "example_key"},
	}

	errorMockMap := map[string]sql_client.MockBenchmarkDataReturn{
		sql_queries.BenchmarkDfSql: {Err: fmt.Errorf("error from df benchmark query")},
	}

	tests := []struct {
		name           string
		params         ServiceParameters
		mockClient     sql_client.MockClient
		expectResponse models.BenchmarkResponse
		expectError    string
	}{
		{
			name:           "should return an error with an unexpected error from the database",
			params:         ServiceParameters{TypeQuery: "df", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockBenchmarkDataMap: errorMockMap},
			expectResponse: models.BenchmarkResponse{Key: ""},
			expectError:    "error from df benchmark query",
		},
		{
			name:           "should return the df benchmark key from the database",
			params:         ServiceParameters{TypeQuery: "df", Project: "", To: 0, From: 0},
			mockClient:     sql_client.MockClient{MockBenchmarkDataMap: dataMockMap},
			expectResponse: models.BenchmarkResponse{Key: "example_key"},
			expectError:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dfService := BenchmarkDfService{Client: tt.mockClient}
			got, err := dfService.ServeRequest(tt.params)

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !reflect.DeepEqual(got, tt.expectResponse) {
				t.Errorf("BenchmarkDfService.ServeRequest() = %v, want %v", got, tt.expectResponse)
			}
		})
	}
}
