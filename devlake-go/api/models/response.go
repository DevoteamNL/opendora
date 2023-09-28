package models

type DataPoint struct {
	Key   string `json:"key" db:"data_key"`
	Value int    `json:"value" db:"data_value"`
}

type Response struct {
	Aggregation string      `json:"aggregation"`
	DataPoints  []DataPoint `json:"dataPoints"`
}
