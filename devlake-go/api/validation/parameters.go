package validation

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/devoteamnl/opendora/api/service"
)

func validTypeQuery(queries url.Values, validTypes []string) (string, bool) {
	metricTypes, exists := queries["type"]
	if !exists || len(metricTypes) != 1 {
		return "", false
	}
	query := metricTypes[0]
	return query, slices.Contains(validTypes, query)
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

func validTimeQueries(queries url.Values, key string) (value time.Time, usesDefault bool, valid bool) {
	times, exists := queries[key]
	usesDefault = !exists || len(times) == 0
	valid = len(times) <= 1
	if usesDefault || !valid {
		return
	}
	value, err := time.Parse(time.RFC3339, times[0])
	valid = err == nil
	return
}

func validToFromQueries(queries url.Values) (to time.Time, from time.Time, err error) {
	now := time.Now()

	to, toShouldUseDefault, valid := validTimeQueries(queries, "to")
	if !valid {
		err = fmt.Errorf("to should be provided as a RFC3339 formatted date string or omitted")
		return
	}
	from, fromShouldUseDefault, valid := validTimeQueries(queries, "from")
	if !valid {
		err = fmt.Errorf("from should be provided as a RFC3339 formatted date string or omitted")
		return
	}

	if toShouldUseDefault && fromShouldUseDefault {
		return now, now.AddDate(0, -6, 0), nil
	}
	if toShouldUseDefault || fromShouldUseDefault {
		err = fmt.Errorf("both to and from should be provided or both should be omitted")
	} else if to.Compare(now) > 0 {
		err = fmt.Errorf("to should not be a date in the future")
	} else if from.Compare(to) > 0 {
		err = fmt.Errorf("from should be a date before to")
	}

	return
}

func validServiceParameters(queries url.Values, validTypes []string) (service.ServiceParameters, error) {
	typeQuery, valid := validTypeQuery(queries, validTypes)
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("type should be provided as one of the following: %s",
			strings.Join(validTypes, ", "))
	}
	project, valid := validProjectQuery(queries)
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("project should be provided as a non-empty string or omitted")
	}

	to, from, err := validToFromQueries(queries)

	if err != nil {
		return service.ServiceParameters{}, err
	}

	return service.ServiceParameters{TypeQuery: typeQuery, Project: project, To: to.Unix(), From: from.Unix()}, nil
}

func ValidBenchmarkServiceParameters(queries url.Values) (service.ServiceParameters, error) {
	return validServiceParameters(queries, []string{"df", "mltc"})
}

func ValidMetricServiceParameters(queries url.Values) (service.ServiceParameters, error) {
	aggregation, valid := validAggregationQuery(queries)
	if !valid {
		return service.ServiceParameters{}, fmt.Errorf("aggregation should be provided as either weekly, monthly or quarterly")
	}

	serviceParameters, err := validServiceParameters(queries, []string{"df_count", "df_average", "mltc"})
	if err != nil {
		return serviceParameters, err
	}
	serviceParameters.Aggregation = aggregation

	return serviceParameters, err
}
