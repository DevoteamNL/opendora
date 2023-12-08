import { Progress, ResponseErrorPanel } from '@backstage/core-components';
import { Grid } from '@material-ui/core';
import React from 'react';
import { useMetricData } from '../../hooks/MetricDataHook';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';

export const ChartGridItem = ({
  type,
  label,
}: {
  type: string;
  label: string;
}) => {
  const { chartData, error } = useMetricData(type);

  const chartOrProgressComponent = chartData ? (
    <BarChartComponent metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const errorOrResponse = error ? (
    <ResponseErrorPanel error={error} />
  ) : (
    chartOrProgressComponent
  );

  return (
    <Grid item xs={12} className="gridBorder">
      <div className="gridBoxText">
        <h1>{label}</h1>
        {errorOrResponse}
      </div>
    </Grid>
  );
};
