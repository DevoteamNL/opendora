package models

type DataPoint struct {
	Key   string  `json:"key" db:"data_key"`
	Value float32 `json:"value" db:"data_value"`
}

type Response interface {
	MetricResponse | BenchmarkResponse
}

type MetricResponse struct {
	Aggregation string      `json:"aggregation"`
	DataPoints  []DataPoint `json:"dataPoints"`
}

type BenchmarkResponse struct {
	Key string `json:"key" db:"data_key"`
}
