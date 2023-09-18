package models

type DataPoint struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type Response struct {
	Aggregation string      `json:"aggregation"`
	DataPoints  []DataPoint `json:"dataPoints"`
}
