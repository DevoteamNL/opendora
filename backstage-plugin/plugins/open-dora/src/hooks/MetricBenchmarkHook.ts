import { useApi } from '@backstage/core-plugin-api';
import { useEffect, useState } from 'react';
import { dfBenchmarkKey } from '../models/DfBenchmarkData';
import { groupDataServiceApiRef } from '../services/GroupDataService';

export const useMetricBenchmark = (type: string) => {
  const groupDataService = useApi(groupDataServiceApiRef);
  const [benchmark, setDfBenchmark] = useState<dfBenchmarkKey | undefined>();
  const [error, setError] = useState<Error | undefined>();

  useEffect(() => {
    groupDataService.retrieveBenchmarkData({ type: type }).then(response => {
      setDfBenchmark(response.key);
    }, setError);
  }, [groupDataService, type]);

  return { error: error, benchmark: benchmark };
};
