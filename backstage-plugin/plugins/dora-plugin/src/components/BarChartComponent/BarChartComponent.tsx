import * as React from 'react';
import { BarChart } from '@mui/x-charts/BarChart';

export default function SimpleCharts({ChartData}) {


  let keys = ["0"];
  let values = [0];
  if (ChartData && ChartData.dataPoints){
   keys = ChartData.dataPoints.map(item => item.key)
   values = ChartData.dataPoints.map(item => item.value)
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
    />
  );
}
