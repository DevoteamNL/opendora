import {
  Content,
  Header,
  Page,
  Progress,
  ResponseErrorPanel,
  SupportButton,
} from '@backstage/core-components';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { CircularProgress, Grid } from '@material-ui/core';
import React from 'react';
import { useTranslation } from 'react-i18next';
import { useMetricBenchmark } from '../../hooks/MetricBenchmarkHook';
import { useMetricData } from '../../hooks/MetricDataHook';
import '../../i18n';
import { MetricContext } from '../../services/MetricContext';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import './DashboardComponent.css';

export interface DashboardComponentProps {
  entityName?: string;
  entityGroup?: string;
}

const ChartGridItem = ({ type, label }: { type: string; label: string }) => {
  const { chartData, error } = useMetricData(type);

  const chartOrProgressComponent = chartData ? (
    <BarChartComponent metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const errorOrResponse = error ? (
    <ResponseErrorPanel error={error} />
  ) : (
    chartOrProgressComponent
  );

  return (
    <Grid item xs={12} className="gridBorder">
      <div className="gridBoxText">
        <h1>{label}</h1>
        {errorOrResponse}
      </div>
    </Grid>
  );
};

const BenchmarkGridItem = ({ type }: { type: string }) => {
  const [t] = useTranslation();
  const { benchmark, error } = useMetricBenchmark(type);

  const testOrProgressComponent = benchmark ? (
    <HighlightTextBoxComponent
      title=""
      text=""
      highlight={t(`deployment_frequency.overall_labels.${benchmark}`)}
      // to do: think of text colouring for different scenarios
      textColour="positiveHighlight"
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

export const DashboardComponent = ({
  entityName,
  entityGroup,
}: DashboardComponentProps) => {
  const [t] = useTranslation();
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');

  return (
    <MetricContext.Provider
      value={{
        aggregation: selectedTimeUnit,
        team: entityGroup,
        project: entityName,
      }}
    >
      <Page themeId="tool">
        <Header
          title="OpenDORA (by Devoteam)"
          subtitle="Through insight to perfection"
        >
          <SupportButton>Plugin for displaying DORA Metrics</SupportButton>
        </Header>
        <Content>
          <Grid container spacing={3} direction="column">
            <Grid container>
              <Grid item xs={12} className="gridBorder">
                <div className="gridBoxText">
                  <Grid container>
                    <Grid item xs={4}>
                      <DropdownComponent
                        onSelect={setSelectedTimeUnit}
                        selection={selectedTimeUnit}
                      />
                    </Grid>
                  </Grid>
                </div>
              </Grid>
              <BenchmarkGridItem type="df" />
              <ChartGridItem
                type="df_count"
                label={t('deployment_frequency.labels.deployment_frequency')}
              />
              <ChartGridItem
                type="df_average"
                label={t(
                  'deployment_frequency.labels.deployment_frequency_average',
                )}
              />
            </Grid>
            <Grid item />
          </Grid>
        </Content>
      </Page>
    </MetricContext.Provider>
  );
};

export const EntityDashboardComponent = () => {
  const { entity } = useEntity();
  const entityGroup = getEntityRelations(entity, 'ownedBy')[0]?.name;
  const entityName = entity.metadata.name;

  return (
    <DashboardComponent entityName={entityName} entityGroup={entityGroup} />
  );
};
