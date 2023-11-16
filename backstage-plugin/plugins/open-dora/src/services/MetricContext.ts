import { createContext } from 'react';

export const MetricContext = createContext<{
  aggregation: string;
  team?: string;
  project?: string;
}>({
  aggregation: 'weekly',
  team: undefined,
  project: undefined,
});
