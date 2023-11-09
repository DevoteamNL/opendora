package service

import (
	"github.com/devoteamnl/opendora/api/models"
)

type ServiceParameters struct {
	TypeQuery   string
	Aggregation string
	Project     string
	To          int64
	From        int64
}

type MetricService interface {
	ServeRequest(params ServiceParameters) (models.MetricResponse, error)
}

type BenchmarkService interface {
	ServeRequest(params ServiceParameters) (models.BenchmarkResponse, error)
}
