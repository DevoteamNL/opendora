import { BarChart } from '@mui/x-charts/BarChart';
import * as React from 'react';
import { DataPoint, MetricData } from '../../models/MetricData';

interface DeploymentFrequencyDataProp {
  metricData: MetricData;
}

export const BarChartComponent = ({
  metricData,
}: DeploymentFrequencyDataProp) => {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const dataPoints = metricData.dataPoints.slice(
    currentIndex,
    currentIndex + 5,
  );
  const keys = dataPoints.map((item: DataPoint) => item.key);
  const values = dataPoints.map((item: DataPoint) => item.value);

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
      margin={{ top: 20, right: 35, bottom: 35, left: 35 }}
    />
  );
};
