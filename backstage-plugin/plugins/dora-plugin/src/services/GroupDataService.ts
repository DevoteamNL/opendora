export const getMockData = async (
  // TODO #13: Use these parameters in fetch and remove underscore 
  _groupQueryParam: string,
  _selectedTimeUnit: string,
) => {

  try {
    const data = await fetch('http://localhost:8080/mock-data', {
      method: 'GET',
    });

    return await data.json();
  } catch (error) {
    return {};
  }
};

export const getAncestry = async (component: string | undefined) => {
  try {
    const data = await fetch(
      `http://localhost:7007/api/catalog/entities/by-name/component/default/${component}`,
      {
        method: 'GET',
      },
    );

    return await data.json();
  } catch (error) {
    return {};
  }
};

const GroupDataService = {
  getMockData,
  getAncestry,
};

export default GroupDataService;
