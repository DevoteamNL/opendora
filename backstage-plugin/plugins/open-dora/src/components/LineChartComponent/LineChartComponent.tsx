import { LineChart } from '@mui/x-charts/LineChart';
import * as React from 'react';
import { DataPoint, MetricData } from '../../models/MetricData';
import { useTheme } from '@material-ui/core';

interface DeploymentFrequencyDataProp {
  metricData: MetricData;
}

export const LineChartComponent = ({
  metricData,
}: DeploymentFrequencyDataProp) => {
  const theme = useTheme();
  const keys = metricData.dataPoints.map((item: DataPoint) => item.key);
  const values = metricData.dataPoints.map((item: DataPoint) => item.value);
  return (
    <LineChart
      xAxis={[
        {
          data: keys,
          scaleType: 'point',
          valueFormatter: v => v.replace(/ /, '\n'),
        },
      ]}
      series={[
        {
          data: values,
          area: true,
          color: theme.palette.primary.main,
        },
      ]}
      height={300}
    />
  );
};
