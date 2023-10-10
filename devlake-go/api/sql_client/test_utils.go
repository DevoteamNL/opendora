package sql_client

import "github.com/devoteamnl/opendora/api/models"

type MockDataReturn struct {
	Data []models.DataPoint
	Err  error
}
type MockClient struct {
	MockDataMap map[string]MockDataReturn
}

func (client MockClient) QueryDeployments(query string, params QueryParams) ([]models.DataPoint, error) {
	if client.MockDataMap[query].Data != nil {
		return client.MockDataMap[query].Data, nil
	}
	if client.MockDataMap[query].Err != nil {
		return nil, client.MockDataMap[query].Err
	}
	return []models.DataPoint{}, nil
}
