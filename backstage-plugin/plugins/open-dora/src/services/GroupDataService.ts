import { ConfigApi, createApiRef } from '@backstage/core-plugin-api';
import { MetricData } from '../models/MetricData';
import { dfBenchmarkData } from '../models/DfBenchmarkData';

export const groupDataServiceApiRef = createApiRef<GroupDataService>({
  id: 'plugin.open-dora.group-data',
});

export class GroupDataService {
  constructor(private options: { configApi: ConfigApi }) {}

  async retrieveMetricDataPoints(params: {
    type: string;
    team?: string;
    project?: string;
    aggregation?: string;
  }) {
    const baseUrl = this.options.configApi.getString('open-dora.apiBaseUrl');

    const url = new URL(baseUrl);
    url.pathname = 'dora/api/metric';
    for (const [key, value] of Object.entries(params)) {
      if (value) {
        url.searchParams.append(key, value);
      }
    }
    const data = await fetch(url.toString(), {
      method: 'GET',
    });

    if (!data.ok) {
      throw new Error(data.statusText);
    }

    const response = await data.json();
    if (
      response.aggregation === undefined ||
      response.dataPoints === undefined
    ) {
      throw new Error('Unexpected response');
    }

    return response as MetricData;
  }

  async retrieveBenchmarkData(params: { type: string }) {
    const baseUrl = this.options.configApi.getString('open-dora.apiBaseUrl');

    const url = new URL(baseUrl);
    url.pathname = 'dora/api/benchmark';
    for (const [key, value] of Object.entries(params)) {
      if (value) {
        url.searchParams.append(key, value);
      }
    }
    const data = await fetch(url.toString(), {
      method: 'GET',
    });

    if (!data.ok) {
      throw new Error(data.statusText);
    }

    const response = await data.json();
    if (response.key === undefined) {
      throw new Error('Unexpected response');
    }

    return response as dfBenchmarkData;
  }
}
