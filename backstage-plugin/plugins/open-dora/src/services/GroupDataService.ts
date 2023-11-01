// TODO: Add tests and remove ignore.
/* istanbul ignore file */
import { ConfigApi, createApiRef } from '@backstage/core-plugin-api';
import { MetricData } from '../models/MetricData';

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
    return (await data.json()) as MetricData;
  }
}
