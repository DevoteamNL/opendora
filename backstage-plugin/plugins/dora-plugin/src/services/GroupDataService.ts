import { DeploymentFrequencyData } from '../models/DeploymentFrequencyData';

export const getMockData = async (
  // TODO #13: Use these parameters in fetch and remove underscore
  _groupQueryParam: string,
  _selectedTimeUnit: string,
) => {
  try {
    const data = await fetch('http://localhost:8080/mock-data', {
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
