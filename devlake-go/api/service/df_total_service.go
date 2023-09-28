package service

import (
	"devlake-go/group-sync/api/models"
	"devlake-go/group-sync/api/sql_client"
)

type DfTotalService struct {
	Client sql_client.ClientInterface
}

func (dfTotalService DfTotalService) ServeRequest(params ServiceParameters) (models.Response, error) {
	if params.Aggregation == "quarterly" {
		// TODO implement quarterly aggregation sql
		return models.Response{Aggregation: params.Aggregation, DataPoints: []models.DataPoint{}}, nil
	}
	aggregationQueryMap := map[string]string{
		"weekly":  sql_client.WEEKLY_DEPLOYMENT_SQL,
		"monthly": sql_client.MONTHLY_DEPLOYMENT_SQL,
	}

	dataPoints, err := dfTotalService.Client.QueryDeployments(aggregationQueryMap[params.Aggregation], sql_client.QueryParams{To: params.To, From: params.From, Project: params.Project})

	return models.Response{Aggregation: params.Aggregation, DataPoints: dataPoints}, err
}
