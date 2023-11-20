export interface dfBenchmarkData {
  key: dfEnum;
}

export enum dfEnum {
  OnDemand = 'on-demand',
  WeekMont = 'week-month',
  Month6Month = 'month-6month',
  LT6Month = 'lt-6month',
}
