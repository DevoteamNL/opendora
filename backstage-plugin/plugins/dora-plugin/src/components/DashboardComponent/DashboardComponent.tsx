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
import { MetricData } from '../../models/MetricData';
import { groupDataServiceApiRef } from '../../services/GroupDataService';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';
import './DashboardComponent.css';

export interface DashboardComponentProps {
  entityName?: string;
  entityGroup?: string;
}
export const DashboardComponent = ({
  entityName,
  entityGroup,
}: DashboardComponentProps) => {
  const [chartData, setChartData] = React.useState<MetricData | null>(null);
  const [chartDataAverage, setChartDataAverage] = React.useState<MetricData | null>(null);
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');
  const [dataError, setDataError] = React.useState<{countError: Error | null, averageError: Error | null}>(
    {countError: null, averageError: null}
  );

  const groupDataService = useApi(groupDataServiceApiRef);

  useEffect(() => {
    groupDataService
      .retrieveMetricDataPoints({
        type: 'df_count',
        team: entityGroup,
        aggregation: selectedTimeUnit,
        project: entityName,
      })
      .then(
        response => {
          setChartDataAverage(response);
        },
        (error: Error) => {
          setDataError({
            ...dataError,
            countError: error,
          });
        }
      );


    groupDataService
    .retrieveMetricDataPoints({
      type: 'df_average',
      team: entityGroup,
      aggregation: selectedTimeUnit,
      project: entityName,
    })
    .then(
      response => {
        setChartData(response);
      },
      (error: Error) => {
        setDataError({
          ...dataError,
          averageError: error,
        });
      },
    );
  }, [entityGroup, entityName, selectedTimeUnit, groupDataService, dataError]);

  const chartOrProgressComponent = chartData ? (
    <BarChartComponent metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const chartOrProgressComponentAverage = chartDataAverage ? (
    <BarChartComponent metricData={chartDataAverage} />
  ) : (
    <Progress variant="indeterminate" />
  );

  return (
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
                <h1>Deployment Frequency</h1>
                {dataError.countError ? (
                  <ResponseErrorPanel error={dataError.countError} />
                ) : (
                  chartOrProgressComponent
                )}
              </div>
            </Grid>
            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <h1>Deployment Frequency Average</h1>
                {dataError.averageError ? (
                  <ResponseErrorPanel error={dataError.averageError} />
                ) : (
                  chartOrProgressComponentAverage
                )}
              </div>
            </Grid>
          </Grid>

          <Grid item />
        </Grid>
      </Content>
    </Page>
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
