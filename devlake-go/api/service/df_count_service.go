package service

import (
	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

type DfCountService struct {
	Client sql_client.ClientInterface
}

func (dfCountService DfCountService) ServeRequest(params ServiceParameters) (models.Response, error) {
	aggregationQueryMap := map[string]string{
		"weekly":    sql_queries.WeeklyDeploymentSql,
		"monthly":   sql_queries.MonthlyDeploymentSql,
		"quarterly": sql_queries.QuarterlyDeploymentSql,
	}

	dataPoints, err := dfCountService.Client.QueryDeployments(aggregationQueryMap[params.Aggregation], sql_client.QueryParams{To: params.To, From: params.From, Project: params.Project})

	return models.Response{Aggregation: params.Aggregation, DataPoints: dataPoints}, err
}
