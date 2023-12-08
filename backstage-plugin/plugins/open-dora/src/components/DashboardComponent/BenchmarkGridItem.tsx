import { ResponseErrorPanel } from '@backstage/core-components';
import { CircularProgress, Grid } from '@material-ui/core';
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
      highlight={t(`deployment_frequency.overall_labels.${benchmark}`)}
      healthStatus={
        {
          'on-demand': 'positive',
          'lt-6month': 'critical',
          'week-month': 'neutral',
          'month-6month': 'negative',
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
    <Grid item xs={12} className="gridBorder">
      <div className="gridBoxText">
        <Grid container>
          <Grid item xs={3}>
            {errorOrResponse}
          </Grid>
        </Grid>
      </div>
    </Grid>
  );
};
