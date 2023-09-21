export interface DeploymentFrequencyData {
  aggregation: string;
  dataPoints: DataPoint[];
}

export interface DataPoint {
  key: 'string';
  value: number;
}
