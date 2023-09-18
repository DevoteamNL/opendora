package main

import (
	"devlake-go/group-sync/api/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockClient struct {
}

func (client MockClient) QueryTotalDeploymentsWeekly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	return []models.DataPoint{}, nil
}
func (client MockClient) QueryTotalDeploymentsMonthly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	return []models.DataPoint{}, nil
}

func Test_metricHandler(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		expectBody string
	}{
		{
			name:       "should throw 400 response when not specifying metric type",
			req:        httptest.NewRequest(http.MethodGet, "/dora/api/metric?", nil),
			expectBody: "type should be provided as either df_average or df_total\n",
		},
		{
			name:       "should throw 400 response when specifying nonsense metric type",
			req:        httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=not_metric", nil),
			expectBody: "type should be provided as either df_average or df_total\n",
		},
		{
			name:       "should throw 400 response when specifying multiple metric types",
			req:        httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_total&type=df_average", nil),
			expectBody: "type should be provided as either df_average or df_total\n",
		},
		{
			name:       "should return data response when specifying df_total",
			req:        httptest.NewRequest(http.MethodGet, "/dora/api/metric?project=my-project&type=df_total", nil),
			expectBody: `{"aggregation":"weekly","dataPoints":[]}` + "\n",
		},
		{
			name:       "should return todo response when specifying df_average",
			req:        httptest.NewRequest(http.MethodGet, "/dora/api/metric?type=df_average", nil),
			expectBody: "todo average",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			metricHandler(MockClient{})(w, tt.req)

			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectBody {
				t.Errorf("expected '%v' got '%v'", tt.expectBody, string(data))
			}
		})
	}
}
