package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
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
			expectBody:       "type should be provided as one of the following: df_count, df_average, mltc, cfr, mttr\n",
		},
		{
			name:             "should throw 400 response when specifying nonsense metric type",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=not_metric", nil),
			expectBody:       "type should be provided as one of the following: df_count, df_average, mltc, cfr, mttr\n",
			expectStatusCode: 400,
		},
		{
			name:             "should throw 400 response when specifying multiple metric types",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_count&type=df_average", nil),
			expectBody:       "type should be provided as one of the following: df_count, df_average, mltc, cfr, mttr\n",
			expectStatusCode: 400,
		},
		{
			name:             "should return data response when specifying df_count and a project",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?project=my-project&type=df_count", nil),
			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
			expectStatusCode: 200,
		},
		{
			name:             "should return data response when specifying df_average",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_average", nil),
			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
			expectStatusCode: 200,
		},
		{
			name:             "should return data response when specifying mltc",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=mltc", nil),
			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
			expectStatusCode: 200,
		},
		{
			name:             "should return data response when specifying cfr",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=cfr", nil),
			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
			expectStatusCode: 200,
		},
		{
			name:             "should return data response when specifying mttr",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=mttr", nil),
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

func Test_benchmarkHandler(t *testing.T) {
	tests := []struct {
		name             string
		req              *http.Request
		expectBody       string
		expectStatusCode int
	}{
		{
			name:             "should throw 400 response when not specifying metric type",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/benchmark?", nil),
			expectStatusCode: 400,
			expectBody:       "type should be provided as one of the following: df, mltc\n",
		},
		{
			name:             "should return data response when specifying df",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df", nil),
			expectBody:       `{"key":""}` + "\n",
			expectStatusCode: 200,
		},
		{
			name:             "should return data response when specifying mltc",
			req:              httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=mltc", nil),
			expectBody:       `{"key":""}` + "\n",
			expectStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			benchmarkHandler(sql_client.MockClient{})(w, tt.req)

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

func Test_metricDatabaseError(t *testing.T) {
	w := httptest.NewRecorder()
	errorClient := sql_client.MockClient{
		MockDeploymentsDataMap: map[string]sql_client.MockDeploymentsDataReturn{
			sql_queries.WeeklyDeploymentSql + sql_queries.CountSql: {Err: fmt.Errorf("error from weekly query")},
		},
	}

	metricHandler(errorClient)(w, httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_count", nil))

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

func Test_benchmarkDatabaseError(t *testing.T) {
	w := httptest.NewRecorder()
	errorClient := sql_client.MockClient{
		MockBenchmarkDataMap: map[string]sql_client.MockBenchmarkDataReturn{
			sql_queries.BenchmarkDfSql: {Err: fmt.Errorf("error from df benchmark query")},
		},
	}

	benchmarkHandler(errorClient)(w, httptest.NewRequest(http.MethodGet, "/dora/api/benchmark?type=df", nil))

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "error from df benchmark query\n" {
		t.Errorf("expected 'error from df benchmark query' got '%v'", string(data))
	}
	if res.StatusCode != 500 {
		t.Errorf("expected '500' got '%v'", res.StatusCode)
	}
}
