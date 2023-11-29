import { useApi } from '@backstage/core-plugin-api';
import { useEffect, useState } from 'react';
import { dfBenchmarkKey } from '../models/DfBenchmarkData';
import { doraDataServiceApiRef } from '../services/DoraDataService';

export const useMetricBenchmark = (type: string) => {
  const doraDataService = useApi(doraDataServiceApiRef);
  const [benchmark, setDfBenchmark] = useState<dfBenchmarkKey | undefined>();
  const [error, setError] = useState<Error | undefined>();

  useEffect(() => {
    doraDataService.retrieveBenchmarkData({ type: type }).then(response => {
      setDfBenchmark(response.key);
    }, setError);
  }, [doraDataService, type]);

  return { error: error, benchmark: benchmark };
};
