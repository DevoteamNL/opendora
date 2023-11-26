export interface dfBenchmarkData {
  key: dfBenchmarkKey;
}
export type dfBenchmarkKey =
  | 'on-demand'
  | 'week-month'
  | 'month-6month'
  | 'lt-6month';
