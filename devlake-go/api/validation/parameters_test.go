package validation

import (
	"math"
	"net/url"
	"testing"
	"time"

	"github.com/devoteamnl/opendora/api/service"
)

func Test_validTypeQuery(t *testing.T) {
	tests := []struct {
		name            string
		values          url.Values
		expectTypeQuery string
		expectValid     bool
	}{
		{
			name:            "should not be valid when no type is provided",
			values:          url.Values{},
			expectTypeQuery: "",
			expectValid:     false,
		},
		{
			name:            "should not be valid when no type is provided",
			values:          url.Values{"type": {}},
			expectTypeQuery: "",
			expectValid:     false,
		},
		{
			name:            "should not be valid when multiple types are provided",
			values:          url.Values{"type": {"df_count", "df_average"}},
			expectTypeQuery: "",
			expectValid:     false,
		},
		{
			name:            "should not be valid when an unrecognized type is provided",
			values:          url.Values{"type": {"not_a_type"}},
			expectTypeQuery: "",
			expectValid:     false,
		},
		{
			name:            "should be valid for df_count",
			values:          url.Values{"type": {"df_count"}},
			expectTypeQuery: "df_count",
			expectValid:     true,
		},
		{
			name:            "should be valid for df_average",
			values:          url.Values{"type": {"df_average"}},
			expectTypeQuery: "df_average",
			expectValid:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeQuery, valid := validTypeQuery(tt.values, []string{"df_average", "df_count"})

			if valid != tt.expectValid {
				t.Errorf("expected '%t' got '%t'", tt.expectValid, valid)
			}
			if tt.expectValid && typeQuery != tt.expectTypeQuery {
				t.Errorf("expected '%v' got '%v'", tt.expectTypeQuery, typeQuery)
			}
		})
	}
}

func Test_validAggregationQuery(t *testing.T) {
	tests := []struct {
		name              string
		values            url.Values
		expectAggregation string
		expectValid       bool
	}{
		{
			name:              "should default to weekly when no aggregation is provided",
			values:            url.Values{},
			expectAggregation: "weekly",
			expectValid:       true,
		},
		{
			name:              "should default to weekly when no aggregation is provided",
			values:            url.Values{"aggregation": {}},
			expectAggregation: "weekly",
			expectValid:       true,
		},
		{
			name:              "should not be valid when multiple aggregations are provided",
			values:            url.Values{"aggregation": {"weekly", "monthly"}},
			expectAggregation: "",
			expectValid:       false,
		},
		{
			name:              "should not be valid when an unrecognized aggregation is provided",
			values:            url.Values{"aggregation": {"not_an_aggregation"}},
			expectAggregation: "",
			expectValid:       false,
		},
		{
			name:              "should be valid for weekly aggregation",
			values:            url.Values{"aggregation": {"weekly"}},
			expectAggregation: "weekly",
			expectValid:       true,
		},
		{
			name:              "should be valid for monthly aggregation",
			values:            url.Values{"aggregation": {"monthly"}},
			expectAggregation: "monthly",
			expectValid:       true,
		},
		{
			name:              "should be valid for quarterly aggregation",
			values:            url.Values{"aggregation": {"quarterly"}},
			expectAggregation: "quarterly",
			expectValid:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregation, valid := validAggregationQuery(tt.values)

			if valid != tt.expectValid {
				t.Errorf("expected '%t' got '%t'", tt.expectValid, valid)
			}
			if tt.expectValid && aggregation != tt.expectAggregation {
				t.Errorf("expected '%v' got '%v'", tt.expectAggregation, aggregation)
			}
		})
	}
}

func Test_validProjectQuery(t *testing.T) {
	tests := []struct {
		name          string
		values        url.Values
		expectProject string
		expectValid   bool
	}{
		{
			name:          "should be valid when no project is provided",
			values:        url.Values{},
			expectProject: "",
			expectValid:   true,
		},
		{
			name:          "should be valid when no project is provided",
			values:        url.Values{"project": {}},
			expectProject: "",
			expectValid:   true,
		},
		{
			name:          "should not be valid when multiple projects are provided",
			values:        url.Values{"project": {"project_a", "project_b"}},
			expectProject: "",
			expectValid:   false,
		},
		{
			name:          "should not be valid when an empty project is provided",
			values:        url.Values{"project": {""}},
			expectProject: "",
			expectValid:   false,
		},
		{
			name:          "should be valid for a single project",
			values:        url.Values{"project": {"project_a"}},
			expectProject: "project_a",
			expectValid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, valid := validProjectQuery(tt.values)

			if valid != tt.expectValid {
				t.Errorf("expected '%t' got '%t'", tt.expectValid, valid)
			}
			if tt.expectValid && project != tt.expectProject {
				t.Errorf("expected '%v' got '%v'", tt.expectProject, project)
			}
		})
	}
}

func Test_validTimeQueries(t *testing.T) {
	tests := []struct {
		name              string
		values            url.Values
		key               string
		expectTime        time.Time
		expectUsesDefault bool
		expectValid       bool
	}{
		{
			name:              "should use default when no to value is provided",
			values:            url.Values{},
			key:               "to",
			expectTime:        time.Time{},
			expectUsesDefault: true,
			expectValid:       true,
		},
		{
			name:              "should use default when no to value is provided",
			values:            url.Values{"to": {}},
			key:               "to",
			expectTime:        time.Time{},
			expectUsesDefault: true,
			expectValid:       true,
		},
		{
			name:              "should return an error when multiple to values are provided",
			values:            url.Values{"to": {"val1", "val2"}},
			key:               "to",
			expectTime:        time.Time{},
			expectUsesDefault: false,
			expectValid:       false,
		},
		{
			name:              "should return an error when an empty to is provided",
			values:            url.Values{"to": {""}},
			key:               "to",
			expectTime:        time.Time{},
			expectUsesDefault: false,
			expectValid:       false,
		},
		{
			name:              "should return an error when to is in the wrong format",
			values:            url.Values{"to": {"not-a-date"}},
			key:               "to",
			expectTime:        time.Time{},
			expectUsesDefault: false,
			expectValid:       false,
		},
		{
			name:              "should return an error when to is in the wrong format",
			values:            url.Values{"to": {"Mon, 02 Jan 2006 15:04:05 MST"}},
			key:               "to",
			expectTime:        time.Time{},
			expectUsesDefault: false,
			expectValid:       false,
		},
		{
			name:              "should return the parsed time when to is a single formatted time value",
			values:            url.Values{"to": {"2023-01-01T00:00:00Z"}},
			key:               "to",
			expectTime:        time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectUsesDefault: false,
			expectValid:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeValue, usesDefault, valid := validTimeQueries(tt.values, tt.key)

			if usesDefault != tt.expectUsesDefault {
				t.Errorf("expected usesDefault '%t' got '%t'", tt.expectUsesDefault, usesDefault)
			}
			if valid != tt.expectValid {
				t.Errorf("expected valid '%t' got '%t'", tt.expectValid, valid)
			}
			if timeValue != tt.expectTime {
				t.Errorf("expected '%v' got '%v'", tt.expectTime, timeValue)
			}
		})
	}
}

func unixTimesAreAlmostEqual(timeA int64, timeB int64) bool {
	return math.Abs(float64(timeA)-float64(timeB)) < float64(time.Hour.Seconds())
}

func Test_validToAndFromQueries(t *testing.T) {
	now := time.Now()
	nowString := now.Format(time.RFC3339)
	sixMonthsAgo := now.AddDate(0, -6, 0)
	sixMonthsAgoString := sixMonthsAgo.Format(time.RFC3339)
	tests := []struct {
		name        string
		values      url.Values
		expectTo    time.Time
		expectFrom  time.Time
		expectError string
	}{
		{
			name:        "should return error when to is provided but not from",
			values:      url.Values{"to": {sixMonthsAgoString}},
			expectTo:    sixMonthsAgo,
			expectFrom:  time.Time{},
			expectError: "both to and from should be provided or both should be omitted",
		},
		{
			name:        "should return error when from is provided but not to",
			values:      url.Values{"from": {sixMonthsAgoString}},
			expectTo:    time.Time{},
			expectFrom:  sixMonthsAgo,
			expectError: "both to and from should be provided or both should be omitted",
		},
		{
			name:        "should return default values when both to and from are not provided",
			values:      url.Values{},
			expectTo:    now,
			expectFrom:  sixMonthsAgo,
			expectError: "",
		},
		{
			name:        "should return error if to is invalid",
			values:      url.Values{"to": {"not-a-date"}, "from": {sixMonthsAgoString}},
			expectTo:    time.Time{},
			expectFrom:  time.Time{},
			expectError: "to should be provided as a RFC3339 formatted date string or omitted",
		},
		{
			name:        "should return error if from is invalid",
			values:      url.Values{"to": {sixMonthsAgoString}, "from": {"not-a-date"}},
			expectTo:    sixMonthsAgo,
			expectFrom:  time.Time{},
			expectError: "from should be provided as a RFC3339 formatted date string or omitted",
		},
		{
			name:        "should return error if to is in the future",
			values:      url.Values{"to": {now.AddDate(1, 0, 0).Format(time.RFC3339)}, "from": {sixMonthsAgoString}},
			expectTo:    now.AddDate(1, 0, 0),
			expectFrom:  sixMonthsAgo,
			expectError: "to should not be a date in the future",
		},
		{
			name:        "should return error if from is after to",
			values:      url.Values{"to": {sixMonthsAgoString}, "from": {nowString}},
			expectTo:    sixMonthsAgo,
			expectFrom:  now,
			expectError: "from should be a date before to",
		},
		{
			name:        "should return parsed to and from if they are both provided and valid",
			values:      url.Values{"to": {now.AddDate(0, -1, 0).Format(time.RFC3339)}, "from": {sixMonthsAgoString}},
			expectTo:    now.AddDate(0, -1, 0),
			expectFrom:  sixMonthsAgo,
			expectError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toValue, fromValue, err := validToFromQueries(tt.values)
			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !unixTimesAreAlmostEqual(toValue.Unix(), tt.expectTo.Unix()) {
				t.Errorf("expected '%v' got '%v'", tt.expectTo, toValue)
			}
			if !unixTimesAreAlmostEqual(fromValue.Unix(), tt.expectFrom.Unix()) {
				t.Errorf("expected '%v' got '%v'", tt.expectFrom, fromValue)
			}
		})
	}
}

func partialParameterMatch(parametersA service.ServiceParameters, parametersB service.ServiceParameters) bool {
	return parametersA.TypeQuery == parametersB.TypeQuery && parametersA.Aggregation == parametersB.Aggregation && parametersA.Project == parametersB.Project &&
		unixTimesAreAlmostEqual(parametersA.To, parametersB.To) && unixTimesAreAlmostEqual(parametersA.From, parametersB.From)
}

func Test_validServiceParameters(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                    string
		values                  url.Values
		expectServiceParameters service.ServiceParameters
		expectError             string
	}{
		{
			name:                    "should return an error for an invalid type parameter",
			values:                  url.Values{"type": {"not_a_type"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "type should be provided as one of the following: df_average, df_count",
		},
		{
			name:                    "should return an error for an invalid project parameter",
			values:                  url.Values{"type": {"df_count"}, "project": {""}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "project should be provided as a non-empty string or omitted",
		},
		{
			name:                    "should return an error for an invalid to parameter",
			values:                  url.Values{"type": {"df_count"}, "to": {"not-a-date"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "to should be provided as a RFC3339 formatted date string or omitted",
		},
		{
			name:                    "should return an error for an invalid from parameter",
			values:                  url.Values{"type": {"df_count"}, "from": {"not-a-date"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "from should be provided as a RFC3339 formatted date string or omitted",
		},
		{
			name:                    "should return service parameters with defaults for project, to and from",
			values:                  url.Values{"type": {"df_count"}},
			expectServiceParameters: service.ServiceParameters{TypeQuery: "df_count", Project: "", To: now.Unix(), From: now.AddDate(0, -6, 0).Unix()},
			expectError:             "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parameters, err := validServiceParameters(tt.values, []string{"df_average", "df_count"})

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !partialParameterMatch(parameters, tt.expectServiceParameters) {
				t.Errorf("expected '%v' got '%v'", tt.expectServiceParameters, parameters)
			}
		})
	}
}

func Test_ValidMetricServiceParameters(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                    string
		values                  url.Values
		expectServiceParameters service.ServiceParameters
		expectError             string
	}{
		{
			name:                    "should return an error for an invalid type parameter",
			values:                  url.Values{"type": {"df"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "type should be provided as one of the following: df_count, df_average, mltc, cfr, mttr",
		},
		{
			name:                    "should return an error for an invalid aggregation parameter",
			values:                  url.Values{"type": {"df_count"}, "aggregation": {"not_an_aggregation"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "aggregation should be provided as either weekly, monthly or quarterly",
		},
		{
			name:                    "should return service parameters with defaults for aggregation, project, to and from",
			values:                  url.Values{"type": {"df_count"}},
			expectServiceParameters: service.ServiceParameters{TypeQuery: "df_count", Aggregation: "weekly", Project: "", To: now.Unix(), From: now.AddDate(0, -6, 0).Unix()},
			expectError:             "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parameters, err := ValidMetricServiceParameters(tt.values)

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !partialParameterMatch(parameters, tt.expectServiceParameters) {
				t.Errorf("expected '%v' got '%v'", tt.expectServiceParameters, parameters)
			}
		})
	}
}

func Test_ValidBenchmarkServiceParameters(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                    string
		values                  url.Values
		expectServiceParameters service.ServiceParameters
		expectError             string
	}{
		{
			name:                    "should return an error for an invalid type parameter",
			values:                  url.Values{"type": {"df_average"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "type should be provided as one of the following: df, mltc",
		},
		{
			name:                    "should return service parameters with defaults for project, to and from",
			values:                  url.Values{"type": {"df"}},
			expectServiceParameters: service.ServiceParameters{TypeQuery: "df", Project: "", To: now.Unix(), From: now.AddDate(0, -6, 0).Unix()},
			expectError:             "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parameters, err := ValidBenchmarkServiceParameters(tt.values)

			if err == nil && tt.expectError != "" {
				t.Errorf("expected '%v' got no error", tt.expectError)
			}
			if err != nil && err.Error() != tt.expectError {
				t.Errorf("expected '%v' got '%v'", tt.expectError, err)
			}
			if !partialParameterMatch(parameters, tt.expectServiceParameters) {
				t.Errorf("expected '%v' got '%v'", tt.expectServiceParameters, parameters)
			}
		})
	}
}
