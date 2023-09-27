package sql_client

import "devlake-go/group-sync/api/models"

type MockClient struct {
	WeeklyDataPointsToReturn  []models.DataPoint
	WeeklyErrorToReturn       error
	MonthlyDataPointsToReturn []models.DataPoint
	MonthlyErrorToReturn      error
}

func (client MockClient) QueryTotalDeploymentsWeekly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	if client.WeeklyDataPointsToReturn != nil {
		return client.WeeklyDataPointsToReturn, nil
	}
	if client.WeeklyErrorToReturn != nil {
		return nil, client.WeeklyErrorToReturn
	}
	return []models.DataPoint{}, nil
}
func (client MockClient) QueryTotalDeploymentsMonthly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	if client.MonthlyDataPointsToReturn != nil {
		return client.MonthlyDataPointsToReturn, nil
	}
	if client.MonthlyErrorToReturn != nil {
		return nil, client.MonthlyErrorToReturn
	}
	return []models.DataPoint{}, nil
}
