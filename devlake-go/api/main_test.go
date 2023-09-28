package main

import (
	"devlake-go/group-sync/api/sql_client"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_metricHandler(t *testing.T) {
	tests := []struct {
		name             string
		req              *http.Request
		expectBody       string
		expectStatusCode int
	}{
		{
			name:             "should throw 400 response when not specifying metric type",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?", nil),
			expectStatusCode: 400,
			expectBody:       "type should be provided as either df_average or df_total\n",
		},
		{
			name:             "should throw 400 response when specifying nonsense metric type",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=not_metric", nil),
			expectBody:       "type should be provided as either df_average or df_total\n",
			expectStatusCode: 400,
		},
		{
			name:             "should throw 400 response when specifying multiple metric types",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_total&type=df_average", nil),
			expectBody:       "type should be provided as either df_average or df_total\n",
			expectStatusCode: 400,
		},
		{
			name:             "should return data response when specifying df_total and a project",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?project=my-project&type=df_total", nil),
			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
			expectStatusCode: 200,
		},
		{
			name:             "should return data response when specifying df_average",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_average", nil),
			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
			expectStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			metricHandler(sql_client.MockClient{})(w, tt.req)

			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectBody {
				t.Errorf("expected '%v' got '%v'", tt.expectBody, string(data))
			}
			if res.StatusCode != tt.expectStatusCode {
				t.Errorf("expected '%v' got '%v'", tt.expectStatusCode, res.StatusCode)
			}
		})
	}
}

func Test_databaseError(t *testing.T) {
	w := httptest.NewRecorder()
	errorClient := sql_client.MockClient{
		MockDataMap: map[string]sql_client.MockDataReturn{
			sql_client.WEEKLY_DEPLOYMENT_SQL: {Err: fmt.Errorf("error from weekly query")},
		},
	}

	metricHandler(errorClient)(w, httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_total", nil))

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "error from weekly query\n" {
		t.Errorf("expected 'error from weekly query' got '%v'", string(data))
	}
	if res.StatusCode != 500 {
		t.Errorf("expected '500' got '%v'", res.StatusCode)
	}
}
