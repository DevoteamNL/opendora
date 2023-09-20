import * as React from 'react';
import { BarChart } from '@mui/x-charts/BarChart';
import { DataPoint, DeploymentFrequencyDataProp } from '../../models/DeploymentFrequencyData';

export default function SimpleCharts({deploymentFrequencyData}: DeploymentFrequencyDataProp) {

  let keys = ['0'];
  let values = [0];
  if (deploymentFrequencyData && deploymentFrequencyData.dataPoints) {
    keys = deploymentFrequencyData.dataPoints.map((item:DataPoint) => item.key);
    values = deploymentFrequencyData.dataPoints.map((item:DataPoint)  => item.value);
  }
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
