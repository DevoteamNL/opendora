import { ResponseErrorPanel } from '@backstage/core-components';
import { Box, CircularProgress } from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import { useMetricBenchmark } from '../../hooks/MetricBenchmarkHook';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';

export const BenchmarkGridItem = ({ type }: { type: string }) => {
  const [t] = useTranslation();
  const { benchmarkKey, benchmarkValue, error } = useMetricBenchmark(type);

  const testOrProgressComponent = benchmarkKey ? (
    <HighlightTextBoxComponent
      title=""
      text=""
      highlight={`${benchmarkValue} 
        ${t(
          `software_delivery_performance_metrics.overall_labels.${benchmarkKey}`,
        )}`}
      healthStatus={
        {
          'week-elite': 'elite',
          'week-high': 'high',
          'month-medium': 'medium',
          'month-low': 'low',
          elite: 'elite',
          high: 'high',
          medium: 'medium',
          low: 'low',
        }[benchmarkKey]
      }
    />
  ) : (
    <CircularProgress />
  );

  const errorOrResponse = error ? (
    <ResponseErrorPanel error={error} />
  ) : (
    testOrProgressComponent
  );

  return (
    <Box sx={{ bgcolor: '#424242', flex: 1 }}>
      <h3>
        {t(`software_delivery_performance_metrics.labels.benchmark_${type}`)}
      </h3>
      {errorOrResponse}
    </Box>
  );
};
