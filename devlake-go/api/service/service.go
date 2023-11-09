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

type Service[R models.Response] interface {
	ServeRequest(params ServiceParameters) (R, error)
}
