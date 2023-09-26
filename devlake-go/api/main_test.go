package main

import (
	"devlake-go/group-sync/api/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockClient struct {
	weeklyDataPointsToReturn  []models.DataPoint
	weeklyErrorToReturn       error
	monthlyDataPointsToReturn []models.DataPoint
	monthlyErrorToReturn      error
}

func (client MockClient) QueryTotalDeploymentsWeekly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	if client.weeklyDataPointsToReturn != nil {
		return client.weeklyDataPointsToReturn, nil
	}
	if client.weeklyErrorToReturn != nil {
		return nil, client.weeklyErrorToReturn
	}
	return []models.DataPoint{}, nil
}
func (client MockClient) QueryTotalDeploymentsMonthly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	if client.monthlyDataPointsToReturn != nil {
		return client.monthlyDataPointsToReturn, nil
	}
	if client.monthlyErrorToReturn != nil {
		return nil, client.monthlyErrorToReturn
	}
	return []models.DataPoint{}, nil
}

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

			metricHandler(MockClient{})(w, tt.req)

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

// func Test_dfTotalHandler(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		values           url.Values
// 		expectBody       string
// 		expectStatusCode int
// 	}{
// 		{
// 			name:             "should throw 400 response when not specifying a project",
// 			values:           url.Values{},
// 			expectBody:       "project should be provided as a non-empty string\n",
// 			expectStatusCode: 400,
// 		},
// 		{
// 			name:             "should throw 400 response when specifying project without a value",
// 			values:           url.Values{"project": {""}},
// 			expectBody:       "project should be provided as a non-empty string\n",
// 			expectStatusCode: 400,
// 		},
// 		{
// 			name:             "should return data response when specifying a project and default to weekly aggregation",
// 			values:           url.Values{"project": {"my-project"}},
// 			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
// 			expectStatusCode: 200,
// 		},
// 		{
// 			name:             "should return the aggregation from the queries",
// 			values:           url.Values{"project": {"my-project"}, "aggregation": {"weekly"}},
// 			expectBody:       `{"aggregation":"weekly","dataPoints":[]}` + "\n",
// 			expectStatusCode: 200,
// 		},
// 		{
// 			name:             "should return the aggregation from the queries",
// 			values:           url.Values{"project": {"my-project"}, "aggregation": {"monthly"}},
// 			expectBody:       `{"aggregation":"monthly","dataPoints":[]}` + "\n",
// 			expectStatusCode: 200,
// 		},
// 		{
// 			name:             "should return the aggregation from the queries",
// 			values:           url.Values{"project": {"my-project"}, "aggregation": {"quarterly"}},
// 			expectBody:       `{"aggregation":"quarterly","dataPoints":[]}` + "\n",
// 			expectStatusCode: 200,
// 		},
// 		{
// 			name:             "should throw 400 response when specifying an unrecognized aggregation",
// 			values:           url.Values{"project": {"my-project"}, "aggregation": {"not-real-aggregation"}},
// 			expectBody:       "aggregation should be provided as either weekly, monthly or quarterly\n",
// 			expectStatusCode: 400,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()

// 			dfTotalHandler(MockClient{}, w, tt.values)

// 			res := w.Result()
// 			defer res.Body.Close()
// 			data, err := io.ReadAll(res.Body)
// 			if err != nil {
// 				t.Errorf("expected error to be nil got %v", err)
// 			}
// 			if string(data) != tt.expectBody {
// 				t.Errorf("expected '%v' got '%v'", tt.expectBody, string(data))
// 			}
// 			if res.StatusCode != tt.expectStatusCode {
// 				t.Errorf("expected '%v' got '%v'", tt.expectStatusCode, res.StatusCode)
// 			}
// 		})
// 	}
// }

// func Test_dfTotalHandlerDbResults(t *testing.T) {

// 	tests := []struct {
// 		name             string
// 		values           url.Values
// 		mockClient       MockClient
// 		expectBody       string
// 		expectStatusCode int
// 	}{
// 		{
// 			name:             "should throw 500 response with an unexpected error from the database",
// 			values:           url.Values{"project": {"my-project"}, "aggregation": {"weekly"}},
// 			mockClient:       MockClient{weeklyErrorToReturn: fmt.Errorf("error from weekly query"), monthlyErrorToReturn: fmt.Errorf("error from monthly query")},
// 			expectBody:       "error from weekly query\n",
// 			expectStatusCode: 500,
// 		},
// 		{
// 			name:             "should throw 500 response with an unexpected error from the database",
// 			values:           url.Values{"project": {"my-project"}, "aggregation": {"monthly"}},
// 			mockClient:       MockClient{weeklyErrorToReturn: fmt.Errorf("error from weekly query"), monthlyErrorToReturn: fmt.Errorf("error from monthly query")},
// 			expectBody:       "error from monthly query\n",
// 			expectStatusCode: 500,
// 		},
// 		{
// 			name:   "should return weekly data points from the database",
// 			values: url.Values{"project": {"my-project"}, "aggregation": {"weekly"}},
// 			mockClient: MockClient{
// 				weeklyDataPointsToReturn:  []models.DataPoint{{Key: "202338", Value: 0}, {Key: "202337", Value: 1}, {Key: "202336", Value: 2}},
// 				monthlyDataPointsToReturn: []models.DataPoint{{Key: "23/04", Value: 6}, {Key: "23/03", Value: 5}, {Key: "23/02", Value: 4}},
// 			},
// 			expectBody:       `{"aggregation":"weekly","dataPoints":[{"key":"202338","value":0},{"key":"202337","value":1},{"key":"202336","value":2}]}` + "\n",
// 			expectStatusCode: 200,
// 		},
// 		{
// 			name:   "should return monthly data points from the database",
// 			values: url.Values{"project": {"my-project"}, "aggregation": {"monthly"}},
// 			mockClient: MockClient{
// 				weeklyDataPointsToReturn:  []models.DataPoint{{Key: "202338", Value: 0}, {Key: "202337", Value: 1}, {Key: "202336", Value: 2}},
// 				monthlyDataPointsToReturn: []models.DataPoint{{Key: "23/04", Value: 6}, {Key: "23/03", Value: 5}, {Key: "23/02", Value: 4}},
// 			},
// 			expectBody:       `{"aggregation":"monthly","dataPoints":[{"key":"23/04","value":6},{"key":"23/03","value":5},{"key":"23/02","value":4}]}` + "\n",
// 			expectStatusCode: 200,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()

// 			dfTotalHandler(tt.mockClient, w, tt.values)

// 			res := w.Result()
// 			defer res.Body.Close()
// 			data, err := io.ReadAll(res.Body)
// 			if err != nil {
// 				t.Errorf("expected error to be nil got %v", err)
// 			}
// 			if string(data) != tt.expectBody {
// 				t.Errorf("expected '%v' got '%v'", tt.expectBody, string(data))
// 			}
// 			if res.StatusCode != tt.expectStatusCode {
// 				t.Errorf("expected '%v' got '%v'", tt.expectStatusCode, res.StatusCode)
// 			}
// 		})
// 	}
// }
