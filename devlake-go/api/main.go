package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/devoteamnl/opendora/api/service"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/validation"
)

func handler(validateParameters func(queries url.Values) (service.ServiceParameters, error), serveRequest func(service.ServiceParameters) (any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		parameters, err := validateParameters(queries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := serveRequest(parameters)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		responseError := json.NewEncoder(w).Encode(response)
		if responseError != nil {
			http.Error(w, responseError.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func metricHandler(client sql_client.ClientInterface) func(w http.ResponseWriter, r *http.Request) {
	dfService := service.MetricDfService{Client: client}
	serviceMap := map[string]service.MetricService{
		"df_count":   dfService,
		"df_average": dfService,
	}

	return handler(validation.ValidMetricServiceParameters, func(parameters service.ServiceParameters) (any, error) {
		return serviceMap[parameters.TypeQuery].ServeRequest(parameters)
	})
}

func benchmarkHandler(client sql_client.ClientInterface) func(w http.ResponseWriter, r *http.Request) {
	dfService := service.BenchmarkDfService{Client: client}
	serviceMap := map[string]service.BenchmarkService{
		"df": dfService,
	}

	return handler(validation.ValidBenchmarkServiceParameters, func(parameters service.ServiceParameters) (any, error) {
		return serviceMap[parameters.TypeQuery].ServeRequest(parameters)
	})
}

func main() {
	sqlClient := sql_client.New()
	http.HandleFunc("/dora/api/metric", metricHandler(sqlClient))
	http.HandleFunc("/dora/api/benchmark", benchmarkHandler(sqlClient))
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "ok")
	})
	log.Fatal(http.ListenAndServe(":10666", nil))
}
