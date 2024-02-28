import { ResponseErrorPanel } from '@backstage/core-components';
import { Box, CircularProgress } from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import { useMetricBenchmark } from '../../hooks/MetricBenchmarkHook';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';

export const BenchmarkGridItem = ({ type }: { type: string }) => {
  const [t] = useTranslation();
  const { benchmark, error } = useMetricBenchmark(type);

  const testOrProgressComponent = benchmark ? (
    <HighlightTextBoxComponent
      title=""
      text=""
      highlight={t(`software_delivery_performance_metrics.overall_labels.${benchmark}`)}
      healthStatus={
        {
          'on-demand': 'positive',
          'lt-6month': 'critical',
          'week-month': 'neutral',
          'month-6month': 'negative',
          'lt-1hour': 'positive',
          'lt-1week': 'neutral',
          'week-6month': 'negative',
          'mt-6month': 'critical',
          'lt-15%': 'positive',
          '15-24.99%': 'neutral',
          '25-30%': 'negative',
          'mt-30%': 'critical',
        }[benchmark]
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
      <h1>{t(`software_delivery_performance_metrics.labels.${type}`)}</h1>
      {errorOrResponse}
    </Box>
  );
};
