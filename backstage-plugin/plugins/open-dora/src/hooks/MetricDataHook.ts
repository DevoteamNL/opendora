import { useApi } from '@backstage/core-plugin-api';
import { useContext, useEffect, useState } from 'react';
import { MetricData } from '../models/MetricData';
import { doraDataServiceApiRef } from '../services/DoraDataService';
import { MetricContext } from '../services/MetricContext';

export const useMetricData = (type: string) => {
  const doraDataService = useApi(doraDataServiceApiRef);
  const [chartData, setChartData] = useState<MetricData | undefined>();
  const [error, setError] = useState<Error | undefined>();
  const { aggregation, team, project } = useContext(MetricContext);

  useEffect(() => {
    doraDataService
      .retrieveMetricDataPoints({
        type: type,
        team: team,
        aggregation: aggregation,
        project: project,
      })
      .then(response => {
        if (response.dataPoints.length > 0) {
          setChartData(response);
        } else {
          setError(new Error('No data found'));
        }
      }, setError);
  }, [aggregation, team, project, doraDataService, type]);

  return { error: error, chartData: chartData };
};
