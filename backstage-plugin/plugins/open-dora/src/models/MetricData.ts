export interface MetricData {
  aggregation: string;
  dataPoints: DataPoint[];
}

export interface DataPoint {
  key: string;
  value: number;
}
