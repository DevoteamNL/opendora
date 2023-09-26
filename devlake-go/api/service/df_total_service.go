package service

import (
	"devlake-go/group-sync/api/models"
	"devlake-go/group-sync/api/sql_client"
)

type DfTotalService struct {
	Client sql_client.ClientInterface
}

func (dfTotalService DfTotalService) ServeRequest(params ServiceParameters) (models.Response, error) {

	var dataPoints []models.DataPoint
	var err error

	if params.Aggregation == "weekly" {
		dataPoints, err = dfTotalService.Client.QueryTotalDeploymentsWeekly(params.Project, params.From, params.To)
	}
	if params.Aggregation == "monthly" {
		dataPoints, err = dfTotalService.Client.QueryTotalDeploymentsMonthly(params.Project, params.From, params.To)
	}
	if params.Aggregation == "quarterly" {
		// TODO implement quarterly aggregation sql
		dataPoints = []models.DataPoint{}
	}

	return models.Response{Aggregation: params.Aggregation, DataPoints: dataPoints}, err
}
