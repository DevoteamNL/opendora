import { useApi } from '@backstage/core-plugin-api';
import { useEffect, useState } from 'react';
import { doraDataServiceApiRef } from '../services/DoraDataService';

export const useMetricBenchmark = (type: string) => {
  const doraDataService = useApi(doraDataServiceApiRef);
  const [benchmarkKey, setBenchmarkKey] = useState<string | undefined>();
  const [benchmarkValue, setBenchmarkValue] = useState<string | undefined>();
  const [error, setError] = useState<Error | undefined>();

  useEffect(() => {
    doraDataService.retrieveBenchmarkData({ type: type }).then(response => {
      setBenchmarkKey(response.key);
      setBenchmarkValue(response.value);
    }, setError);
  }, [doraDataService, type]);

  return { error: error, benchmarkKey, benchmarkValue };
};
