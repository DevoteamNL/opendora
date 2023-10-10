package service

import (
	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
)

type DfAverageService struct {
	Client sql_client.ClientInterface
}

func (dfAverageService DfAverageService) ServeRequest(params ServiceParameters) (models.Response, error) {
	return models.Response{Aggregation: params.Aggregation, DataPoints: []models.DataPoint{}}, nil
}
