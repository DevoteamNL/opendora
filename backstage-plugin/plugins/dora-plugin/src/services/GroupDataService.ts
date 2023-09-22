import { DeploymentFrequencyData } from '../models/DeploymentFrequencyData';

export const getMockData = async (
  groupQueryParam: string,
  selectedTimeUnit: string,
) => {
  try {
    const url = new URL('http://localhost:10666/dora/api/metric');
    url.searchParams.append('type', 'df_count');
    url.searchParams.append('aggregation', selectedTimeUnit);
    url.searchParams.append('team', groupQueryParam);
    const data = await fetch(url, {
      method: 'GET',
    });
    return (await data.json()) as DeploymentFrequencyData;
  } catch (e) {
    throw e;
  }
};

const GroupDataService = {
  getMockData,
};

export default GroupDataService;
