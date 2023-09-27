import { ConfigApi, createApiRef } from '@backstage/core-plugin-api';
import { DeploymentFrequencyData } from '../models/DeploymentFrequencyData';

export const groupDataServiceApiRef = createApiRef<GroupDataService>({
  id: 'plugin.dora-metrics.group-data',
});

export class GroupDataService {
  constructor(private options: { configApi: ConfigApi }) {}

  async retrieveDeploymentFrequencyTotal(
    groupQueryParam: string,
    selectedTimeUnit: string,
  ) {
    const baseUrl = this.options.configApi.getString('dora-metrics.apiBaseUrl');
    const url = new URL(baseUrl);
    url.pathname = 'dora/api/metric';
    url.searchParams.append('type', 'df_count');
    url.searchParams.append('aggregation', selectedTimeUnit);
    url.searchParams.append('team', groupQueryParam);
    const data = await fetch(url.toString(), {
      method: 'GET',
    });
    return (await data.json()) as DeploymentFrequencyData;
  }
}
