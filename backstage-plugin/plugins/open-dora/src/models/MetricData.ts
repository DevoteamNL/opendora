export interface MetricData {
  aggregation: 'weekly' | 'monthly' | 'quarterly';
  dataPoints: DataPoint[];
}

export interface DataPoint {
  key: string;
  value: number;
}
