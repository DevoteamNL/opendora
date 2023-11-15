package sql_client

import "github.com/devoteamnl/opendora/api/models"

type MockDeploymentsDataReturn struct {
	Data []models.DataPoint
	Err  error
}
type MockBenchmarkDataReturn struct {
	Data string
	Err  error
}
type MockClient struct {
	MockDeploymentsDataMap map[string]MockDeploymentsDataReturn
	MockBenchmarkDataMap   map[string]MockBenchmarkDataReturn
}

func (client MockClient) QueryDeployments(query string, params QueryParams) ([]models.DataPoint, error) {
	if client.MockDeploymentsDataMap[query].Data != nil {
		return client.MockDeploymentsDataMap[query].Data, nil
	}
	if client.MockDeploymentsDataMap[query].Err != nil {
		return nil, client.MockDeploymentsDataMap[query].Err
	}
	return []models.DataPoint{}, nil
}

func (client MockClient) QueryBenchmark(query string, params QueryParams) (string, error) {
	if client.MockBenchmarkDataMap[query].Data != "" {
		return client.MockBenchmarkDataMap[query].Data, nil
	}
	if client.MockBenchmarkDataMap[query].Err != nil {
		return "", client.MockBenchmarkDataMap[query].Err
	}
	return "", nil
}
