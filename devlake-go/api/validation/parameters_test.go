package validation

import (
	"devlake-go/group-sync/api/service"
	"net/url"
	"testing"
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
			typeQuery, valid := validTypeQuery(tt.values)

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

func partialParameterMatch(parametersA service.ServiceParameters, parametersB service.ServiceParameters) bool {
	return parametersA.TypeQuery == parametersB.TypeQuery && parametersA.Aggregation == parametersB.Aggregation && parametersA.Project == parametersB.Project
}

func Test_ValidServiceParameters(t *testing.T) {
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
			expectError:             "type should be provided as either df_average or df_count",
		},
		{
			name:                    "should return an error for an invalid aggregation parameter",
			values:                  url.Values{"type": {"df_count"}, "aggregation": {"not_an_aggregation"}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "aggregation should be provided as either weekly, monthly or quarterly",
		},
		{
			name:                    "should return an error for an invalid project parameter",
			values:                  url.Values{"type": {"df_count"}, "project": {""}},
			expectServiceParameters: service.ServiceParameters{},
			expectError:             "project should be provided as a non-empty string or omitted",
		},
		{
			name:                    "should return service parameters with defaults for aggregation and project",
			values:                  url.Values{"type": {"df_count"}},
			expectServiceParameters: service.ServiceParameters{TypeQuery: "df_count", Aggregation: "weekly", Project: "", To: 0, From: 0},
			expectError:             "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parameters, err := ValidServiceParameters(tt.values)

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

func Test_ServiceParametersDefaultDate(t *testing.T) {
	parameters, err := ValidServiceParameters(url.Values{"type": {"df_count"}})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parameters.From != parameters.To-(60*60*24*30*6) {
		t.Errorf("expected the 'from' parameter to be 6 months before the 'to'  parameter. got from '%v', to '%v'", parameters.From, parameters.To)
	}
}
