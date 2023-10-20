package validation

import (
	"fmt"
	"net/url"
	"time"

	"github.com/devoteamnl/opendora/api/service"
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

func validTimeQuery(queries url.Values, key string, defaultValue time.Time) (time.Time, bool) {
	times, exists := queries[key]

	if !exists || len(times) == 0 {
		return defaultValue, true
	}
	if len(times) > 1 || len(times[0]) == 0 {
		return defaultValue, false
	}

	timeValue, err := time.Parse(time.RFC3339, times[0])

	return timeValue, err == nil
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
	to, valid := validTimeQuery(queries, "to", time.Now())
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("to should be provided as a RFC3339 formatted date string or omitted")
	}
	from, valid := validTimeQuery(queries, "from", time.Now().Add(-time.Hour*24*30*6))
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("from should be provided as a RFC3339 formatted date string or omitted")
	}

	return service.ServiceParameters{TypeQuery: typeQuery, Aggregation: aggregation, Project: project, To: to.Unix(), From: from.Unix()}, nil
}
