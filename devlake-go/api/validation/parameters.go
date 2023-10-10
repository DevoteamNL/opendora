package validation

import (
	"fmt"
	"github.com/devoteamnl/opendora/api/service"
	"net/url"
	"time"
)

func validTypeQuery(queries url.Values) (string, bool) {
	metricTypes, exists := queries["type"]
	if !exists || len(metricTypes) != 1 {
		return "", false
	}
	query := metricTypes[0]
	return query, query == "df_count" || query == "df_average"
}

func validAggregationQuery(queries url.Values) (string, bool) {
	aggregations, exists := queries["aggregation"]
	if !exists || len(aggregations) == 0 {
		return "weekly", true
	}
	if len(aggregations) > 1 {
		return "", false
	}
	aggregation := aggregations[0]

	return aggregation, aggregation == "weekly" || aggregation == "monthly" || aggregation == "quarterly"
}

func validProjectQuery(queries url.Values) (string, bool) {
	projects, exists := queries["project"]

	if !exists || len(projects) == 0 {
		return "", true
	}
	if len(projects) > 1 || len(projects[0]) == 0 {
		return "", false
	}

	return projects[0], true
}

func ValidServiceParameters(queries url.Values) (service.ServiceParameters, error) {
	typeQuery, valid := validTypeQuery(queries)
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("type should be provided as either df_average or df_count")
	}
	aggregation, valid := validAggregationQuery(queries)
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("aggregation should be provided as either weekly, monthly or quarterly")

	}
	project, valid := validProjectQuery(queries)
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("project should be provided as a non-empty string or omitted")
	}
	to := time.Now().Unix()
	from := to - (60 * 60 * 24 * 30 * 6)

	return service.ServiceParameters{TypeQuery: typeQuery, Aggregation: aggregation, Project: project, To: to, From: from}, nil
}
