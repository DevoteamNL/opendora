import { DeploymentFrequencyData } from '../models/DeploymentFrequencyData';

export const getMockData = async (
  groupQueryParam: string,
  selectedTimeUnit: string,
) => {
  try {
    const data = await fetch(`http://localhost:10666/dora/api/metric?type=df_count&aggregation=${selectedTimeUnit}&team=${groupQueryParam}`, {
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
