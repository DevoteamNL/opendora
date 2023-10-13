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

type Service interface {
	ServeRequest(params ServiceParameters) (models.Response, error)
}
