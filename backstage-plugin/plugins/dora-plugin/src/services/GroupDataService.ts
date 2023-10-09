import { ConfigApi, createApiRef } from '@backstage/core-plugin-api';
import { MetricData } from '../models/MetricData';

export const groupDataServiceApiRef = createApiRef<GroupDataService>({
  id: 'plugin.open-dora.group-data',
});

export class GroupDataService {
  constructor(private options: { configApi: ConfigApi }) {}

  async retrieveMetricDataPoints(
    metricQueryParam: string,
    groupQueryParam: string,
    selectedTimeUnit: string,
  ) {
    const baseUrl = this.options.configApi.getString('open-dora.apiBaseUrl');
    const url = new URL(baseUrl);
    url.pathname = 'dora/api/metric';
    url.searchParams.append('type', metricQueryParam);
    url.searchParams.append('aggregation', selectedTimeUnit);
    url.searchParams.append('team', groupQueryParam);
    const data = await fetch(url.toString(), {
      method: 'GET',
    });
    return (await data.json()) as MetricData;
  }
}
