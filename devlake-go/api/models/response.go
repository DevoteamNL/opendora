package models

type DataPoint struct {
	Key   string `json:"key" db:"deployment_key"`
	Value int    `json:"value" db:"deployment_value"`
}

type Response struct {
	Aggregation string      `json:"aggregation"`
	DataPoints  []DataPoint `json:"dataPoints"`
}
