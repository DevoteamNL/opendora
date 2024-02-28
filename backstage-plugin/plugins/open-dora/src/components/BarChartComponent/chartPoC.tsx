import { LineChart, Line } from 'recharts';
import { DataPoint, MetricData } from '../../models/MetricData';
import React, { PureComponent } from 'react';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';

const data = [
  {
    name: 'Page A',
    uv: 4000,
    pv: 2400,
    amt: 2400,
  },
  {
    name: 'Page B',
    uv: 3000,
    pv: 1398,
    amt: 2210,
  },
  {
    name: 'Page C',
    uv: 2000,
    pv: 9800,
    amt: 2290,
  },
  {
    name: 'Page D',
    uv: 2780,
    pv: 3908,
    amt: 2000,
  },
  {
    name: 'Page E',
    uv: 1890,
    pv: 4800,
    amt: 2181,
  },
  {
    name: 'Page F',
    uv: 2390,
    pv: 3800,
    amt: 2500,
  },
  {
    name: 'Page G',
    uv: 3490,
    pv: 4300,
    amt: 2100,
  },
];

const renderLineChart = (
  <ResponsiveContainer width="100%" height="100%">
    <AreaChart
      width={500}
      height={400}
      data={data}
      margin={{
        top: 10,
        right: 30,
        left: 0,
        bottom: 0,
      }}
    >
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="name" />
      <YAxis />
      <Tooltip />
      <Area type="monotone" dataKey="uv" stroke="#8884d8" fill="#8884d8" />
    </AreaChart>
  </ResponsiveContainer>
);

interface DeploymentFrequencyDataProp {
  metricData: MetricData;
}

export const Example = ({ metricData }: DeploymentFrequencyDataProp) => {
  const myMetricData = metricData;
  const keys = metricData.dataPoints.map((item: DataPoint) => item.key);
  const values = metricData.dataPoints.map((item: DataPoint) => item.value);
  return (
    <div style={{ width: 'auto', height: ' 300px' }}>
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart
          width={500}
          height={400}
          data={myMetricData.dataPoints}
          margin={{
            top: 10,
            right: 30,
            left: 0,
            bottom: 0,
          }}
        >
          <defs>
            <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#ff964f" stopOpacity={1} />
              <stop offset="95%" stopColor="#ff964f" stopOpacity={0.2} />
            </linearGradient>
          </defs>
          {/* <CartesianGrid strokeDasharray="3 3" /> */}
          <XAxis dataKey="key" />
          <YAxis />
          <Tooltip />
          <Area
            type="monotone"
            dataKey="value"
            stroke="#ff964f"
            // fill="rgba(255, 250, 160, 0.7)"
            fillOpacity={1}
            fill="url(#colorUv)"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};
