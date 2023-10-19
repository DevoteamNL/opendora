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

const useEntityDetails = () => {
  try {
    const { entity } = useEntity();
    const groupName = getEntityRelations(entity, 'ownedBy')[0]?.name;
    const entityName = entity.metadata.name;

    return { group: groupName, name: entityName };
  } catch (e: unknown) {
    return { group: undefined, name: undefined };
  }
};

export const DashboardComponent = () => {
  const [chartData, setChartData] = React.useState<MetricData | null>(null);
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');
  const [dataError, setDataError] = React.useState<Error | undefined>(
    undefined,
  );

  const groupDataService = useApi(groupDataServiceApiRef);
  const entity = useEntityDetails();

  useEffect(() => {
    groupDataService
      .retrieveMetricDataPoints({
        type: 'df_count',
        team: entity.group,
        aggregation: selectedTimeUnit,
        project: entity.name,
      })
      .then(
        response => {
          setChartData(response);
        },
        (error: Error) => {
          setDataError(error);
        },
      );
  }, [entity, selectedTimeUnit, groupDataService]);

  const chartOrProgressComponent = chartData ? (
    <BarChartComponent metricData={chartData} />
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
                {dataError ? (
                  <ResponseErrorPanel error={dataError} />
                ) : (
                  chartOrProgressComponent
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
