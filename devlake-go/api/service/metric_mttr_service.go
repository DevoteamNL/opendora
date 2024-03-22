package service

import (
	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

type MetricMttrService struct {
	Client sql_client.ClientInterface
}

func (service MetricMttrService) ServeRequest(params ServiceParameters) (models.MetricResponse, error) {
	aggregationQueryMap := map[string]string{
		"weekly":    sql_queries.WeeklyMttrSql,
		"monthly":   sql_queries.MonthlyMttrSql,
		"quarterly": sql_queries.QuarterlyMttrSql,
	}

	query := aggregationQueryMap[params.Aggregation]

	dataPoints, err := service.Client.QueryDeployments(query, sql_client.QueryParams{To: params.To, From: params.From, Project: params.Project})

	return models.MetricResponse{Aggregation: params.Aggregation, DataPoints: dataPoints}, err
}
