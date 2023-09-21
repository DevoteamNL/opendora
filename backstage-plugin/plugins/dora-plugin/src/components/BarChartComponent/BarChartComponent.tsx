import * as React from 'react';
import { BarChart } from '@mui/x-charts/BarChart';
import {
  DataPoint,
  DeploymentFrequencyData,
} from '../../models/DeploymentFrequencyData';

interface DeploymentFrequencyDataProp {
  deploymentFrequencyData: DeploymentFrequencyData;
}

export default function SimpleCharts({
  deploymentFrequencyData,
}: DeploymentFrequencyDataProp) {
  const keys = deploymentFrequencyData.dataPoints.map(
    (item: DataPoint) => item.key,
  );
  const values = deploymentFrequencyData.dataPoints.map(
    (item: DataPoint) => item.value,
  );

  return (
    <BarChart
      xAxis={[
        {
          id: 'barCategories',
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
    />
  );
}
