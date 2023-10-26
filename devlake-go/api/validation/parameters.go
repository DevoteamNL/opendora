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

	to, from, err := validToFromQueries(queries)

	if err != nil {
		return service.ServiceParameters{}, err
	}

	return service.ServiceParameters{TypeQuery: typeQuery, Aggregation: aggregation, Project: project, To: to.Unix(), From: from.Unix()}, nil
}
