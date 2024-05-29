/* eslint-disable no-console */
import { Progress, ResponseErrorPanel } from '@backstage/core-components';
import { Box, useTheme } from '@material-ui/core';
import React from 'react';
import { useMetricData } from '../../hooks/MetricDataHook';
import { LineChartComponent } from '../LineChartComponent/LineChartComponent';

export const ChartGridItem = ({
  type,
  label,
}: {
  type: string;
  label: string;
}) => {
  const { chartData, error } = useMetricData(type);
  const theme = useTheme();
  const chartOrProgressComponent = chartData ? (
    <LineChartComponent metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const errorOrResponse = error ? (
    <ResponseErrorPanel error={error} />
  ) : (
    chartOrProgressComponent
  );
  return (
    <Box
      sx={{
        flex: 1,
        bgcolor: theme.palette.background.paper,
        boxShadow: `
        0px 2px 2px -1px rgba(0,0,0,0.05), 
        0px 2px 2px 0px rgba(0,0,0,0.07),
        0px 1px 5px 0px rgba(0,0,0,0.06)`,
        borderRadius: 10,
      }}
    >
      <h3
        style={{
          color: theme.palette.primary.main,
          fontStyle: 'normal',
          fontSize: '1.3em',
          marginBottom: -15,
          marginLeft: 30,
        }}
      >
        {label}
      </h3>
      {errorOrResponse}
    </Box>
  );
};
