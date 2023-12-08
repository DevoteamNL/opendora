package service

import (
	"github.com/devoteamnl/opendora/api/models"
	"github.com/devoteamnl/opendora/api/sql_client"
	"github.com/devoteamnl/opendora/api/sql_client/sql_queries"
)

type MltcService struct {
	Client sql_client.ClientInterface
}

func (service MltcService) ServeRequest(params ServiceParameters) (models.MetricResponse, error) {
	typeQueryMap := map[string]string{
		"mltc": sql_queries.MltcSql,
	}

	query := typeQueryMap[params.TypeQuery]

	dataPoints, err := service.Client.QueryDeployments(query, sql_client.QueryParams{To: params.To, From: params.From, Project: params.Project})

	return models.MetricResponse{Aggregation: params.Aggregation, DataPoints: dataPoints}, err
}
