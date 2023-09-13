import React from 'react';

export const getMockData = async () => {
  const data = await fetch('http://localhost:8080/mock-data', {
    method: 'GET',
  });

  return await data.json();
};

const GroupDataService = {
  getMockData,
};

export default GroupDataService;
