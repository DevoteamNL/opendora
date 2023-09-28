import { BarChart } from '@mui/x-charts/BarChart';
import * as React from 'react';
import {
  DataPoint,
  DeploymentFrequencyData,
} from '../../models/DeploymentFrequencyData';

interface DeploymentFrequencyDataProp {
  deploymentFrequencyData: DeploymentFrequencyData;
}

export const BarChartComponent = ({
  deploymentFrequencyData,
}: DeploymentFrequencyDataProp) => {
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
};
