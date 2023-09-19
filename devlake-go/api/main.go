package main

import (
	"devlake-go/group-sync/api/models"
	"devlake-go/group-sync/api/sql_client"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func dfTotalHandler(client sql_client.ClientInterface, w http.ResponseWriter, queries url.Values) {
	projects, exists := queries["project"]
	if !exists || len(projects) < 1 || len(projects[0]) < 1 {
		http.Error(w, "project should be provided as a non-empty string", http.StatusBadRequest)
		return
	}
	project := projects[0]

	aggregations, exists := queries["aggregation"]
	aggregation := "weekly"
	if exists && len(aggregations) > 0 {
		aggregation = aggregations[0]
	}

	// TODO Make these query parameters
	to := time.Now().Unix()
	from := to - (60 * 60 * 24 * 30 * 6)

	var dataPoints []models.DataPoint
	var err error

	switch aggregation {
	case "weekly":
		dataPoints, err = client.QueryTotalDeploymentsWeekly(project, from, to)
	case "monthly":
		dataPoints, err = client.QueryTotalDeploymentsMonthly(project, from, to)
	case "quarterly":
		// TODO implement quarterly aggregation sql
		dataPoints = []models.DataPoint{}
	default:
		http.Error(w, "aggregation should be provided as either weekly, monthly or quarterly", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.Response{Aggregation: aggregation, DataPoints: dataPoints})
}

func metricHandler(client sql_client.ClientInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		metricTypes, exists := queries["type"]
		if !exists || len(metricTypes) != 1 {
			http.Error(w, "type should be provided as either df_average or df_total", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		switch metricTypes[0] {
		case "df_total":
			dfTotalHandler(client, w, queries)
		case "df_average":
			fmt.Fprint(w, "todo average")
		default:
			http.Error(w, "type should be provided as either df_average or df_total", http.StatusBadRequest)
		}
	}
}

func main() {
	http.HandleFunc("/dora/api/metric", metricHandler(sql_client.New()))

	log.Fatal(http.ListenAndServe(":10666", nil))
}
