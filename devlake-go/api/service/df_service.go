package service

import (
	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

type DfService struct {
	Client sql_client.ClientInterface
}

func (service DfService) ServeRequest(params ServiceParameters) (models.Response, error) {
	aggregationQueryMap := map[string]string{
		"weekly":    sql_queries.WeeklyDeploymentSql,
		"monthly":   sql_queries.MonthlyDeploymentSql,
		"quarterly": sql_queries.QuarterlyDeploymentSql,
	}

	typeQueryMap := map[string]string{
		"df_count":   sql_queries.CountSql,
		"df_average": sql_queries.AverageSql,
	}

	query := aggregationQueryMap[params.Aggregation] + typeQueryMap[params.TypeQuery]

	dataPoints, err := service.Client.QueryDeployments(query, sql_client.QueryParams{To: params.To, From: params.From, Project: params.Project})

	return models.Response{Aggregation: params.Aggregation, DataPoints: dataPoints}, err
}
