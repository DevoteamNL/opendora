package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
)

func dfTotalHandler(w http.ResponseWriter, queries url.Values) {
	aggregations, exists := queries["aggregation"]
	aggregation := "weekly"
	if exists && len(aggregations) > 0 {
		aggregation = aggregations[0]
	}
	switch aggregation {
	case "weekly":
		fmt.Fprintf(w, "Hello, %q", html.EscapeString("weekly"))
	case "monthly":
		fmt.Fprintf(w, "Hello, %q", html.EscapeString("monthly"))
	case "quarterly":
		fmt.Fprintf(w, "Hello, %q", html.EscapeString("quarterly"))
	default:
		http.Error(w, "aggregation should be provided as either weekly, monthly or quarterly", http.StatusBadRequest)
	}
}

func main() {

	http.HandleFunc("/dora/api/metric", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		metricTypes, exists := queries["type"]
		if !exists || len(metricTypes) != 1 {
			http.Error(w, "type should be provided as either df_average or df_total", http.StatusBadRequest)
			return
		}
		switch metricTypes[0] {
		case "df_total":
			dfTotalHandler(w, queries)
		case "df_average":
			fmt.Fprintf(w, "Hello, %q", html.EscapeString("average"))
		default:
			http.Error(w, "type should be provided as either df_average or df_total", http.StatusBadRequest)
		}
	})

	log.Fatal(http.ListenAndServe(":10666", nil))
}
