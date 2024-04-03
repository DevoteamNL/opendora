import { ConfigApi, createApiRef } from '@backstage/core-plugin-api';
import { dfBenchmarkData } from '../models/BenchmarkData';
import { MetricData } from '../models/MetricData';

export const doraDataServiceApiRef = createApiRef<DoraDataService>({
  id: 'plugin.open-dora.data-service',
});

export class DoraDataService {
  constructor(private options: { configApi: ConfigApi }) {}

  private async retrieveData(params: Record<string, string>, path: string) {
    const baseUrl = this.options.configApi.getString('open-dora.apiBaseUrl');

    const url = new URL(baseUrl);
    url.pathname = path;
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

    return await data.json();
  }

  async retrieveMetricDataPoints(params: {
    type: string;
    team?: string;
    project?: string;
    aggregation?: string;
  }) {
    const response = await this.retrieveData(params, 'dora/api/metric');
    if (
      response.aggregation === undefined ||
      response.dataPoints === undefined
    ) {
      throw new Error('Unexpected response');
    }

    return response as MetricData;
  }

  async retrieveBenchmarkData(params: { type: string }) {
    const response = await this.retrieveData(params, 'dora/api/benchmark');

    if (response.key === undefined) {
      throw new Error('Unexpected response');
    }

    return response as dfBenchmarkData;
  }
}
