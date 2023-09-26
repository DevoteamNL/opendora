package main

import (
	"devlake-go/group-sync/api/service"
	"devlake-go/group-sync/api/sql_client"
	"devlake-go/group-sync/api/validation"
	"encoding/json"
	"log"
	"net/http"
)

func metricHandler(client sql_client.ClientInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		parameters, err := validation.ValidServiceParameters(queries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		services := map[string]service.Service{
			"df_total": service.DfTotalService{Client: client},
		}

		response, err := services[parameters.TypeQuery].ServeRequest(parameters)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	http.HandleFunc("/dora/api/metric", metricHandler(sql_client.New()))

	log.Fatal(http.ListenAndServe(":10666", nil))
}