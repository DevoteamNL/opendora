import { BarChart } from '@mui/x-charts/BarChart';
import * as React from 'react';
import { DataPoint, MetricData } from '../../models/MetricData';

interface DeploymentFrequencyDataProp {
  metricData: MetricData;
}

export const BarChartComponent = ({
  metricData,
}: DeploymentFrequencyDataProp) => {
  const keys = metricData.dataPoints.map((item: DataPoint) => item.key);
  const values = metricData.dataPoints.map((item: DataPoint) => item.value);

  return (
    <BarChart
      xAxis={[
        {
          data: keys,
          scaleType: 'band',
        },
      ]}
      series={[
        {
          data: values,
        },
      ]}
      height={300}
      margin={{ top: 20, right: 35, bottom: 35, left: 35 }}
    />
  );
};
