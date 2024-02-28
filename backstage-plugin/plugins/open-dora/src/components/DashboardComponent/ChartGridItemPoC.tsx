import { Progress, ResponseErrorPanel } from '@backstage/core-components';
import { Box } from '@material-ui/core';
import React from 'react';
import { useMetricData } from '../../hooks/MetricDataHook';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { Example } from '../BarChartComponent/chartPoC';

export const ChartGridItemPoC = ({
  type,
  label,
}: {
  type: string;
  label: string;
}) => {
  const { chartData, error } = useMetricData(type);

  const chartOrProgressComponent = chartData ? (
    <Example metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const errorOrResponse = error ? (
    <ResponseErrorPanel error={error} />
  ) : (
    chartOrProgressComponent
  );

  return (
    <Box sx={{ bgcolor: '#424242', flex: 1 }}>
      <h1>{label}</h1>
      {errorOrResponse}
    </Box>
  );
};
