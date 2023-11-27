import {
  Content,
  Header,
  Page,
  Progress,
  ResponseErrorPanel,
  SupportButton,
} from '@backstage/core-components';
import { useApi } from '@backstage/core-plugin-api';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { Grid } from '@material-ui/core';
import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import '../../i18n';
import { dfBenchmarkKey } from '../../models/DfBenchmarkData';
import { groupDataServiceApiRef } from '../../services/GroupDataService';
import { MetricContext } from '../../services/MetricContext';
import { useMetricData } from '../../services/MetricDataHook';
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

export const DashboardComponent = ({
  entityName,
  entityGroup,
}: DashboardComponentProps) => {
  // Overview
  const [dfOverview, setDfOverview] = React.useState<dfBenchmarkKey | null>(
    null,
  );

  const [t] = useTranslation();
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');

  const groupDataService = useApi(groupDataServiceApiRef);

  useEffect(() => {
    groupDataService.retrieveBenchmarkData({ type: 'df' }).then(
      response => {
        setDfOverview(response.key);
      },
      (error: Error) => {
        console.error(error);
        // dispatch({
        //   type: 'dfBenchmarkError',
        //   error: error,
        // });
      },
    );
  }, [entityGroup, entityName, selectedTimeUnit, groupDataService]);

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
              <Grid item xs={12} className="gridBorder">
                <div className="gridBoxText">
                  <Grid container>
                    <Grid item xs={3}>
                      <HighlightTextBoxComponent
                        title=""
                        text=""
                        highlight={
                          dfOverview
                            ? t(
                                `deployment_frequency.overall_labels.${dfOverview}`,
                              )
                            : t('custom_errors.data_unavailable')
                        }
                        // to do: think of text colouring for different scenarios
                        textColour="positiveHighlight"
                      />
                    </Grid>
                  </Grid>
                </div>
              </Grid>
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
