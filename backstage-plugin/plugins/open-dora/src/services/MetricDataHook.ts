import { useApi } from '@backstage/core-plugin-api';
import { useContext, useEffect, useState } from 'react';
import { dfBenchmarkKey } from '../models/DfBenchmarkData';
import { MetricData } from '../models/MetricData';
import { groupDataServiceApiRef } from './GroupDataService';
import { MetricContext } from './MetricContext';

export const useMetricData = (type: string) => {
  const groupDataService = useApi(groupDataServiceApiRef);
  const [chartData, setChartData] = useState<MetricData | undefined>();
  const [error, setError] = useState<Error | undefined>();
  const { aggregation, team, project } = useContext(MetricContext);

  useEffect(() => {
    groupDataService
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
  }, [aggregation, team, project, groupDataService, type]);

  return { error: error, chartData: chartData };
};

export const useMetricOverview = (type: string) => {
  const groupDataService = useApi(groupDataServiceApiRef);
  const [overview, setDfOverview] = useState<dfBenchmarkKey | undefined>();
  const [error, setError] = useState<Error | undefined>();

  useEffect(() => {
    groupDataService.retrieveBenchmarkData({ type: type }).then(response => {
      setDfOverview(response.key);
    }, setError);
  }, [groupDataService, type]);

  return { error: error, overview: overview };
};
