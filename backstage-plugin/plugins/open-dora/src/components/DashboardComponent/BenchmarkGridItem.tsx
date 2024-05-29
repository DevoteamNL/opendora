import { ResponseErrorPanel } from '@backstage/core-components';
import { Box, CircularProgress } from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import { useMetricBenchmark } from '../../hooks/MetricBenchmarkHook';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import { useTheme } from '@mui/material/styles';

export const BenchmarkGridItem = ({ type }: { type: string }) => {
  const [t] = useTranslation();
  const { benchmarkKey, benchmarkValue, error } = useMetricBenchmark(type);
  const theme = useTheme();

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
        }[benchmarkKey] || 'neutral'
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
    <div
      style={{
        display: 'flex',

        flexDirection: 'column',
        color: theme.palette.text.primary,
        fontStyle: 'normal',
        marginRight: 30,
        marginBottom: 25,
        fontSize: '1em',
        textAlign: 'start',
      }}
    >
      {t(`software_delivery_performance_metrics.labels.benchmark_${type}`)}
      {errorOrResponse}
    </div>
  );
};
