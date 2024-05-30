import { ResponseErrorPanel } from '@backstage/core-components';
import { Box, CircularProgress, Grid } from '@material-ui/core';
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
    <Grid
      item
      style={{
        width: '100%',
        backgroundColor: theme.palette.background.paper,
        boxShadow: `
        0px 2px 2px -1px rgba(0,0,0,0.05), 
        0px 2px 2px 0px rgba(0,0,0,0.07),
        0px 1px 5px 0px rgba(0,0,0,0.06)`,
        borderRadius: 10,
        textAlign: 'center',
      }}
    >
      <Box sx={{}}>
        <div
          style={{
            display: 'flex',
            padding: 20,
            flexDirection: 'column',
            color: theme.palette.text.primary,
            fontStyle: 'normal',
            fontSize: '1em',
            textAlign: 'center',
          }}
        >
          {t(`software_delivery_performance_metrics.labels.benchmark_${type}`)}
          {errorOrResponse}
        </div>
      </Box>
    </Grid>
  );
};
